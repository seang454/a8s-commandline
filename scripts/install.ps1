param(
    [string]$BinaryPath = "",
    [string]$InstallDir = "$env:LOCALAPPDATA\a8s\bin"
)

$ErrorActionPreference = "Stop"

if ($BinaryPath -eq "") {
    $BinaryPath = Join-Path (Get-Location) "dist\a8s-windows-amd64.exe"
}

if (-not (Test-Path $BinaryPath)) {
    throw "Binary not found: $BinaryPath"
}

New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
$target = Join-Path $InstallDir "a8s.exe"
Copy-Item -Force -Path $BinaryPath -Destination $target

Write-Host "Installed a8s to $target"
Write-Host "Add this directory to PATH if it is not already present:"
Write-Host "  $InstallDir"
