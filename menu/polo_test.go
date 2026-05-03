package menu

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGoParentRestoresChildSelection(t *testing.T) {
	parentDir := t.TempDir()
	childDir := filepath.Join(parentDir, "child")
	otherDir := filepath.Join(parentDir, "other")

	for _, dir := range []string{childDir, otherDir} {
		if err := os.Mkdir(dir, 0755); err != nil {
			t.Fatalf("make child directory: %v", err)
		}
	}

	browser := testBrowser(t)
	if err := browser.loadDir(childDir); err != nil {
		t.Fatalf("load child directory: %v", err)
	}

	browser.goParent()

	if browser.currentDir != parentDir {
		t.Fatalf("expected parent dir %s, got %s", parentDir, browser.currentDir)
	}

	if browser.entries[browser.selected].Path != childDir {
		t.Fatalf("expected selected entry to restore child path, got %s", browser.entries[browser.selected].Path)
	}
}

func TestLoadDirAtRestoresRequestedSelection(t *testing.T) {
	startDir := t.TempDir()
	alphaDir := filepath.Join(startDir, "alpha")
	betaDir := filepath.Join(startDir, "beta")

	for _, dir := range []string{alphaDir, betaDir} {
		if err := os.Mkdir(dir, 0755); err != nil {
			t.Fatalf("make directory: %v", err)
		}
	}

	browser := testBrowser(t)
	if err := browser.loadDirAt(startDir, betaDir); err != nil {
		t.Fatalf("load directory with selected path: %v", err)
	}

	if browser.entries[browser.selected].Path != betaDir {
		t.Fatalf("expected restored selection %s, got %s", betaDir, browser.entries[browser.selected].Path)
	}
}

func TestScrollClampsToPreviewLimit(t *testing.T) {
	startDir := t.TempDir()
	filePath := filepath.Join(startDir, "preview.txt")
	content := strings.Join([]string{
		"line one",
		"line two",
		"line three",
		"line four",
		"line five",
		"line six",
	}, "\n")

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("write preview file: %v", err)
	}

	browser := testBrowser(t)
	if err := browser.loadDir(startDir); err != nil {
		t.Fatalf("load directory: %v", err)
	}

	browser.scroll(100)

	maxScroll := previewFileScrollLimit(filePath, browser.filePreviewRows(browser.entries[browser.selected]))
	if browser.scrollLine != maxScroll {
		t.Fatalf("expected scroll line %d, got %d", maxScroll, browser.scrollLine)
	}
}

func testBrowser(t *testing.T) *poloBrowser {
	t.Helper()

	lines := make([]string, 18)
	return &poloBrowser{
		menu: &Menu{
			Row:    poloListStartRow,
			Col:    1,
			Pane:   0,
			Height: 18,
			Width:  80,
			Panes: []Pane{
				{Lines: append([]string(nil), lines...)},
				{Lines: append([]string(nil), lines...)},
			},
		},
	}
}
