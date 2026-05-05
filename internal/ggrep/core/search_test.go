package core

import (
	"testing"
)

func TestNewMatcher(t *testing.T) {
	tests := []struct {
		name       string
		pattern    string
		isRegex    bool
		ignoreCase bool
		wantErr    bool
	}{
		{"Valid string match", "hello", false, false, false},
		{"Valid regex match", "^hello", true, false, false},
		{"Valid case insensitive regex match", "^hello", true, true, false},
		{"Invalid regex", "[def", true, false, true},
		{"Empty string", "", false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewMatcher(tt.pattern, tt.isRegex, tt.ignoreCase)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMatcher() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMatcherMatch(t *testing.T) {
	tests := []struct {
		name       string
		pattern    string
		isRegex    bool
		ignoreCase bool
		input      string
		expected   bool
	}{
		{"Str exact match", "foo", false, false, "this is foo bar", true},
		{"Str no match", "foo", false, false, "no f o o here", false},
		{"Str case insensitive match", "FOO", false, true, "this is foo", true},
		{"Str case sensitive fail", "FOO", false, false, "this is foo", false},
		{"Regex match basic", "^[0-9]+", true, false, "123 testing", true},
		{"Regex no match", "^[0-9]+", true, false, "abc testing", false},
		{"Regex case insensitive", "FOO", true, true, "this is foo", true},
		{"Str empty strings", "", false, false, "hello from empty match", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matcher, err := NewMatcher(tt.pattern, tt.isRegex, tt.ignoreCase)
			if err != nil {
				t.Fatalf("Unexpected error creating matcher: %v", err)
			}
			result := matcher.Match(tt.input)
			if result != tt.expected {
				t.Errorf("Match(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}
