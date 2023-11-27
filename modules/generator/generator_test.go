package generator

import (
	"github.com/go-mods/tagsvar/modules/parser"
	"testing"
)

func TestGenerator_generateFile(t *testing.T) {
	// List of files to parse
	toParse := []string{
		"../../.testdata/user.go",
		"../../.testdata/blog_author.go",
		"../../.testdata/exclude_struct.go",
	}

	// Create a parser
	p := parser.NewParser()

	// Create a generator
	g := NewGenerator()

	// Parse the files
	for _, filename := range toParse {
		// Parse the file
		parsed, err := p.ParseFile(filename)
		if err != nil {
			t.Errorf("ParseFile()")
			return
		}
		// Generate the file
		if parsed != nil {
			err = g.generateFile(parsed)
			if err != nil {
				t.Errorf("generateFile()")
				return
			}
		}
	}
}
