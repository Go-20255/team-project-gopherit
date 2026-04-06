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
