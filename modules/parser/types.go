package parser

import "github.com/go-mods/tags"

// File represents a project file
// It contains the path of the file, the package name and the structs
// This information are extracted from the file and will be used to generate the variables files
type File struct {
	Path    string
	Package string
	Structs []Struct
}

// Struct represents a struct in a project file
// It contains the name of the struct and the fields
// This information are extracted from the file and will be used to generate the variables files
type Struct struct {
	Name    string
	Comment string
	Fields  []Field
}

// Field represents a field in a struct
// It contains the name of the field, the type and the tags
// This information are extracted from the file and will be used to generate the variables files
type Field struct {
	Name    string
	Comment string
	Type    string
	Tags    []tags.Tag
}
