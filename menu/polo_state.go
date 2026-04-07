package menu

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type poloEntry struct {
	Name  string
	Path  string
	IsDir bool
	Size  int64
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
		if !dirEntry.IsDir() {
			info, err := dirEntry.Info()
			if err == nil {
				entry.Size = info.Size()
			}
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

	if maxLines > len(entries) {
		maxLines = len(entries)
	}

	lines := make([]string, 0, maxLines+1)
	for i := 0; i < maxLines; i++ {
		prefix := "[F]"
		if entries[i].IsDir {
			prefix = "[D]"
		}

		lines = append(lines, prefix+" "+entries[i].Name)
	}

	if len(entries) > maxLines {
		lines = append(lines, fmt.Sprintf("... %d more", len(entries)-maxLines))
	}

	return lines
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
