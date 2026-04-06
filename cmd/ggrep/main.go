package main

import (
	"fmt"
	"os"

	"github.com/team-project-gopherit/internal/ggrep/cli"
	"github.com/team-project-gopherit/internal/ggrep/core"
)

func main() {
	opts, err := cli.ParseFlags()
	if err != nil {
		os.Exit(2)
	}

	matcher, err := core.NewMatcher(opts.Pattern, opts.IsRegex, opts.IgnoreCase)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid pattern: %v\n", err)
		os.Exit(2)
	}

	scanner := &core.Scanner{
		Matcher: matcher,
		After:   opts.After,
		Before:  opts.Before,
	}

	printFn := func(res core.MatchResult) {
		// Output format: [filename:]<line> if multiple files?
		prefix := ""
		if len(opts.Files) > 1 {
			prefix = res.FileName
			if res.IsMatch {
				prefix += ":"
			} else {
				prefix += "-"
			}
		}
		if prefix != "" {
			fmt.Printf("%s%s\n", prefix, res.Line)
		} else {
			fmt.Println(res.Line)
		}
	}

	if len(opts.Files) == 0 {
		// Read from stdin
		resultsChan := make(chan core.MatchResult)
		go func() {
			_ = scanner.Scan(os.Stdin, "(standard input)", resultsChan)
			close(resultsChan)
		}()
		for res := range resultsChan {
			printFn(res)
		}
		return
	}

	// Read from files concurrently
	numWorkers := 4 // Arbitrary default worker count for MVP
	workerPool := core.NewWorkerPool(scanner, opts.Files, printFn)
	workerPool.Run(numWorkers)
}
