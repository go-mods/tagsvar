package parser

import (
	"github.com/go-mods/tags"
	"github.com/go-mods/tagsvar/modules/config"
	"github.com/go-mods/tagsvar/modules/fs"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

type Parser struct {

	// This is the tagsvar preprocessor
	// It is used by the parser to know if a struct should
	// be included or not and to know which tags to include
	//
	// The default value is #tagsvar
	preprocessor *Preprocessor

	// This is used to indicate if the file should be processed or not even if the preprocessor is not found
	forceProcess bool
}

// NewParser creates a new Parser
func NewParser() *Parser {
	return &Parser{
		preprocessor: NewPreprocessor(),
	}
}

// ParseDir parses a directory and returns a map of parsed File
// It extracts the package name, the structs, the fields and the tags from the files
// It will be used to generate the variables files
func (p *Parser) ParseDir(path string, recursive bool) (map[string]*File, error) {
	// List files to parse
	files, err := fs.ListFiles(path, recursive, config.C.IsProjectFile)
	if err != nil {
		return nil, err
	}
	// Slice of parsed files
	parsedFiles := make(map[string]*File)

	// Parse the files
	for _, filename := range files {
		parsedFile, err := p.ParseFile(filename)
		if err != nil {
			return nil, err
		}
		parsedFiles[filename] = parsedFile
	}
	return parsedFiles, nil
}

// ParseFile parses a file and returns a File
// It extracts the package name, the structs, the fields and the tags from the file
// It will be used to generate the variables files
func (p *Parser) ParseFile(filename string) (*File, error) {
	var err error

	// Read the file
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Parse the file and get the File
	parsedFile, err := p.parseFile(filename, content)
	if err != nil {
		return nil, err
	}
	return parsedFile, nil
}

func (p *Parser) parseFile(filename string, content []byte) (*File, error) {
	// Parse the file and get the AST
	fileSet := token.NewFileSet()
	astFile, err := parser.ParseFile(fileSet, filename, content, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// Create the parsed File
	parsedFile := &File{}
	parsedFile.Path = filename
	parsedFile.Package = astFile.Name.Name

	// Inspect the AST
	ast.Inspect(astFile, func(node ast.Node) bool {
		switch node := node.(type) {
		// Generic declaration node
		// (Const, Type, Var)
		case *ast.GenDecl:
			{
				// Get the comment
				comment, process := p.processComment(node.Doc.Text())
				if !process {
					return true
				}

				// Loop over the specs to get the structs
				for _, spec := range node.Specs {
					switch spec := spec.(type) {
					case *ast.TypeSpec:
						{
							switch spec.Type.(type) {
							case *ast.StructType:
								parsedStruct, parseErr := p.parseStruct(spec, comment)
								if parseErr != nil {
									err = parseErr
									return false
								}
								// Add the struct to the file
								if parsedStruct != nil && len(parsedStruct.Fields) > 0 {
									parsedFile.Structs = append(parsedFile.Structs, *parsedStruct)
								}
							}
						}
					}
				}
			}
		}
		return true
	})

	// If the File is empty, return nil
	if len(parsedFile.Structs) == 0 {
		return nil, nil
	}

	return parsedFile, nil
}

func (p *Parser) processComment(comment string) (string, bool) {
	// Split the comment lines
	lines := strings.Split(comment, "\n")

	// By default, the file will not be processed
	// unless the preprocessor is found in the comment
	process := false

	// Check if the comment is a preprocessor
	if p.preprocessor.preprocessor != "" && len(lines) > 0 {
		for i, line := range lines {
			if strings.HasPrefix(line, p.preprocessor.preprocessor) {
				process = true
				// Initialize the preprocessor tags
				p.preprocessor.Parse(line)
				// Remove the comment
				lines = append(lines[:i], lines[i+1:]...)
				break
			}
		}
	}

	// Clean the comment
	comment = strings.Join(lines, "\n")
	comment = strings.TrimFunc(comment, func(r rune) bool { return r == ' ' || r == '\r' || r == '\n' })

	//
	process = p.preprocessor.DoProcess()

	// Process all the files even if the preprocessor is not found
	if !process && p.forceProcess {
		process = true
	}

	return comment, process
}

func (p *Parser) parseStruct(typeSpec *ast.TypeSpec, comment string) (*Struct, error) {
	// Convert to *ast.StructType to check if it is a struct
	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return nil, nil
	}

	// Create the parsed Struct
	parsedStruct := &Struct{}
	parsedStruct.Name = typeSpec.Name.Name
	parsedStruct.Comment = comment

	// Iterate over the fields
	for _, field := range structType.Fields.List {
		// Get the comment
		comment := field.Doc.Text()
		// Iterate over the field names
		for _, fieldName := range field.Names {
			// Create the parsed Field
			parsedField := &Field{}
			parsedField.Name = fieldName.Name
			parsedField.Comment = comment
			parsedField.Type = p.parseType(field.Type)
			parsedField.Tags = p.parseTags(field.Tag)

			// Add the field to the struct
			parsedStruct.Fields = append(parsedStruct.Fields, *parsedField)
		}
	}

	return parsedStruct, nil
}

func (p *Parser) parseType(expr ast.Expr) string {
	switch expr := expr.(type) {
	case *ast.Ident:
		return expr.Name
	case *ast.StarExpr:
		return "*" + p.parseType(expr.X)
	case *ast.SelectorExpr:
		return p.parseType(expr.X) + "." + expr.Sel.Name
	case *ast.ArrayType:
		return "[]" + p.parseType(expr.Elt)
	case *ast.MapType:
		return "map[" + p.parseType(expr.Key) + "]" + p.parseType(expr.Value)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.ChanType:
		return "chan " + p.parseType(expr.Value)
	case *ast.FuncType:
		return "func"
	case *ast.StructType:
		return "struct"
	default:
		return ""
	}
}

func (p *Parser) parseTags(tag *ast.BasicLit) []tags.Tag {
	if tag == nil {
		return nil
	}

	// Get the tags value from *ast.Field.Tag
	v := tag.Value

	// Trim backticks from the tags value
	v = strings.Trim(v, "`")

	// Parse the tags
	tagList, err := tags.Parse(v)
	if err != nil {
		return nil
	}
	// Convert the tags a slice of tags
	tagsSlice := make([]tags.Tag, 0, len(tagList))
	for _, t := range tagList {
		// check if the tag is in the tags map
		if p.preprocessor.ShouldProcess(t.Key) {
			tagsSlice = append(tagsSlice, *t)
		}
	}
	return tagsSlice
}
