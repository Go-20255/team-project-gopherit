package core

import (
	"os"
	"sync"
	"testing"
)

func TestWorkerPool(t *testing.T) {
	// Create temporary files to test real file scanning behavior within workers
	f1, err := os.CreateTemp("", "test1_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(f1.Name())
	_, _ = f1.WriteString("hello foo\nhello bar")
	f1.Close()

	f2, err := os.CreateTemp("", "test2_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(f2.Name())
	_, _ = f2.WriteString("testing foo\nnothing here")
	f2.Close()

	matcher, _ := NewMatcher("foo", false, false)
	scanner := &Scanner{Matcher: matcher}

	var printedLines []string
	var mu sync.Mutex

	printFn := func(res MatchResult) {
		mu.Lock()
		defer mu.Unlock()
		printedLines = append(printedLines, res.Line)
	}

	workerPool := NewWorkerPool(scanner, []string{f1.Name(), f2.Name()}, printFn)
	
	// Test concurrency with 10 workers (though we only have 2 files)
	workerPool.Run(10)

	// We expect 2 matches in total
	if len(printedLines) != 2 {
		t.Fatalf("Expected 2 matches, got %d", len(printedLines))
	}

	foundHello := false
	foundTesting := false
	for _, l := range printedLines {
		if l == "hello foo" {
			foundHello = true
		}
		if l == "testing foo" {
			foundTesting = true
		}
	}

	if !foundHello || !foundTesting {
		t.Errorf("Missing expected lines. printedLines: %v", printedLines)
	}
}
