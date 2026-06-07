package manifestcmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	cliruntime "github.com/yourname/a8s/internal/cli/runtime"
	"github.com/yourname/a8s/internal/clierrors"
	manifestregistry "github.com/yourname/a8s/internal/operation/manifest"
)

func NewCommand(runtime *cliruntime.Runtime) *cobra.Command {
	command := &cobra.Command{
		Use:   "manifest",
		Short: "Generate and validate operation manifests",
	}
	command.AddCommand(newKinds(runtime), newSchema(runtime), newInit(runtime), newValidate(runtime))
	return command
}

func newKinds(runtime *cliruntime.Runtime) *cobra.Command {
	return &cobra.Command{
		Use:   "kinds",
		Short: "List supported operation manifest kinds",
		RunE: func(cmd *cobra.Command, args []string) error {
			definitions := manifestregistry.Definitions()
			rows := make([]map[string]any, 0, len(definitions))
			for _, definition := range definitions {
				rows = append(rows, map[string]any{
					"kind":        definition.Kind,
					"description": definition.Description,
					"strict":      definition.Strict != nil,
				})
			}
			return runtime.Printer.Print(rows)
		},
	}
}

func newSchema(runtime *cliruntime.Runtime) *cobra.Command {
	return &cobra.Command{
		Use:   "schema <kind>",
		Short: "Show the manifest schema summary for a kind",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			definition, ok := manifestregistry.Get(args[0])
			if !ok {
				return clierrors.Validation(fmt.Sprintf("unknown operation kind %q", args[0]))
			}
			return runtime.Printer.Print(definition)
		},
	}
}

func newInit(runtime *cliruntime.Runtime) *cobra.Command {
	var outputFile string
	var overwrite bool
	command := &cobra.Command{
		Use:   "init <kind>",
		Short: "Generate a starter manifest for a kind",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			template, err := manifestregistry.Template(args[0])
			if err != nil {
				return err
			}
			if outputFile == "" || outputFile == "-" {
				_, err := runtime.Out.Write([]byte(template))
				return err
			}
			if err := writeOutputFile(outputFile, []byte(template), overwrite); err != nil {
				return err
			}
			definition, _ := manifestregistry.Get(args[0])
			return runtime.Printer.Print(map[string]string{"kind": definition.Kind, "outputFile": outputFile})
		},
	}
	command.Flags().StringVar(&outputFile, "output-file", "", "write the starter manifest to a file instead of stdout")
	command.Flags().BoolVar(&overwrite, "overwrite", false, "replace output file if it already exists")
	return command
}

func newValidate(runtime *cliruntime.Runtime) *cobra.Command {
	var file string
	command := &cobra.Command{
		Use:   "validate",
		Short: "Validate an operation manifest without sending a backend request",
		RunE: func(cmd *cobra.Command, args []string) error {
			if file == "" {
				return clierrors.Validation("--file is required")
			}
			result, err := manifestregistry.ValidateFile(file, runtime.In)
			if err != nil {
				return err
			}
			return runtime.Printer.Print(result)
		},
	}
	command.Flags().StringVar(&file, "file", "", "YAML or JSON operation manifest; use - for stdin")
	return command
}

func writeOutputFile(path string, data []byte, overwrite bool) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	flags := os.O_WRONLY | os.O_CREATE
	if overwrite {
		flags |= os.O_TRUNC
	} else {
		flags |= os.O_EXCL
	}
	file, err := os.OpenFile(path, flags, 0o600)
	if errors.Is(err, os.ErrExist) {
		return clierrors.Validation(fmt.Sprintf("output file %q already exists; use --overwrite to replace it", path))
	}
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	defer file.Close()
	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}
	return nil
}
