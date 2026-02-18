// today is a CLI tool.
package main

import (
	"context"
	"os"

	"github.com/mholtzscher/today/cmd"
	"github.com/mholtzscher/today/internal/output"
)

func main() {
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		output.Stderrln(err.Error())
		os.Exit(1)
	}
}
