# gGrep Progress Report

## What's Done

### Phase 1: Core Engine (MVP)
* **Standard Library:** Built the core search engine entirely with standard Go libraries to avoid bloat.
* **Memory Safety:** Implemented `bufio.Scanner` for line-by-line file reading, ensuring minimal memory footprint even for large files.
* **Concurrency:** Developed a custom worker pool using Goroutines and Channels. Files are scanned concurrently, with outputs synchronized to avoid race conditions and terminal mangling.
* **Testing:** Wrote and passed comprehensive, idiomatic table-driven tests (`search_test.go`, `scanner_test.go`, `worker_test.go`) to verify concurrency and matching logic.

### Phase 2: Context Flags
* **Implementation:** Added support for trailing (`-A`), leading (`-B`), and output (`-C`) context flags.
	* **Ring Buffer:** Engineered a Circular (Ring) Buffer for the `-B` flag to maintain a strict memory limit, only storing the required preceding lines.

### Phase 3: Interactive Mode
* **TUI Implementation:** Developed a dynamic terminal user interface using `charmbracelet/bubbletea` and `bubbles` for real-time search filtering.
* **Debouncing & Performance:** Engineered a 300ms keystroke debouncer to prevent excessive regex compilations and worker pool thrashing.
* **Graceful Degradation:** Safely handles intermediate invalid regex states without crashing, retaining the last valid results on the screen.
* **UNIX Pipeline Compatibility:** Rendered the TUI entirely to `os.Stderr` while explicitly pushing matched lines to `os.Stdout` upon exit, allowing seamless piping (e.g., `./bin/ggrep -I file.txt > output.txt`).
* **Asynchronous Integration:** Integrated the existing asynchronous `core.WorkerPool` into the Bubbletea Update loop via a channel proxy, ensuring smooth, non-blocking UI rendering.

## How to Run & Test

**1. Build the Binary**
```bash
go build -o bin/ggrep ./cmd/ggrep
```

**2. Run Automated Tests**
```bash
go test -v ./...
```

**3. Manual Testing Examples**

*Basic Match:*
```bash
./bin/ggrep "func" ./internal/ggrep/core/*.go
```

*Regex Match (-E) with Case Insensitivity (-i):*
```bash
./bin/ggrep -E -i "^func" ./internal/ggrep/core/*.go
```

*Context Matching:*
(Tip: Use a dummy `numbers.txt` file stretching from 1 to 10 to test context calculations easily).
```bash
./bin/ggrep -C 3 "5" numbers.txt
```

*Interactive Mode (-I):*
```bash
./bin/ggrep -I sample.txt
```

