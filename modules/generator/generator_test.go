package generator

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-mods/tagsvar/modules/config"
	"github.com/go-mods/tagsvar/modules/parser"
)

const testdata = "../../.testdata"

// copyToTemp copies a source file of the test data into a temporary directory,
// so that the generated file is written there and not next to the reference files
func copyToTemp(t *testing.T, name string) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(testdata, name))
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestGenerator_generateFile(t *testing.T) {
	tests := []struct {
		source    string
		reference string // empty when no file must be generated
	}{
		{source: "user.go", reference: "user.vars.go"},
		{source: "blog_author.go", reference: "blog_author.vars.go"},
		{source: "exclude_struct.go"},
	}
	for _, tt := range tests {
		t.Run(tt.source, func(t *testing.T) {
			source := copyToTemp(t, tt.source)

			parsed, err := parser.NewParser().ParseFile(source)
			if err != nil {
				t.Fatalf("ParseFile() error = %v", err)
			}
			if tt.reference == "" {
				if parsed != nil {
					t.Fatalf("ParseFile() = %v, want nil", parsed)
				}
				return
			}
			if parsed == nil {
				t.Fatal("ParseFile() = nil, want a file to generate")
			}

			if err := NewGenerator().generateFile(parsed); err != nil {
				t.Fatalf("generateFile() error = %v", err)
			}

			got, err := os.ReadFile(generatedPath(source))
			if err != nil {
				t.Fatal(err)
			}
			want, err := os.ReadFile(filepath.Join(testdata, tt.reference))
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, want) {
				t.Errorf("generated file differs from %s:\n%s", tt.reference, got)
			}
		})
	}
}

func TestGeneratedPath(t *testing.T) {
	prefix, suffix := config.C.Prefix, config.C.Suffix
	t.Cleanup(func() { config.C.Prefix, config.C.Suffix = prefix, suffix })

	dir := filepath.Join("some", "dir")
	tests := []struct {
		name, prefix, suffix, want string
	}{
		{name: "default configuration", prefix: "", suffix: ".vars", want: filepath.Join(dir, "user.vars.go")},
		{name: "prefix applies to the file name", prefix: "gen_", suffix: ".vars", want: filepath.Join(dir, "gen_user.vars.go")},
		{name: "custom suffix", prefix: "", suffix: "_tags", want: filepath.Join(dir, "user_tags.go")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.C.Prefix, config.C.Suffix = tt.prefix, tt.suffix
			if got := generatedPath(filepath.Join(dir, "user.go")); got != tt.want {
				t.Errorf("generatedPath() = %q, want %q", got, tt.want)
			}
		})
	}
}
