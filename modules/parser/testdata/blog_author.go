//go:build exclude

package testdata

// Author is a struct that represents an author
// #tagsvar
type Author struct {
	ID    int    `json:"id"    xml:"id"    gorm:"id"`
	Name  string `json:"name" xml:"name" gorm:"name" `
	Email string `json:"email" xml:"email" gorm:"email"`
}

// Blog is a struct that represents an author blog
// #tagsvar:exclude:xml
type Blog struct {
	ID     int    `json:"id"    xml:"id"    gorm:"id"`
	Author Author `json:"author" xml:"author" gorm:"author,embedded"`
	Upvote int32  `json:"upvote" xml:"upvote" gorm:"upvote"`
}
