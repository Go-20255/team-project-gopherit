package core

import (
	"sync"
)

// WorkerPool manages concurrent file scanning and synchronized output printing.
type WorkerPool struct {
	Scanner   *Scanner
	FilePaths []string
	Results   chan MatchResult
	PrintFn   func(MatchResult)
}

// NewWorkerPool initializes a new WorkerPool.
func NewWorkerPool(scanner *Scanner, files []string, printFn func(MatchResult)) *WorkerPool {
	return &WorkerPool{
		Scanner:   scanner,
		FilePaths: files,
		Results:   make(chan MatchResult, 100), // Buffered channel for performance
		PrintFn:   printFn,
	}
}

// Run executes the worker pool with the specified number of concurrent workers.
func (wp *WorkerPool) Run(numWorkers int) {
	jobs := make(chan string, len(wp.FilePaths))
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for filePath := range jobs {
				_ = wp.Scanner.ScanFile(filePath, wp.Results)
				// Ignoring errors for now (e.g. file not found, permission denied).
				// We can add a stderr logging channel in the future.
			}
		}()
	}

	// Feed jobs
	for _, f := range wp.FilePaths {
		jobs <- f
	}
	close(jobs)

	// Close results channel when all workers are done
	go func() {
		wg.Wait()
		close(wp.Results)
	}()

	// Single printer goroutine runs in the current thread (or we could wait on it)
	// Consuming from wp.Results ensures output lines are printed entirely, without race conditions.
	for res := range wp.Results {
		wp.PrintFn(res)
	}
}
