package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/team-project-gopherit/internal/ggrep/cli"
	"github.com/team-project-gopherit/internal/ggrep/core"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#E88388")).Italic(true)
	matchStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Bold(true)
	fileStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#A8CC8C"))
)

type tickMsg time.Time

type matchMsg struct {
	id  int
	res core.MatchResult
	sub chan core.MatchResult
}

type searchFinishedMsg struct {
	id int
}

type Model struct {
	textInput textinput.Model
	viewport  viewport.Model
	opts      *cli.Options

	results []core.MatchResult

	err          error
	debounceTime time.Time

	ready    bool
	searchID int // Tracks the current active search to ignore outdated matches
}

func InitialModel(opts *cli.Options) *Model {
	ti := textinput.New()
	ti.Placeholder = "Type regex pattern..."
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50

	if opts.Pattern != "" {
		ti.SetValue(opts.Pattern)
	}

	return &Model{
		textInput: ti,
		opts:      opts,
		results:   make([]core.MatchResult, 0),
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		m.scheduleTick(),
	)
}

func (m *Model) scheduleTick() tea.Cmd {
	m.debounceTime = time.Now().Add(300 * time.Millisecond)
	return tea.Tick(300*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(m.debounceTime)
	})
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
		cmds  []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

		// Any key press updates text input and resets debounce timer
		oldVal := m.textInput.Value()
		m.textInput, tiCmd = m.textInput.Update(msg)
		cmds = append(cmds, tiCmd)

		if m.textInput.Value() != oldVal {
			cmds = append(cmds, m.scheduleTick())
		}

	case tea.WindowSizeMsg:
		headerHeight := 3 // Search bar + error space
		footerHeight := 1
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.viewport.SetContent("Ready. Start typing to search...")
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}
		m.updateViewport()

	case tickMsg:
		// Only trigger search if this is the most recently scheduled tick
		if time.Time(msg).Equal(m.debounceTime) {
			pattern := m.textInput.Value()
			if pattern == "" {
				m.err = nil
				m.results = nil
				m.updateViewport()
				return m, nil
			}

			matcher, err := core.NewMatcher(pattern, m.opts.IsRegex, m.opts.IgnoreCase)
			if err != nil {
				// Graceful failure: record error, keep old results, update UI
				m.err = err
				return m, nil
			}

			// Valid regex! Clear results and start new worker pool
			m.err = nil
			m.results = nil
			m.searchID++
			m.updateViewport()
			m.viewport.GotoTop()

			cmds = append(cmds, startSearch(m.searchID, pattern, m.opts, matcher))
		}

	case matchMsg:
		if msg.id == m.searchID {
			m.results = append(m.results, msg.res)
			m.updateViewport()
		}
		// Wait for next match from this search
		cmds = append(cmds, waitForMatch(msg.id, msg.sub))

	case searchFinishedMsg:
		// Do nothing, just stop waiting
	}

	if m.ready {
		m.viewport, vpCmd = m.viewport.Update(msg)
		cmds = append(cmds, vpCmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	header := fmt.Sprintf("%s\n%s",
		titleStyle.Render(" ggrep Interactive Mode "),
		m.textInput.View(),
	)

	errSpace := "\n"
	if m.err != nil {
		errSpace = errorStyle.Render(fmt.Sprintf("\nInvalid regex: %v", m.err))
	}

	footer := fmt.Sprintf("\n%d matches found • press Esc to exit", len(m.results))

	return fmt.Sprintf("%s%s\n%s%s", header, errSpace, m.viewport.View(), footer)
}

func (m *Model) updateViewport() {
	if len(m.results) == 0 {
		m.viewport.SetContent("No matches found.")
		return
	}

	var sb strings.Builder
	for _, res := range m.results {
		prefix := ""
		if len(m.opts.Files) > 1 {
			prefix = fileStyle.Render(res.FileName)
			if res.IsMatch {
				prefix += matchStyle.Render(":")
			} else {
				prefix += matchStyle.Render("-")
			}
		}
		
		line := res.Line
		if res.IsMatch {
			// Highlight the entire matching line for simplicity in TUI
			line = matchStyle.Render(line)
		}
		
		if prefix != "" {
			sb.WriteString(fmt.Sprintf("%s%s\n", prefix, line))
		} else {
			sb.WriteString(fmt.Sprintf("%s\n", line))
		}
	}
	m.viewport.SetContent(sb.String())
}

func startSearch(id int, pattern string, opts *cli.Options, matcher core.Matcher) tea.Cmd {
	return func() tea.Msg {
		scanner := &core.Scanner{
			Matcher: matcher,
			After:   opts.After,
			Before:  opts.Before,
		}

		sub := make(chan core.MatchResult, 100)
		
		wp := core.NewWorkerPool(scanner, opts.Files, func(res core.MatchResult) {
			sub <- res
		})

		go func() {
			numWorkers := 4
			wp.Run(numWorkers)
			close(sub)
		}()

		// Start reading from the channel immediately
		return waitForMatch(id, sub)()
	}
}

func waitForMatch(id int, sub chan core.MatchResult) tea.Cmd {
	return func() tea.Msg {
		res, ok := <-sub
		if !ok {
			return searchFinishedMsg{id: id}
		}
		return matchMsg{id: id, res: res, sub: sub}
	}
}

// GetFinalResults allows the main application to retrieve the matches 
// after exiting so they can be piped to stdout.
func (m *Model) GetFinalResults() []core.MatchResult {
	return m.results
}
