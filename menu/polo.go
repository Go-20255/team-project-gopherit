package menu

import (
	"os"
	"path/filepath"

	"github.com/eiannone/keyboard"
)

const (
	poloListStartRow = 4
	poloLeftWidth    = 38
	poloRightWidth   = 38
	poloPreviewLines = 8
)

type poloBrowser struct {
	menu       *Menu
	currentDir string
	entries    []poloEntry
	selected   int
	status     string
}

// Starts the interactive browser from the provided directory
func RunPolo(startDir string) error {
	if startDir == "" {
		var err error
		startDir, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	absDir, err := filepath.Abs(startDir)
	if err != nil {
		return err
	}

	browser, err := newPoloBrowser(absDir)
	if err != nil {
		return err
	}

	return browser.Run()
}

// Builds the browser state and loads the first directory view
func newPoloBrowser(startDir string) (*poloBrowser, error) {
	browser := &poloBrowser{
		menu:       Init(),
		currentDir: startDir,
	}

	if err := browser.loadDir(startDir); err != nil {
		return nil, err
	}

	return browser, nil
}

// Listens for key presses and updates the browser until the given user exits
func (b *poloBrowser) Run() error {
	b.menu.setup = func() {
		_ = b.render()
	}
	b.menu.handle = func(key keyboard.Key) bool {
		switch key {
		case keyboard.KeyArrowUp:
			b.move(-1)
		case keyboard.KeyArrowDown:
			b.move(1)
		case keyboard.KeyArrowLeft:
			b.goParent()
		case keyboard.KeyArrowRight, keyboard.KeyEnter:
			b.enterSelected()
		case keyboard.KeySpace:
			return true
		default:
			return false
		}

		_ = b.render()
		return true
	}

	b.menu.Run()
	return nil
}

// Refreshes the browser state for a new current directory
func (b *poloBrowser) loadDir(dir string) error {
	entries, err := loadPoloEntries(dir)
	if err != nil {
		return err
	}

	b.entries = entries
	b.selected = clamp(b.selected, 0, len(entries)-1)
	b.status = ""
	b.menu.prevPane()
	b.menu.Col = 1
	// When returning to parent, loop over entries to find
	// The directory we were previously in
	for i, entry := range entries {
		if entry.Path == b.currentDir {
			b.move(i - b.selected)
		}
	}
	b.currentDir = dir
	b.setCursorRow(b.cursorRow())

	return nil
}

// Changes the selected entry while keeping it inside the visible list
func (b *poloBrowser) move(delta int) {
	if len(b.entries) == 0 {
		return
	}

	b.selected = clamp(b.selected+delta, 0, len(b.entries)-1)
	b.setCursorRow(b.cursorRow())
}

// Moves the browser one directory up if possible
func (b *poloBrowser) goParent() {
	nextDir := parentDir(b.currentDir)
	if nextDir == b.currentDir {
		b.status = "Already at the filesystem root."
		return
	}

	if err := b.loadDir(nextDir); err != nil {
		b.status = "Unable to open parent directory."
	}
}

// Opens the highlighted directory
func (b *poloBrowser) enterSelected() {
	if len(b.entries) == 0 {
		b.status = "This directory is empty."
		return
	}

	entry := b.entries[b.selected]
	if !entry.IsDir {
		b.status = "Selected item is not a directory."
		return
	}

	if err := b.loadDir(entry.Path); err != nil {
		b.status = "Unable to open selected directory."
	}
}

// Rebuilds both panes from the current browser state and draws them
func (b *poloBrowser) render() error {
	b.clearPane(0)
	b.clearPane(1)

	b.writePaneLines(b.leftPaneLines(), 0, poloLeftWidth)
	b.writePaneLines(b.rightPaneLines(), 1, poloRightWidth)

	b.menu.Pane = 0
	b.menu.Col = 1
	b.setCursorRow(b.cursorRow())
	b.menu.draw()

	return nil
}

// Removes any old lines before the next render pass
func (b *poloBrowser) clearPane(pane int) {
	for i := range b.menu.Panes[pane].Lines {
		b.menu.write("", i, pane)
	}
}

// Uses the existing batch writer to fill a pane from the top
func (b *poloBrowser) writePaneLines(lines []string, pane int, width int) {
	b.menu.batchWrite(formatPaneLines(lines, width, b.lastDrawRow()), pane, 1)
}

// Formats the directory listing shown in the main pane
func (b *poloBrowser) leftPaneLines() []string {
	lines := []string{
		"Polo",
		b.currentDir,
		"",
	}

	if len(b.entries) == 0 {
		return append(lines, "(empty directory)")
	}

	start, end := b.visibleRange()
	for i := start; i < end; i++ {
		entry := b.entries[i]
		prefix := "  "
		if i == b.selected {
			prefix = "> "
		}

		entryType := "[F]"
		if entry.IsDir {
			entryType = "[D]"
		}

		lines = append(lines, prefix+entryType+" "+entry.Name)
	}

	return lines
}

// Formats details for the currently selected entry
func (b *poloBrowser) rightPaneLines() []string {
	lines := []string{"Selected"}

	if len(b.entries) == 0 {
		lines = append(lines, "(directory is empty)")
		return b.appendStatusAndHelp(lines)
	}
	entry := b.entries[b.selected]
	entryType := "[F]"
	if entry.IsDir {
		entryType = "[D]"
	}

	lines = append(lines, entryType+" "+entry.Name)
	if entry.IsDir {
		lines = append(lines, "Contents:")
		lines = append(lines, previewDirLines(entry.Path, poloPreviewLines)...)
	} else {
		lines = append(lines, "Size: "+formatSize(entry.Size))
	}

	return b.appendStatusAndHelp(lines)
}

// Adds status text and the key hints
func (b *poloBrowser) appendStatusAndHelp(lines []string) []string {
	if b.status != "" {
		lines = append(lines, "", b.status)
	}
	lines = append(lines,
		"",
		"Up/Down: move",
		"Right/Enter: open",
		"Left: parent",
		"Esc: quit",
	)
	return lines
}

// Choosses which slice of entries fits in the current window
func (b *poloBrowser) visibleRange() (int, int) {
	maxVisible := b.maxVisibleEntries()
	if maxVisible <= 0 || len(b.entries) <= maxVisible {
		return 0, len(b.entries)
	}
	start := 0
	if b.selected >= maxVisible {
		start = b.selected - maxVisible + 1
	}

	end := start + maxVisible
	if end > len(b.entries) {
		end = len(b.entries)
	}
	return start, end
}

// Maps the selected entry to the row used by the menu cursor
func (b *poloBrowser) cursorRow() int {
	start, _ := b.visibleRange()
	row := poloListStartRow + (b.selected - start)
	return clamp(row, poloListStartRow, maxInt(poloListStartRow, b.lastDrawRow()))
}

// Reuses the menu cursor movement helpers to reach a target row
func (b *poloBrowser) setCursorRow(target int) {
	switch {
	case target > b.menu.Row:
		b.menu.down(target - b.menu.Row)
	case target < b.menu.Row:
		b.menu.up(b.menu.Row - target)
	}
}

// Reports how many list rows are available for entries
func (b *poloBrowser) maxVisibleEntries() int {
	maxVisible := b.lastDrawRow() - poloListStartRow + 1
	if maxVisible < 0 {
		return 0
	}
	return maxVisible
}

// Matches the last row used by the existing menu draw loop
func (b *poloBrowser) lastDrawRow() int {
	lastRow := b.menu.Height - 2
	if lastRow < 1 {
		return 1
	}
	return lastRow
}

// Keeps the cursor bounds code a little easier to read
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
