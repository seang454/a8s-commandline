package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/yourname/a8s/internal/cli"
)

func main() {
	root := cli.NewRootCommand(strings.NewReader(""), ioDiscard{}, ioDiscard{})
	root.DisableAutoGenTag = true
	disableAutoGen(root)

	outputDir := filepath.Join("docs", "commands")
	if err := os.RemoveAll(outputDir); err != nil {
		fatal(err)
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		fatal(err)
	}
	var commands []*cobra.Command
	walk(root, &commands)
	for _, command := range commands {
		if command.Hidden {
			continue
		}
		if err := writeCommandPage(command, outputDir); err != nil {
			fatal(err)
		}
	}
	if err := writeReference(commands); err != nil {
		fatal(err)
	}
}

func disableAutoGen(command *cobra.Command) {
	command.DisableAutoGenTag = true
	for _, child := range command.Commands() {
		disableAutoGen(child)
	}
}

func writeReference(commands []*cobra.Command) error {
	sort.Slice(commands, func(i, j int) bool {
		return commands[i].CommandPath() < commands[j].CommandPath()
	})
	var buffer bytes.Buffer
	buffer.WriteString("# A8S CLI Command Reference\n\n")
	buffer.WriteString("Generated from the Cobra command tree. Do not manually edit command behavior here; update the command implementation and regenerate this file.\n\n")
	buffer.WriteString("Regenerate with:\n\n")
	buffer.WriteString("```bash\nmake generate-docs\n```\n\n")
	buffer.WriteString("## Command Pages\n\n")
	for _, command := range commands {
		if command.Hidden {
			continue
		}
		name := strings.ReplaceAll(command.CommandPath(), " ", "_") + ".md"
		buffer.WriteString(fmt.Sprintf("- [%s](commands/%s) - %s\n", command.CommandPath(), name, command.Short))
	}
	buffer.WriteString("\n## Generation Policy\n\n")
	buffer.WriteString("- Command pages are generated under `docs/commands/`.\n")
	buffer.WriteString("- Endpoint coverage remains tracked in `backend-api-cli-catalog.md`.\n")
	buffer.WriteString("- Commands with configurable mutation payloads should support both `--file` and equivalent flags.\n")
	buffer.WriteString("- Asynchronous operations should expose `--wait` when the backend provides a status URL, operation ID, or known polling endpoint.\n")
	return os.WriteFile(filepath.Join("docs", "command-reference.md"), buffer.Bytes(), 0o644)
}

func writeCommandPage(command *cobra.Command, outputDir string) error {
	var buffer bytes.Buffer
	buffer.WriteString("# " + command.CommandPath() + "\n\n")
	if command.Short != "" {
		buffer.WriteString(command.Short + "\n\n")
	}
	if command.Long != "" && command.Long != command.Short {
		buffer.WriteString("## Description\n\n")
		buffer.WriteString(strings.TrimSpace(command.Long) + "\n\n")
	}
	buffer.WriteString("## Usage\n\n")
	buffer.WriteString("```text\n" + command.UseLine() + "\n```\n\n")
	if command.Example != "" {
		buffer.WriteString("## Examples\n\n")
		buffer.WriteString("```bash\n" + strings.TrimSpace(command.Example) + "\n```\n\n")
	}
	writeFlags(&buffer, "Flags", command.NonInheritedFlags())
	writeFlags(&buffer, "Inherited Flags", command.InheritedFlags())
	var children []string
	for _, child := range command.Commands() {
		if child.Hidden {
			continue
		}
		children = append(children, fmt.Sprintf("- [%s](%s) - %s", child.CommandPath(), fileName(child), child.Short))
	}
	if len(children) > 0 {
		sort.Strings(children)
		buffer.WriteString("## Child Commands\n\n")
		for _, child := range children {
			buffer.WriteString(child + "\n")
		}
		buffer.WriteString("\n")
	}
	for _, key := range []string{"a8s.io/method", "a8s.io/endpoint", "a8s.io/controller"} {
		if value := command.Annotations[key]; value != "" {
			if key == "a8s.io/method" {
				buffer.WriteString("## Backend Endpoint\n\n")
			}
			buffer.WriteString("- `" + strings.TrimPrefix(key, "a8s.io/") + "`: `" + value + "`\n")
		}
	}
	if command.Annotations["a8s.io/method"] != "" {
		buffer.WriteString("\n")
	}
	return os.WriteFile(filepath.Join(outputDir, fileName(command)), buffer.Bytes(), 0o644)
}

func writeFlags(buffer *bytes.Buffer, title string, flags *pflag.FlagSet) {
	rows := []string{}
	flags.VisitAll(func(flag *pflag.Flag) {
		if flag.Hidden {
			return
		}
		name := "--" + flag.Name
		if flag.Shorthand != "" {
			name = "-" + flag.Shorthand + ", " + name
		}
		if flag.DefValue != "" && flag.DefValue != "false" {
			rows = append(rows, fmt.Sprintf("- `%s` `%s` - %s (default `%s`)", name, flag.Value.Type(), flag.Usage, flag.DefValue))
			return
		}
		rows = append(rows, fmt.Sprintf("- `%s` `%s` - %s", name, flag.Value.Type(), flag.Usage))
	})
	if len(rows) == 0 {
		return
	}
	sort.Strings(rows)
	buffer.WriteString("## " + title + "\n\n")
	for _, row := range rows {
		buffer.WriteString(row + "\n")
	}
	buffer.WriteString("\n")
}

func fileName(command *cobra.Command) string {
	return strings.ReplaceAll(command.CommandPath(), " ", "_") + ".md"
}

func walk(command *cobra.Command, result *[]*cobra.Command) {
	*result = append(*result, command)
	for _, child := range command.Commands() {
		walk(child, result)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

type ioDiscard struct{}

func (ioDiscard) Write(p []byte) (int, error) {
	return len(p), nil
}
