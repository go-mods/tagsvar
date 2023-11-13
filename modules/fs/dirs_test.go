package fs

import (
	"os"
	"path/filepath"
	"testing"
)

type DirTest struct {
	path, result string
}

var cwd string
var dirtests []DirTest

func init() {
	var err error

	// Get the current directory
	cwd, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	// Create the tests
	dirtests = []DirTest{
		{"", cwd},
		{".", cwd},
		{"..", filepath.Join(cwd, "..")},
		{"../fs", filepath.Join(cwd, "../fs")},
	}
}

func TestNilPath(t *testing.T) {
	result, _ := WorkDir()
	if result != cwd {
		t.Errorf("WorkDir() = %q, want %q", result, cwd)
	}
}

func TestWorkDir(t *testing.T) {
	for _, test := range dirtests {
		result, _ := WorkDir(test.path)
		if result != test.result {
			t.Errorf("WorkDir(%q) = %q, want %q", test.path, result, test.result)
		}
	}
}
