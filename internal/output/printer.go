package output

import (
	"encoding/json"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

type Printer struct {
	Out    io.Writer
	ErrOut io.Writer
	Format string
}

func (p Printer) Print(value any) error {
	switch p.Format {
	case "json":
		encoder := json.NewEncoder(p.Out)
		encoder.SetIndent("", "  ")
		return encoder.Encode(value)
	case "yaml":
		encoder := yaml.NewEncoder(p.Out)
		encoder.SetIndent(2)
		defer encoder.Close()
		return encoder.Encode(value)
	default:
		_, err := fmt.Fprintln(p.Out, formatTableFallback(value))
		return err
	}
}

func (p Printer) Progress(format string, args ...any) {
	fmt.Fprintf(p.ErrOut, format+"\n", args...)
}

func formatTableFallback(value any) string {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Sprint(value)
	}
	return string(data)
}
