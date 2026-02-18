//nolint:testpackage // testscript requires package name to be testscript
package testscript

import (
	"context"
	"os"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"

	"github.com/mholtzscher/today/cmd"
	"github.com/mholtzscher/today/internal/output"
)

func TestMain(m *testing.M) {
	testscript.Main(m, map[string]func(){
		"today": func() {
			if err := cmd.Run(context.Background(), os.Args); err != nil {
				output.Stderrln(err.Error())
				os.Exit(1)
			}
		},
	})
}

func TestScript(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "scripts",
		Setup: func(env *testscript.Env) error {
			env.Setenv("HOME", env.Getenv("WORK"))
			return nil
		},
	})
}
