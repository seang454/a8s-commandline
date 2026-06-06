package output

import (
	"encoding/json"
	"fmt"
	"os"
)

func PrintJSON(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		fmt.Fprintln(os.Stderr, "Error encoding JSON:", err)
	}
}
