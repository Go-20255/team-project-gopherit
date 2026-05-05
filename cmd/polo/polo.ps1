$stateFile = Join-Path ([System.IO.Path]::GetTempPath()) "goline-polo-path.txt"
$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$binaryPath = Join-Path $scriptDir "polo.exe"

if (-not (Test-Path $binaryPath)) {
    $binaryPath = Join-Path $scriptDir "polo"
}

& $binaryPath @args

if (Test-Path $stateFile) {
    $target = (Get-Content $stateFile -Raw).Trim()
    if ($target.Length -gt 0) {
        Set-Location $target
    }

    Remove-Item $stateFile -ErrorAction SilentlyContinue
}
