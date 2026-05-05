package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/team-project-gopherit/internal/ggrep/cli"
	"github.com/team-project-gopherit/internal/ggrep/core"
	"github.com/team-project-gopherit/internal/ggrep/tui"
)

func main() {
	opts, err := cli.ParseFlags()
	if err != nil {
		os.Exit(2)
	}

	matcher, err := core.NewMatcher(opts.Pattern, opts.IsRegex, opts.IgnoreCase)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid pattern: %v\n", err)
		os.Exit(2)
	}

	scanner := &core.Scanner{
		Matcher: matcher,
		After:   opts.After,
		Before:  opts.Before,
	}

	printFn := func(res core.MatchResult) {
		// Output format: [filename:]<line> if multiple files?
		prefix := ""
		if len(opts.Files) > 1 {
			prefix = res.FileName
			if res.IsMatch {
				prefix += ":"
			} else {
				prefix += "-"
			}
		}
		if prefix != "" {
			fmt.Fprintf(os.Stdout, "%s%s\n", prefix, res.Line)
		} else {
			fmt.Fprintln(os.Stdout, res.Line)
		}
	}

	if opts.Interactive {
		app := tui.InitialModel(opts)
		p := tea.NewProgram(app, tea.WithAltScreen(), tea.WithOutput(os.Stderr))
		finalModel, err := p.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
			os.Exit(1)
		}

		// Print final results to stdout for piping
		if m, ok := finalModel.(*tui.Model); ok {
			for _, res := range m.GetFinalResults() {
				printFn(res)
			}
		}
		return
	}

	if len(opts.Files) == 0 {
		// Read from stdin
		resultsChan := make(chan core.MatchResult)
		go func() {
			_ = scanner.Scan(os.Stdin, "(standard input)", resultsChan)
			close(resultsChan)
		}()
		for res := range resultsChan {
			printFn(res)
		}
		return
	}

	// Read from files concurrently
	numWorkers := 4 // Arbitrary default worker count for MVP
	workerPool := core.NewWorkerPool(scanner, opts.Files, printFn)
	workerPool.Run(numWorkers)
}
