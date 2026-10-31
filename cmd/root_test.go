package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-mods/tagsvar/modules/config"
)

const userSource = "package x\n\n// #tagsvar\ntype User struct {\n\tID int `json:\"id\"`\n}\n"

// run executes the root command with args and returns what it wrote on stdout and stderr
func run(t *testing.T, args ...string) (stdout, stderr string, err error) {
	t.Helper()
	verbose, silent := config.C.Verbose, config.C.Silent
	t.Cleanup(func() { config.C.Verbose, config.C.Silent = verbose, silent })

	var out, errOut bytes.Buffer
	root := newRootCmd()
	root.SetArgs(args)
	root.SetOut(&out)
	root.SetErr(&errOut)
	err = root.Execute()
	return out.String(), errOut.String(), err
}

// newProject creates a directory holding user.go, and sub/user.go in a subdirectory
func newProject(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	for _, path := range []string{filepath.Join(dir, "user.go"), filepath.Join(dir, "sub", "user.go")} {
		if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(userSource), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func TestGenAndClean(t *testing.T) {
	tests := []struct {
		name    string
		flags   []string
		wantSub bool
	}{
		{name: "directory only", flags: nil, wantSub: false},
		{name: "recursive", flags: []string{"-r"}, wantSub: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := newProject(t)
			root, sub := filepath.Join(dir, "user.vars.go"), filepath.Join(dir, "sub", "user.vars.go")

			stdout, stderr, err := run(t, append([]string{"gen", "-d", dir}, tt.flags...)...)
			if err != nil {
				t.Fatalf("gen error = %v", err)
			}
			if stdout != "" {
				t.Errorf("gen wrote %q on stdout, want nothing: messages go to stderr", stdout)
			}
			if !strings.Contains(stderr, "Parsing files in") {
				t.Errorf("gen stderr = %q, want the progress messages", stderr)
			}
			if !exists(root) || exists(sub) != tt.wantSub {
				t.Fatalf("after gen: %s exists = %v, %s exists = %v (want true, %v)",
					root, exists(root), sub, exists(sub), tt.wantSub)
			}

			if _, _, err := run(t, append([]string{"clean", "-d", dir}, tt.flags...)...); err != nil {
				t.Fatalf("clean error = %v", err)
			}
			if exists(root) || exists(sub) {
				t.Errorf("after clean: %s exists = %v, %s exists = %v (want false, false)",
					root, exists(root), sub, exists(sub))
			}
			if !exists(filepath.Join(dir, "user.go")) {
				t.Error("clean deleted the source file")
			}
		})
	}
}

func TestSilent(t *testing.T) {
	dir := newProject(t)
	stdout, stderr, err := run(t, "gen", "-s", "-d", dir)
	if err != nil {
		t.Fatalf("gen error = %v", err)
	}
	if stdout != "" || stderr != "" {
		t.Errorf("gen -s wrote stdout = %q, stderr = %q, want nothing", stdout, stderr)
	}
}

func TestErrors(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing")
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "gen with an argument", args: []string{"gen", "extra"}, want: "unknown command"},
		{name: "clean with an argument", args: []string{"clean", "extra"}, want: "unknown command"},
		{name: "missing directory", args: []string{"gen", "-d", missing}, want: "get working directory"},
		{name: "missing directory in silent mode", args: []string{"clean", "-s", "-d", missing}, want: "get working directory"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, stderr, err := run(t, tt.args...)
			if err == nil {
				t.Fatal("error = nil, want an error")
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("error = %q, want it to contain %q", err, tt.want)
			}
			// The error is returned to main, which prints it once: no usage, no duplicate
			if strings.Contains(stderr, "Usage:") || strings.Contains(stderr, err.Error()) {
				t.Errorf("stderr = %q, want neither the usage nor the error", stderr)
			}
		})
	}
}
