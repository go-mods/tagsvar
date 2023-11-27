package parser

import (
	"github.com/go-mods/tags"
	"testing"
)

func TestParser_parseFile(t *testing.T) {
	var structTests = []struct {
		filename string
		parsed   *File
	}{
		{
			filename: "../../.testdata/user.go",
			parsed: &File{
				Path:    "../../.testdata/user.go",
				Package: "testdata",
				Structs: []Struct{
					{
						Name:    "User",
						Comment: "User is a struct that represents a user",
						Fields: []Field{
							{
								Name:    "ID",
								Comment: "",
								Type:    "int",
								Tags: []tags.Tag{
									{Key: "json", Name: "id"},
									{Key: "xml", Name: "id"},
									{Key: "gorm", Name: "id"},
								},
							},
							{
								Name:    "Name",
								Comment: "",
								Type:    "string",
								Tags: []tags.Tag{
									{Key: "json", Name: "name"},
									{Key: "xml", Name: "name"},
									{Key: "gorm", Name: "name"},
								},
							},
						},
					},
				},
			},
		},
		{
			filename: "../../.testdata/blog_author.go",
			parsed: &File{
				Path:    "../../.testdata/blog_author.go",
				Package: "testdata",
				Structs: []Struct{
					{
						Name:    "Author",
						Comment: "Author is a struct that represents an author",
						Fields: []Field{
							{
								Name:    "ID",
								Comment: "",
								Type:    "int",
								Tags: []tags.Tag{
									{Key: "json", Name: "id"},
									{Key: "xml", Name: "id"},
									{Key: "gorm", Name: "id",
										Options: []*tags.Option{
											{Key: "type", Value: "uuid"},
											{Key: "default", Value: "uuid_generate_v4()"},
											{Key: "primary_key", Value: nil},
										},
									},
								},
							},
							{
								Name:    "Name",
								Comment: "",
								Type:    "string",
								Tags: []tags.Tag{
									{Key: "json", Name: "name"},
									{Key: "xml", Name: "name"},
									{Key: "gorm", Name: "name",
										Options: []*tags.Option{
											{Key: "type", Value: "varchar(255)"},
											{Key: "not null", Value: nil},
										},
									},
								},
							},
							{
								Name:    "Email",
								Comment: "",
								Type:    "string",
								Tags: []tags.Tag{
									{Key: "json", Name: "email"},
									{Key: "xml", Name: "email"},
									{Key: "gorm", Name: "email"},
								},
							},
						},
					},
					{
						Name:    "Blog",
						Comment: "Blog is a struct that represents an author blog",
						Fields: []Field{
							{
								Name:    "ID",
								Comment: "",
								Type:    "int",
								Tags: []tags.Tag{
									{Key: "json", Name: "id"},
									//{Key: "xml", Name: "id"},
									{Key: "gorm", Name: "id"},
								},
							},
							{
								Name:    "Author",
								Comment: "",
								Type:    "Author",
								Tags: []tags.Tag{
									{Key: "json", Name: "author"},
									//{Key: "xml", Name: "author"},
									{Key: "gorm", Name: "author", Options: []*tags.Option{{Key: "embedded"}}},
								},
							},
							{
								Name:    "Upvote",
								Comment: "",
								Type:    "int32",
								Tags: []tags.Tag{
									{Key: "json", Name: "upvote"},
									//{Key: "xml", Name: "upvote"},
									{Key: "gorm", Name: "upvote"},
								},
							},
						},
					},
				},
			},
		},
		{
			filename: "../../.testdata/exclude_struct.go",
			parsed:   nil,
		},
	}

	parser := NewParser()

	for _, test := range structTests {
		parsed, err := parser.ParseFile(test.filename)

		if err == nil && parsed != nil {
			if parsed == nil && test.parsed != nil {
				t.Errorf("Parse() got = %v, want %v", parsed, test.parsed)
			}
			if parsed != nil && test.parsed == nil {
				t.Errorf("Parse() got = %v, want %v", parsed, test.parsed)
			}
			if parsed != nil && test.parsed != nil && len(parsed.Structs) != len(test.parsed.Structs) {
				t.Errorf("Parse() got = %v, want %v", len(parsed.Structs), len(test.parsed.Structs))
			}
			for i, str := range parsed.Structs {
				if str.Name != test.parsed.Structs[i].Name {
					t.Errorf("Parse() got = %v, want %v", str.Name, test.parsed.Structs[i].Name)
				}
				if str.Comment != test.parsed.Structs[i].Comment {
					t.Errorf("Parse() got = %v, want %v", str.Comment, test.parsed.Structs[i].Comment)
				}
				if len(str.Fields) != len(test.parsed.Structs[i].Fields) {
					t.Errorf("Parse() got = %v, want %v", len(str.Fields), len(test.parsed.Structs[i].Fields))
				}
				for j, field := range str.Fields {
					if field.Name != test.parsed.Structs[i].Fields[j].Name {
						t.Errorf("Parse() got = %v, want %v", field.Name, test.parsed.Structs[i].Fields[j].Name)
					}
					if field.Comment != test.parsed.Structs[i].Fields[j].Comment {
						t.Errorf("Parse() got = %v, want %v", field.Comment, test.parsed.Structs[i].Fields[j].Comment)
					}
					if field.Type != test.parsed.Structs[i].Fields[j].Type {
						t.Errorf("Parse() got = %v, want %v", field.Type, test.parsed.Structs[i].Fields[j].Type)
					}
					if len(field.Tags) != len(test.parsed.Structs[i].Fields[j].Tags) {
						t.Errorf("Parse() got = %v, want %v", len(field.Tags), len(test.parsed.Structs[i].Fields[j].Tags))
					}
					for k, tag := range field.Tags {
						//if tag.Tag != test.parsed.Structs[i].Fields[j].Tags[k].Tag {
						//	t.Errorf("Parse() got = %v, want %v", tag.Tag, test.parsed.Structs[i].Fields[j].Tags[k].Tag)
						//}
						if tag.Key != test.parsed.Structs[i].Fields[j].Tags[k].Key {
							t.Errorf("Parse() got = %v, want %v", tag.Key, test.parsed.Structs[i].Fields[j].Tags[k].Key)
						}
						//if tag.Value != test.parsed.Structs[i].Fields[j].Tags[k].Value {
						//	t.Errorf("Parse() got = %v, want %v", tag.Value, test.parsed.Structs[i].Fields[j].Tags[k].Value)
						//}
						if tag.Name != test.parsed.Structs[i].Fields[j].Tags[k].Name {
							t.Errorf("Parse() got = %v, want %v", tag.Name, test.parsed.Structs[i].Fields[j].Tags[k].Name)
						}
						if len(tag.Options) != len(test.parsed.Structs[i].Fields[j].Tags[k].Options) {
							t.Errorf("Parse() got = %v, want %v", len(tag.Options), len(test.parsed.Structs[i].Fields[j].Tags[k].Options))
						}
						for l, opt := range tag.Options {
							if opt.Key != test.parsed.Structs[i].Fields[j].Tags[k].Options[l].Key {
								t.Errorf("Parse() got = %v, want %v", opt.Key, test.parsed.Structs[i].Fields[j].Tags[k].Options[l].Key)
							}
							if opt.Value != test.parsed.Structs[i].Fields[j].Tags[k].Options[l].Value {
								t.Errorf("Parse() got = %v, want %v", opt.Value, test.parsed.Structs[i].Fields[j].Tags[k].Options[l].Value)
							}
						}
					}
				}
			}

		} else if err != nil {
			t.Errorf("Parse() error = %v", err)
		}

	}
}
