//go:build exclude

package testdata

// ExcludeStruct is a struct which should be excluded
// #tagsvar:exclude
type ExcludeStruct struct {
	ID int `json:"id"    xml:"id"    gorm:"id"`
}

// AnotherExcludeStruct is a struct which should be excluded
// #tagsvar:exclude:all
type AnotherExcludeStruct struct {
	ID int `json:"id"    xml:"id"    gorm:"id"`
}
