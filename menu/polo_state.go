package menu

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type poloEntry struct {
	Name    string
	Path    string
	IsDir   bool
	Size    int64
	ModTime time.Time
}

// Reads a directory and sorts folders before files
func loadPoloEntries(dir string) ([]poloEntry, error) {
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	entries := make([]poloEntry, 0, len(dirEntries))
	for _, dirEntry := range dirEntries {
		entry := poloEntry{
			Name:  dirEntry.Name(),
			Path:  filepath.Join(dir, dirEntry.Name()),
			IsDir: dirEntry.IsDir(),
		}

		if info, err := dirEntry.Info(); err == nil {
			entry.Size = info.Size()
			entry.ModTime = info.ModTime()
		}

		entries = append(entries, entry)
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir != entries[j].IsDir {
			return entries[i].IsDir
		}

		return strings.ToLower(entries[i].Name) < strings.ToLower(entries[j].Name)
	})

	return entries, nil
}

// Returns a short listing for the right pane directory preview
func previewDirLines(dir string, maxLines int) []string {
	entries, err := loadPoloEntries(dir)
	if err != nil {
		return []string{"(unable to read directory)"}
	}

	if len(entries) == 0 {
		return []string{"(empty directory)"}
	}

	if maxLines <= 0 {
		return []string{"(preview area too small)"}
	}

	visibleLines := maxLines
	showRemainder := len(entries) > maxLines
	if showRemainder && maxLines > 1 {
		visibleLines = maxLines - 1
	}

	if visibleLines > len(entries) {
		visibleLines = len(entries)
	}

	lines := make([]string, 0, maxLines)
	for i := 0; i < visibleLines; i++ {
		prefix := "[F]"
		if entries[i].IsDir {
			prefix = "[D]"
		}

		lines = append(lines, prefix+" "+entries[i].Name)
	}

	if showRemainder {
		if maxLines == 1 {
			return []string{fmt.Sprintf("... %d items", len(entries))}
		}

		lines = append(lines, fmt.Sprintf("... %d more", len(entries)-visibleLines))
	}

	return lines
}

// Reads a window of text preview lines and reports the furthest scroll point
func previewFileLines(path string, scrollLine int, maxLines int) ([]string, int) {
	if maxLines <= 0 {
		return []string{"(preview area too small)"}, 0
	}

	file, err := os.Open(path)
	if err != nil {
		return []string{"(unable to open file)"}, 0
	}
	defer file.Close()

	if !isTextFile(file) {
		return []string{"(binary or unsupported file)"}, 0
	}

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024), 1024*1024)

	allLines := make([]string, 0)
	for scanner.Scan() {
		allLines = append(allLines, strings.ReplaceAll(scanner.Text(), "\t", "    "))
	}

	if err := scanner.Err(); err != nil {
		return []string{"(error reading file)"}, 0
	}

	if len(allLines) == 0 {
		return []string{"(empty file)"}, 0
	}

	maxScroll := clamp(len(allLines)-maxLines, 0, len(allLines))
	scrollLine = clamp(scrollLine, 0, maxScroll)

	end := scrollLine + maxLines
	if end > len(allLines) {
		end = len(allLines)
	}

	return append([]string(nil), allLines[scrollLine:end]...), maxScroll
}

func previewFileScrollLimit(path string, maxLines int) int {
	_, maxScroll := previewFileLines(path, 0, maxLines)
	return maxScroll
}

func formatPreviewWindow(scrollLine int, visibleLines int, maxScroll int) string {
	if maxScroll == 0 {
		return "Preview position: top"
	}

	startLine := scrollLine + 1
	endLine := scrollLine + visibleLines
	if endLine < startLine {
		endLine = startLine
	}

	return fmt.Sprintf("Preview lines: %d-%d", startLine, endLine)
}

// Returns the next directory up while stopping at the filesystem root
func parentDir(dir string) string {
	parent := filepath.Dir(dir)
	if parent == "." || parent == "" {
		return dir
	}

	return parent
}

// Shortens long strings so they fit inside a pane
func trimForPane(text string, width int) string {
	if width <= 0 {
		return ""
	}

	if len(text) <= width {
		return text
	}

	if width <= 3 {
		return text[:width]
	}

	return text[:width-3] + "..."
}

// Trims and clips pane content so it fits the current window
func formatPaneLines(lines []string, width int, maxRows int) []string {
	if maxRows <= 0 {
		return nil
	}

	if len(lines) > maxRows {
		lines = lines[:maxRows]
	}

	formatted := make([]string, 0, len(lines))
	for _, line := range lines {
		formatted = append(formatted, trimForPane(line, width))
	}

	return formatted
}

// Clamp keeps a value within the provided bounds
func clamp(value, min, max int) int {
	if max < min {
		return min
	}

	if value < min {
		return min
	}

	if value > max {
		return max
	}

	return value
}

// Prints file sizes in a compact human-readable form
func formatSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	}

	if size < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(size)/1024)
	}

	return fmt.Sprintf("%.1f MB", float64(size)/(1024*1024))
}

func formatFileType(entry poloEntry) string {
	if entry.IsDir {
		return "directory"
	}

	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(entry.Name)), ".")
	if ext == "" {
		return "file"
	}

	return ext + " file"
}

func formatModTime(modTime time.Time) string {
	if modTime.IsZero() {
		return "unknown"
	}

	return modTime.Format("2006-01-02 15:04")
}

func findEntryIndex(entries []poloEntry, path string) int {
	for i, entry := range entries {
		if entry.Path == path {
			return i
		}
	}

	return 0
}

func selectedDirectoryPath(currentDir string, entries []poloEntry, selected int) string {
	if len(entries) == 0 {
		return currentDir
	}

	selected = clamp(selected, 0, len(entries)-1)
	if entries[selected].IsDir {
		return entries[selected].Path
	}

	return currentDir
}

func poloStateFilePath() string {
	return filepath.Join(os.TempDir(), "goline-polo-path.txt")
}

func isTextFile(file *os.File) bool {
	sample := make([]byte, 512)
	readCount, err := file.Read(sample)
	if err != nil && err != io.EOF {
		return false
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return false
	}

	return !bytes.Contains(sample[:readCount], []byte{0})
}
