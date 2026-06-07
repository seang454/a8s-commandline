<#
.SYNOPSIS
Generate a markdown command index from the generated route registry.

This script parses internal/cli/catalog/generated_routes.go and emits
docs/command-index.md listing every implemented command, its HTTP method,
endpoint, whether it typically accepts a request body (naively based on the
HTTP method), and a link to the existing per-command doc in docs/commands
if present.

Run from the repository root:

  powershell -ExecutionPolicy Bypass -File .\scripts\generate-command-index.ps1

#>
Param()

$repoRoot = Split-Path -Parent $MyInvocation.MyCommand.Path -Resolve
Set-Location $repoRoot

$routesFile = Join-Path $repoRoot 'internal\cli\catalog\generated_routes.go'
$outFile = Join-Path $repoRoot 'docs\command-index.md'

if (-not (Test-Path $routesFile)) {
    Write-Error "Routes file not found: $routesFile"
    exit 2
}

$content = Get-Content -Raw -Path $routesFile

# Match each route entry. This is a tolerant regex and may be extended.
$pattern = '\{[^}]*Method:\s*"(?<method>[^"]+)"[^}]*Endpoint:\s*"(?<endpoint>[^"]+)"[^}]*Command:\s*\[]string\{(?<command>[^}]*)\}[^}]*Args:\s*\[]string\{(?<args>[^}]*)\}[^}]*\}'

$matches = [regex]::Matches($content, $pattern)

$lines = @()
$lines += '# A8S CLI Command Index'
$lines += ''
$lines += '> Generated from internal/cli/catalog/generated_routes.go — run scripts/generate-command-index.ps1 to refresh.'
$lines += ''
$lines += '| Command | Method | Endpoint | Request Body? | Doc |'
$lines += '|---|---|---|---:|---|'

foreach ($m in $matches) {
    $method = $m.Groups['method'].Value
    $endpoint = $m.Groups['endpoint'].Value
    $cmdRaw = $m.Groups['command'].Value.Trim()
    # command parts are quoted strings separated by commas
    $parts = @()
    if ($cmdRaw -ne '') {
        $rawParts = $cmdRaw -split ',' | ForEach-Object { $_.Trim() }
        foreach ($p in $rawParts) {
            if ($p -match '"(?<q>[^"]+)"') { $parts += $Matches['q'] }
        }
    }
    if ($parts.Count -eq 0) { continue }
    $cmd = 'a8s ' + ($parts -join ' ')
    $fileName = 'a8s_' + ($parts -join '_') + '.md'
    $docPath = Join-Path $repoRoot ('docs\commands\' + $fileName)
    $docLink = if (Test-Path $docPath) { "[doc](commands/$fileName)" } else { '' }
    $requestBody = if ($method -in @('GET','HEAD','OPTIONS','DELETE')) { 'No' } else { 'Usually' }
    $lines += "| $cmd | $method | $endpoint | $requestBody | $docLink |"
}

$lines | Out-File -FilePath $outFile -Encoding UTF8
Write-Output "Wrote $outFile"
