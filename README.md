[![Review Assignment Due Date](https://classroom.github.com/assets/deadline-readme-button-22041afd0340ce965d47ae6ef1cefeee28c7c493a6346c4f15d667ab976d596c.svg)](https://classroom.github.com/a/YG0u44_d)

# GoLine

## Authors
Leo Farmerie, Josh Justice, and Harshvardhan Singh

GoLine currently includes the interactive `polo` browser work and the in-progress `ggrep` work

## Polo Build And Run

`polo` uses a wrapper so it can update your shell to the directory you choose when you finish browsing

### Build The Binary

From the repo root

```powershell
cd cmd/polo
go build -o polo
```

### Run In Bash

Use `source` so the selected directory can carry back into your current shell

```bash
cd cmd/polo
source ./polo.sh
```

If you run `./polo.sh` directly, the browser will open, but the final `cd` will stay inside the subshell

### Run In PowerShell

From the `cmd/polo` directory

```powershell
.\polo.ps1
```

If your execution policy blocks the script, you can allow it for the current session with

```powershell
Set-ExecutionPolicy -Scope Process Bypass
```

### Optional Start Directory

You can also start `polo` in a specific directory

```powershell
go run . ..\..
```

## Ggrep Build And Run

### Build Instructions

From the repo root, navigate to the `ggrep` directory and build the binary:

```bash
cd cmd/ggrep
go build -o ggrep
```

### Standard Search

`ggrep` utilizes a highly concurrent worker pool along with `bufio.Scanner` to enable efficient, low-memory streaming and line-by-line file reading.

```bash
./ggrep -i "pattern" file.txt
```

### Context Flags

The tool supports `-A` (after), `-B` (before), and `-C` (context) flags to print lines surrounding your matches. It efficiently maps leading context using a memory-safe Circular Ring Buffer.

```bash
./ggrep -C 2 "pattern" file.txt
```

### Interactive Mode (TUI)

The `-I` flag launches the crown jewel of `ggrep`: an asynchronous, live terminal interface built with Bubbletea. It features a 300ms keystroke debouncer for optimal performance and handles graceful regex error degradation (preventing crashes if an incomplete regex is typed).

```bash
./ggrep -I file.txt
```

### POSIX Compliance

To maintain standard UNIX pipeline philosophy, the TUI renders entirely to `stderr`. This allows your standard output to remain clean so your final filtered data can be cleanly piped into other commands upon exiting.

```bash
./ggrep -I file.txt | wc -l
```

## Testing

From the repo root

```powershell
go test ./...
```
