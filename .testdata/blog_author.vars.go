// Code generated by tagsvar. DO NOT EDIT.

package testdata

// File: ../../.testdata/blog_author.go

// Struct: Author
// Author is a struct that represents an author
const (
	// Tag: json
	JsonAuthorId    = "id"
	JsonAuthorName  = "name"
	JsonAuthorEmail = "email"

	// Tag: xml
	XmlAuthorId    = "id"
	XmlAuthorName  = "name"
	XmlAuthorEmail = "email"

	// Tag: gorm
	GormAuthorId    = "id"
	GormAuthorName  = "name"
	GormAuthorEmail = "email"
)

var (
	// Tag: json

	// Tag: xml

	// Tag: gorm
	GormAuthorIdOptions = map[string]any{
		"type":        "uuid",
		"default":     "uuid_generate_v4()",
		"primary_key": nil,
	}
	GormAuthorNameOptions = map[string]any{
		"type":     "varchar(255)",
		"not null": nil,
	}
)

// Struct: Blog
// Blog is a struct that represents an author blog
const (
	// Tag: json
	JsonBlogId     = "id"
	JsonBlogAuthor = "author"
	JsonBlogUpvote = "upvote"

	// Tag: gorm
	GormBlogId     = "id"
	GormBlogAuthor = "author"
	GormBlogUpvote = "upvote"
)

var (
	// Tag: json

	// Tag: gorm
	GormBlogAuthorOptions = map[string]any{
		"embedded": nil,
	}
)
