//go:build exclude

package testdata

// User is a struct that represents a user
// #tagsvar
type User struct {
	ID   int    `json:"id"    xml:"id"    gorm:"id"`
	Name string `json:"name"  xml:"name"  gorm:"name"`
}
