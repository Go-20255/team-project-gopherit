package menu

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPreviewDirLinesHandlesSmallPreviewArea(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "notes.txt")

	if err := os.WriteFile(filePath, []byte("hello"), 0644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	lines := previewDirLines(tempDir, 0)
	if len(lines) != 1 || lines[0] != "(preview area too small)" {
		t.Fatalf("expected small preview message, got %v", lines)
	}
}

func TestPreviewFileLinesReturnsWindowAndScrollLimit(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "notes.txt")
	content := strings.Join([]string{
		"line one",
		"line two",
		"line three",
		"line four",
	}, "\n")

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("write preview file: %v", err)
	}

	lines, maxScroll := previewFileLines(filePath, 1, 2)
	if maxScroll != 2 {
		t.Fatalf("expected max scroll 2, got %d", maxScroll)
	}

	if len(lines) != 2 || lines[0] != "line two" || lines[1] != "line three" {
		t.Fatalf("unexpected preview lines %v", lines)
	}
}

func TestPreviewFileLinesHandlesBinaryFiles(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "data.bin")

	if err := os.WriteFile(filePath, []byte{0, 1, 2, 3}, 0644); err != nil {
		t.Fatalf("write binary file: %v", err)
	}

	lines, maxScroll := previewFileLines(filePath, 0, 4)
	if maxScroll != 0 {
		t.Fatalf("expected no scroll for binary preview, got %d", maxScroll)
	}

	if len(lines) != 1 || lines[0] != "(binary or unsupported file)" {
		t.Fatalf("unexpected binary preview lines %v", lines)
	}
}

func TestPreviewFileLinesHandlesEmptyFiles(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "empty.txt")

	if err := os.WriteFile(filePath, []byte{}, 0644); err != nil {
		t.Fatalf("write empty file: %v", err)
	}

	lines, maxScroll := previewFileLines(filePath, 0, 4)
	if maxScroll != 0 {
		t.Fatalf("expected no scroll for empty file, got %d", maxScroll)
	}

	if len(lines) != 1 || lines[0] != "(empty file)" {
		t.Fatalf("unexpected empty file preview lines %v", lines)
	}
}

func TestPreviewDirLinesShowsRemainderMessage(t *testing.T) {
	tempDir := t.TempDir()
	for _, name := range []string{"alpha", "beta", "gamma"} {
		filePath := filepath.Join(tempDir, name)
		if err := os.WriteFile(filePath, []byte(name), 0644); err != nil {
			t.Fatalf("write preview file: %v", err)
		}
	}

	lines := previewDirLines(tempDir, 2)
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %v", lines)
	}

	if lines[1] != "... 2 more" {
		t.Fatalf("expected remainder line, got %v", lines)
	}
}

func TestSelectedDirectoryPathUsesHighlightedDirectory(t *testing.T) {
	entries := []poloEntry{
		{Name: "notes.txt", Path: "/tmp/notes.txt"},
		{Name: "project", Path: "/tmp/project", IsDir: true},
	}

	selectedPath := selectedDirectoryPath("/tmp", entries, 1)
	if selectedPath != "/tmp/project" {
		t.Fatalf("expected selected directory path, got %s", selectedPath)
	}
}

func TestSelectedDirectoryPathFallsBackToCurrentDirectoryForFiles(t *testing.T) {
	entries := []poloEntry{
		{Name: "notes.txt", Path: "/tmp/notes.txt"},
	}

	selectedPath := selectedDirectoryPath("/tmp", entries, 0)
	if selectedPath != "/tmp" {
		t.Fatalf("expected current directory fallback, got %s", selectedPath)
	}
}

func TestFormatPreviewWindowReportsScrolledRange(t *testing.T) {
	label := formatPreviewWindow(3, 5, 10)
	if label != "Preview lines: 4-8" {
		t.Fatalf("unexpected preview window label %s", label)
	}
}

func TestFormatFileTypeUsesExtension(t *testing.T) {
	entry := poloEntry{Name: "notes.md"}
	if formatFileType(entry) != "md file" {
		t.Fatalf("unexpected file type label %s", formatFileType(entry))
	}
}
