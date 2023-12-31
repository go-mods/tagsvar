package generator

import (
	"bytes"
	"fmt"
	"github.com/go-mods/tags"
	"github.com/go-mods/tagsvar/modules/config"
	"github.com/go-mods/tagsvar/modules/parser"
	"github.com/rs/zerolog/log"
	"github.com/stoewer/go-strcase"
	"go/format"
	"os"
	"path/filepath"
	"strings"
)

type Generator struct {
}

// NewGenerator creates an instance of Generator
func NewGenerator() *Generator {
	return &Generator{}
}

// Generate generates the variables files
func (g *Generator) Generate(files map[parser.FilePath]*parser.File) error {
	// Loop through the files
	for _, file := range files {
		// Generate the file
		if file != nil {
			log.Info().Msgf("Generating file for %s", string(file.Path))
			err := g.generateFile(file)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// generateFile generates the variables file
func (g *Generator) generateFile(file *parser.File) error {

	// Construct the file path where the variables file will be generated
	filePath := string(file.Path)
	filePath = strings.TrimSuffix(filePath, filepath.Ext(filePath))
	filePath = config.C.Prefix + filePath + config.C.Suffix + ".go"

	// Generate the code to write to the file
	genCode, err := g.generateCode(file)
	if err != nil {
		return err
	}

	// Create the file
	genFile, err := os.Create(filepath.Clean(filePath))
	if err != nil {
		return err
	}

	// Write the code to the file
	_, err = genFile.Write(genCode)
	if err != nil {
		return err
	}

	return nil
}

// generateCode generates the code for the variables file
func (g *Generator) generateCode(file *parser.File) ([]byte, error) {
	if file == nil {
		return nil, nil
	}

	// Buffer to write the code to
	genCode := bytes.Buffer{}
	genCode.WriteString("// Code generated by tagsvar. DO NOT EDIT.\n\n")
	genCode.WriteString("package " + file.Package + "\n\n")
	genCode.WriteString("// File: " + string(file.Path) + "\n")

	// loop through file.Structs
	for _, s := range file.Structs {
		// Title
		g.writeTitle(&genCode, s)

		// Const
		g.writeConst(&genCode, s)

		// Var
		g.writeVars(&genCode, s)
	}

	return format.Source(genCode.Bytes())
}

// writeTitle writes the title of the struct
func (g *Generator) writeTitle(genCode *bytes.Buffer, s parser.Struct) {
	genCode.WriteString("\n")
	genCode.WriteString("// Struct: " + s.Name + "\n")
	if s.Comment != "" {
		genCode.WriteString("// " + s.Comment + "\n")
	}
}

// writeConst writes the const variable
func (g *Generator) writeConst(genCode *bytes.Buffer, s parser.Struct) {
	// at least one of the tag name must be set
	hasName := false
	for _, f := range s.Fields {
		for _, t := range f.Tags {
			if t.Name != "" {
				hasName = true
				break
			}
		}
	}

	// if no tag name is set, do not generate the const
	if !hasName {
		return
	}

	// Const's variables organised by tag key
	genCode.WriteString("const (\n")
	for _, tk := range s.TagKeys {
		genCode.WriteString("// Tag: " + tk + "\n")
		for _, f := range s.Fields {
			for _, tag := range f.Tags {
				if tag.Key == tk {
					constName := g.generateConstName(s, f, tag)
					if constName != "" {
						genCode.WriteString(constName + "\n")
					}
				}
			}
		}
		genCode.WriteString("\n")
	}
	genCode.WriteString(")\n")
}

// writeVars writes the vars variables from the tag options
func (g *Generator) writeVars(genCode *bytes.Buffer, s parser.Struct) {
	// at least one of the tag options must be set
	hasOptions := false
	for _, f := range s.Fields {
		for _, t := range f.Tags {
			if len(t.Options) > 0 {
				hasOptions = true
				break
			}
		}
	}

	// if no tag options are set, do not generate the vars
	if !hasOptions {
		return
	}

	// Var's variables organised by tag key
	genCode.WriteString("var (\n")
	for _, tk := range s.TagKeys {
		genCode.WriteString("// Tag: " + tk + "\n")
		for _, f := range s.Fields {
			for _, tag := range f.Tags {
				if tag.Key == tk {
					varOptions := g.generateVarOptions(s, f, tag)
					if varOptions != "" {
						genCode.WriteString(varOptions + "\n")
					}
				}
			}
		}
		genCode.WriteString("\n")
	}
	genCode.WriteString(")\n")
}

// generateConstName generates the variable from the tag name
func (g *Generator) generateConstName(s parser.Struct, f parser.Field, t tags.Tag) string {
	if t.Name == "" {
		return ""
	}
	return strcase.UpperCamelCase(t.Key) + strcase.UpperCamelCase(s.Name) + strcase.UpperCamelCase(f.Name) + ` = "` + t.Name + `"`
}

// generateVarOptions generates the variable from the tag options
func (g *Generator) generateVarOptions(s parser.Struct, f parser.Field, t tags.Tag) string {
	if len(t.Options) == 0 {
		return ""
	}

	options := strcase.UpperCamelCase(t.Key) + strcase.UpperCamelCase(s.Name) + strcase.UpperCamelCase(f.Name) + `Options = map[string]any{` + "\n"

	for _, o := range t.Options {
		if o.Value != nil {
			options = options + fmt.Sprintf(`"%s": "%v",`, o.Key, o.Value)
		} else {
			options = options + `"` + o.Key + `": nil, `
		}
		options = options + "\n"
	}

	options = options + `}`

	return options
}
