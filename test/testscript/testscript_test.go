//nolint:testpackage // testscript requires package name to be testscript
package testscript

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestScript(t *testing.T) {
	cwd, _ := os.Getwd()
	projectRoot := filepath.Clean(filepath.Join(cwd, "..", ".."))
	homeDir, _ := os.UserHomeDir()

	testscript.Run(t, testscript.Params{
		Dir: "scripts",
		Setup: func(env *testscript.Env) error {
			env.Setenv("PROJECT_ROOT", projectRoot)
			env.Setenv("GOCACHE", filepath.Join(homeDir, ".cache", "go-build"))
			env.Setenv("GOMODCACHE", filepath.Join(homeDir, "go", "pkg", "mod"))
			env.Setenv("HOME", env.Getenv("WORK"))
			return nil
		},
		Cmds: map[string]func(ts *testscript.TestScript, neg bool, args []string){
			"today": func(ts *testscript.TestScript, _ bool, args []string) {
				root := ts.Getenv("PROJECT_ROOT")
				cmdArgs := append([]string{"run", "-C", root, "."}, args...)
				ts.Exec("go", cmdArgs...)
			},
		},
	})
}
