package core

import (
	"strings"
	"testing"
)

func TestScanner_Scan(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		matcherPattern  string
		before          int
		after           int
		expectedCount   int
		expectedLines   []string
		expectedIsMatch []bool
	}{
		{
			name:            "Basic match without context",
			input:           "line1 foo\nline2 bar\nline3 foo again\nline4\nline5 foo",
			matcherPattern:  "foo",
			expectedCount:   3,
			expectedLines:   []string{"line1 foo", "line3 foo again", "line5 foo"},
			expectedIsMatch: []bool{true, true, true},
		},
		{
			name:            "Match with Before context",
			input:           "line1\nline2\ntarget match\nline4\nline5\nline6 target\nline7",
			matcherPattern:  "target",
			before:          2,
			expectedCount:   6,
			expectedLines:   []string{
				"line1", "line2", "target match", // matches target match
				"line4", "line5", "line6 target", // matches line6 target
			},
			expectedIsMatch: []bool{false, false, true, false, false, true},
		},
		{
			name:            "Match with After context",
			input:           "line1\nline2 target\nline3\nline4\nline5 target\nline6",
			matcherPattern:  "target",
			after:           1,
			expectedCount:   4,
			expectedLines:   []string{"line2 target", "line3", "line5 target", "line6"},
			expectedIsMatch: []bool{true, false, true, false},
		},
		{
			name:            "Overlapping context (-C 1 with adjacent matches)",
			input:           "line1\nline2 foo\nline3 foo\nline4",
			matcherPattern:  "foo",
			before:          1,
			after:           1,
			expectedCount:   4, // line1, line2 (match), line3 (match), line4
			expectedLines:   []string{"line1", "line2 foo", "line3 foo", "line4"},
			expectedIsMatch: []bool{false, true, true, false},
		},
		{
			name:            "Overlapping context (-C 2 with matches close by)",
			input:           "line1\nline2\nmatch one\nline4\nmatch two\nline6",
			matcherPattern:  "match",
			before:          2,
			after:           2,
			expectedCount:   6, // All lines
			expectedLines:   []string{"line1", "line2", "match one", "line4", "match two", "line6"},
			expectedIsMatch: []bool{false, false, true, false, true, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher, err := NewMatcher(tt.matcherPattern, false, false)
			if err != nil {
				t.Fatalf("Failed to compile matcher: %v", err)
			}
			
			scanner := &Scanner{
				Matcher: matcher,
				Before:  tt.before,
				After:   tt.after,
			}

			reader := strings.NewReader(tt.input)
			resultsChan := make(chan MatchResult, 100)

			go func() {
				err := scanner.Scan(reader, "testfile", resultsChan)
				if err != nil {
					t.Errorf("Scan failed: %v", err)
				}
				close(resultsChan)
			}()

			var results []MatchResult
			for res := range resultsChan {
				results = append(results, res)
			}

			if len(results) != tt.expectedCount {
				t.Fatalf("Expected %d results, got %d", tt.expectedCount, len(results))
			}

			// Validate emitted lines
			for i, expectedLine := range tt.expectedLines {
				if i < len(results) {
					if results[i].Line != expectedLine {
						t.Errorf("Result %d: Expected line %q, got %q", i, expectedLine, results[i].Line)
					}
					if results[i].IsMatch != tt.expectedIsMatch[i] {
						t.Errorf("Result %d: Expected IsMatch=%v, got %v for line %q", i, tt.expectedIsMatch[i], results[i].IsMatch, expectedLine)
					}
				}
			}
		})
	}
}
