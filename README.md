[![Review Assignment Due Date](https://classroom.github.com/assets/deadline-readme-button-22041afd0340ce965d47ae6ef1cefeee28c7c493a6346c4f15d667ab976d596c.svg)](https://classroom.github.com/a/YG0u44_d)

# GoLine

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

## Testing

From the repo root

```powershell
go test ./...
```
