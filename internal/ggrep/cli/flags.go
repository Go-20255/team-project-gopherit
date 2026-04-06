package cli

import (
	"flag"
	"fmt"
	"os"
)

// Options holds the parsed command-line arguments.
type Options struct {
	IgnoreCase bool
	IsRegex    bool
	Pattern    string
	Files      []string
	After      int
	Before     int
}

// ParseFlags parses standard command-line flags and returns the configuration.
func ParseFlags() (*Options, error) {
	ignoreCase := flag.Bool("i", false, "Ignore case distinctions in both the PATTERN and the input files")
	isRegex := flag.Bool("E", false, "Interpret pattern as an extended regular expression")
	fixedString := flag.Bool("F", false, "Interpret pattern as a fixed string (default)")
	
	after := flag.Int("A", 0, "Print NUM lines of trailing context after matching lines")
	before := flag.Int("B", 0, "Print NUM lines of leading context before matching lines")
	context := flag.Int("C", 0, "Print NUM lines of output context")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTION]... PATTERN [FILE]...\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Search for PATTERN in each FILE.\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		return nil, fmt.Errorf("missing pattern")
	}

	pattern := args[0]
	files := args[1:]

	regex := *isRegex
	if *fixedString {
		regex = false
	}

	a := *after
	b := *before
	if *context > 0 {
		if *after == 0 {
			a = *context
		}
		if *before == 0 {
			b = *context
		}
	}

	return &Options{
		IgnoreCase: *ignoreCase,
		IsRegex:    regex,
		Pattern:    pattern,
		Files:      files,
		After:      a,
		Before:     b,
	}, nil
}
