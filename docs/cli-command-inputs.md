**CLI Command Input Recommendations**

- **Context:** I could not access the two paths you provided from this workspace. I used the generated catalog at `docs/cli-command-catalog.md` (created from `./dist/a8s list all -o json`) to infer which commands should accept an operation manifest (`--file`, YAML/JSON) versus being flags-only.

- **Goal:** Produce a concise mapping of which CLI commands prefer `--file` (operation manifest) and which are flags-only. Use the extraction snippets below to regenerate or to produce per-repo lists locally.

Commands that typically prefer an operation manifest (`--file`, YAML/JSON):

- Examples: `deploy`, `create`, `update`, `apply`, `import`, `set`, `restore`, `upload`, `download`, `redeploy`, `webhook create`, `context create`, `manifest init`, `manifest validate`, and other commands that mutate server state or accept complex nested inputs.

Commands that are usually flags-only:

- Examples: `list`, `get`, `status`, `health`, `help`, `version`, `logs`, `status`, `catalog`, `request` and other read/inspection commands.

Exact extraction (PowerShell):

```powershell
$catalog = 'docs/cli-command-catalog.md'
# Commands that prefer a file (flags, file)
Get-Content $catalog |
  Where-Object { $_ -match 'flags, file\s*\|' -and $_ -notmatch '^\|\s*Command\s*\|' } |
  ForEach-Object { ($_ -split '\|')[1].Trim() } |
  Set-Content docs/commands-prefer-file.txt -Encoding UTF8

# Commands that are flags-only
Get-Content $catalog |
  Where-Object { $_ -match '\|\s*flags\s*\|' -and $_ -notmatch '^\|\s*Command\s*\|' } |
  ForEach-Object { ($_ -split '\|')[1].Trim() } |
  Set-Content docs/commands-prefer-flags.txt -Encoding UTF8
```

If you want, I can run these commands here and commit `docs/commands-prefer-file.txt` and `docs/commands-prefer-flags.txt` into the repo — or I can re-run the same extraction specifically for the two folders you mentioned if you add them to the workspace. Which would you prefer?
