package core

import (
	"regexp"
	"strings"
)

// MatchResult holds a successful match.
type MatchResult struct {
	FileName   string
	LineNumber int
	Line       string
	IsMatch    bool
}

// Matcher defines the interface for checking if a line matches a pattern.
type Matcher interface {
	Match(line string) bool
}

// StringMatcher matches exact substrings.
type StringMatcher struct {
	Pattern    string
	IgnoreCase bool
}

func (s *StringMatcher) Match(line string) bool {
	if s.IgnoreCase {
		return strings.Contains(strings.ToLower(line), strings.ToLower(s.Pattern))
	}
	return strings.Contains(line, s.Pattern)
}

// RegexMatcher matches using a compiled regular expression.
type RegexMatcher struct {
	Regex *regexp.Regexp
}

func (r *RegexMatcher) Match(line string) bool {
	return r.Regex.MatchString(line)
}

// NewMatcher creates the appropriate matcher based on flags.
func NewMatcher(pattern string, isRegex, ignoreCase bool) (Matcher, error) {
	if isRegex {
		if ignoreCase {
			pattern = "(?i)" + pattern
		}
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		return &RegexMatcher{Regex: re}, nil
	}

	return &StringMatcher{
		Pattern:    pattern,
		IgnoreCase: ignoreCase,
	}, nil
}
