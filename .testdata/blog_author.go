//go:build exclude

package testdata

// Author is a struct that represents an author
// #tagsvar
type Author struct {
	ID    int    `json:"id"    xml:"id"    gorm:"id;type:uuid;default:uuid_generate_v4();primary_key"`
	Name  string `json:"name"  xml:"name"  gorm:"name;type:varchar(255);not null"`
	Email string `json:"email" xml:"email" gorm:"email"`
}

// Blog is a struct that represents an author blog
// #tagsvar:exclude:xml
type Blog struct {
	ID     int    `json:"id"     xml:"id"     gorm:"id"`
	Author Author `json:"author" xml:"author" gorm:"author,embedded"`
	Upvote int32  `json:"upvote" xml:"upvote" gorm:"upvote"`
}
