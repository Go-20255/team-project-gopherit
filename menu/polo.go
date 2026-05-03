package menu

import (
	"os"
	"path/filepath"

	"github.com/eiannone/keyboard"
)

const (
	poloListStartRow = 5
	poloHelpLines    = 4
)

type poloBrowser struct {
	menu       *Menu
	currentDir string
	entries    []poloEntry
	selected   int
	scrollLine int
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
		scrollLine: 0,
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
	b.menu.handle = func(key keyboard.Key, letter rune) bool {
		switch letter {
		case ',':
			b.scroll(-5)
		case '.':
			b.scroll(5)
		default:
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
				b.finishSelection()
			default:
				return false
			}
		}

		_ = b.render()
		return true
	}

	b.menu.Run()
	return nil
}

// Refreshes the browser state for a new current directory
func (b *poloBrowser) loadDir(dir string) error {
	return b.loadDirAt(dir, "")
}

// Refreshes the browser state and optionally restores a highlighted path
func (b *poloBrowser) loadDirAt(dir string, selectedPath string) error {
	entries, err := loadPoloEntries(dir)
	if err != nil {
		return err
	}

	b.entries = entries
	b.selected = 0
	if selectedPath != "" {
		b.selected = findEntryIndex(entries, selectedPath)
	}
	b.currentDir = dir
	b.scrollLine = 0
	b.status = ""
	b.menu.prevPane()
	b.menu.Col = 1
	b.setCursorRow(b.cursorRow())

	return nil
}

// Changes the selected entry while keeping it inside the visible list
func (b *poloBrowser) move(delta int) {
	if len(b.entries) == 0 {
		return
	}

	b.selected = clamp(b.selected+delta, 0, len(b.entries)-1)
	b.scrollLine = 0
	b.setCursorRow(b.cursorRow())
}

func (b *poloBrowser) scroll(delta int) {
	// Only files use scroll because directory previews are already clipped
	if len(b.entries) == 0 || b.entries[b.selected].IsDir {
		b.scrollLine = 0
		return
	}

	maxScroll := previewFileScrollLimit(b.entries[b.selected].Path, b.filePreviewRows(b.entries[b.selected]))
	b.scrollLine = clamp(b.scrollLine+delta, 0, maxScroll)
}

// Moves the browser one directory up if possible
func (b *poloBrowser) goParent() {
	nextDir := parentDir(b.currentDir)
	if nextDir == b.currentDir {
		b.status = "Already at the filesystem root."
		return
	}

	if err := b.loadDirAt(nextDir, b.currentDir); err != nil {
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

	paneWidth := b.menu.displayPaneWidth()
	b.writePaneLines(b.leftPaneLines(), 0, paneWidth)
	b.writePaneLines(b.rightPaneLines(), 1, paneWidth)

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
		"Polo Browser",
		"Current: " + b.currentDir,
		"Review the target path before finishing",
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
	lines := b.selectionHeaderLines()

	if len(b.entries) == 0 {
		lines = append(lines, "No entries in this directory")
		return b.appendStatusAndHelp(lines)
	}

	entry := b.entries[b.selected]
	if entry.IsDir {
		lines = append(lines,
			"[D] "+entry.Name,
			"Space will finish in this directory",
			"Directory preview:",
		)
		lines = append(lines, previewDirLines(entry.Path, b.previewRows(len(lines)))...)
		return b.appendStatusAndHelp(lines)
	}

	lines = append(lines, b.filePreviewBaseLines(entry)...)
	// Keeps the preview window separate from the metadata and help footer
	previewLines, maxScroll := previewFileLines(entry.Path, b.scrollLine, b.filePreviewRows(entry))
	b.scrollLine = clamp(b.scrollLine, 0, maxScroll)
	lines = append(lines, formatPreviewWindow(b.scrollLine, len(previewLines), maxScroll))
	lines = append(lines, previewLines...)

	return b.appendStatusAndHelp(lines)
}

// Adds status text and the key hints
func (b *poloBrowser) appendStatusAndHelp(lines []string) []string {
	if b.status != "" {
		lines = append(lines, "", b.status)
	}
	lines = append(lines,
		"",
		"Up/Down move through entries",
		"Right opens and Left goes back",
		"Space uses the target path shown above",
		", up preview  . down preview  Esc quit",
	)
	return lines
}

// Writes the selected directory path for the wrapper script
func (b *poloBrowser) finishSelection() {
	target := b.finishTarget()
	if err := os.WriteFile(poloStateFilePath(), []byte(target), 0644); err != nil {
		b.status = "Unable to save selected directory."
		return
	}

	clear()
	os.Exit(0)
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

// Tracks how many rows the footer uses in the preview pane
func (b *poloBrowser) footerLines() int {
	footerLines := 1 + poloHelpLines
	if b.status != "" {
		footerLines += 2
	}

	return footerLines
}

// Calculates the remaining space for preview content
func (b *poloBrowser) previewRows(headerLines int) int {
	previewRows := b.lastDrawRow() - headerLines - b.footerLines()
	if previewRows < 0 {
		return 0
	}

	return previewRows
}

// Keeps the finish behavior visible while browsing
func (b *poloBrowser) selectionHeaderLines() []string {
	// Reminds the user where Space will leave them before they confirm
	return []string{
		"Selection",
		"Finish target: " + b.finishTarget(),
		"",
	}
}

func (b *poloBrowser) filePreviewBaseLines(entry poloEntry) []string {
	return []string{
		"[F] " + entry.Name,
		"Selected file preview only",
		"Type: " + formatFileType(entry),
		"Size: " + formatSize(entry.Size),
		"Modified: " + formatModTime(entry.ModTime),
		"Preview:",
	}
}

func (b *poloBrowser) filePreviewRows(entry poloEntry) int {
	// Reserves room for the selection header metadata preview label and footer
	headerLines := len(b.selectionHeaderLines()) + len(b.filePreviewBaseLines(entry)) + 1
	return b.previewRows(headerLines)
}

func (b *poloBrowser) finishTarget() string {
	return selectedDirectoryPath(b.currentDir, b.entries, b.selected)
}

// Keeps the cursor bounds code a little easier to read
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
