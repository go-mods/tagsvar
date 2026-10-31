package fs

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestListFiles(t *testing.T) {
	dir := t.TempDir()
	root := filepath.Join(dir, "root.go")
	sub := filepath.Join(dir, "sub", "sub.go")
	other := filepath.Join(dir, "notes.txt")
	for _, path := range []string{root, sub, other} {
		if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("package x\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	tests := []struct {
		name      string
		recursive bool
		want      []string
	}{
		{name: "root directory only", recursive: false, want: []string{root}},
		{name: "recursive", recursive: true, want: []string{root, sub}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListFiles(dir, tt.recursive, IsGoFile)
			if err != nil {
				t.Fatalf("ListFiles() error = %v", err)
			}
			slices.Sort(got)
			if !slices.Equal(got, tt.want) {
				t.Errorf("ListFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}
