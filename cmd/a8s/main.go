package main

import (
	"fmt"
	"os"

	"github.com/yourname/a8s/internal/cli"
	"github.com/yourname/a8s/internal/clierrors"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(clierrors.ExitCode(err))
	}
}
