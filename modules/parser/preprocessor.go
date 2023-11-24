package parser

import (
	"strings"
)

// Preprocessor is the preprocessor
//
// The preprocessor parse the string and extract the options
// This options will be used to know if a struct should
// be included or not and to know which tags to include
//
//	#tagsvar:all -> include the struct and all the fields with all the tags (equivalent to #tagsvar)
//	#tagsvar:include -> include the struct and all the fields (equivalent to #tagsvar)
//	#tagsvar:include:all -> include the struct and all the fields with all the tags (equivalent to #tagsvar)
//	#tagsvar:include:json -> include the struct and all the fields with the json tag
//	#tagsvar:include:json,xml -> include the struct and all the fields with the json and xml tags
//
//	#tagsvar:exclude -> exclude the struct and all the fields
//	#tagsvar:exclude:all -> exclude the struct and all the fields with all the tags (equivalent to #tagsvar:exclude)
//	#tagsvar:exclude:json -> only exclude the fields with the json tag
//	#tagsvar:exclude:json,xml" -> only exclude the fields with the json and xml tags
//
// Exclude is the default behavior, so if the preprocessor is not found in the comment, the struct will be excluded
// If the preprocessor is found and contains the include and exclude keywords, the exclude keyword will have the priority
type Preprocessor struct {
	// This is the preprocessor string name
	// The default value is #tagsvar
	preprocessor string

	// Include is a boolean that indicates if the struct should be included or not
	// The default value is true
	Include bool

	// IncludeTags is a list of tags that should be included
	// The default value is nil
	IncludeTags []string

	// Exclude is a boolean that indicates if the struct should be excluded or not
	// The default value is false
	Exclude bool

	// ExcludeTags is a list of tags that should be excluded
	// The default value is nil
	ExcludeTags []string
}

// Option is a simple key, value pair
type Option struct {
	Key   string
	Value string
}

// NewPreprocessor returns a new preprocessor
func NewPreprocessor() *Preprocessor {
	return &Preprocessor{
		preprocessor: "#tagsvar",
		Include:      true,
		IncludeTags:  nil,
		Exclude:      false,
		ExcludeTags:  nil,
	}
}

// Parse parses the string and extract the options
func (p *Preprocessor) Parse(s string) {
	// Clean the string
	s = strings.TrimSpace(s)
	// Check if the string starts with the preprocessor
	if !strings.HasPrefix(s, p.preprocessor) {
		return
	}

	// Split the string by :
	// The first element is the preprocessor
	// The second element is the include or exclude keyword
	// The third element is the tags
	split := strings.Split(s, ":")

	// If the length of the split is 1, it means that the preprocessor is the only element
	// So the struct should be included and all the tags should be included
	if len(split) == 1 {
		p.Include = true
		p.IncludeTags = nil
		p.Exclude = false
		p.ExcludeTags = nil
		return
	}

	// If the length of the split is 2, it means that the preprocessor and the keyword are the only elements
	// So the struct should be included or excluded and all the tags should be included
	if len(split) == 2 {
		// Check if the keyword is include
		if split[1] == "include" {
			p.Include = true
			p.IncludeTags = nil
			p.Exclude = false
			p.ExcludeTags = nil
			return
		}
		// Check if the keyword is exclude
		if split[1] == "exclude" {
			p.Include = false
			p.IncludeTags = nil
			p.Exclude = true
			p.ExcludeTags = nil
			return
		}
		return
	}

	// If the length of the split is 3 or more, it means that the preprocessor, the keyword and the tags are the only elements
	// So the struct should be included or excluded and the tags should be included or excluded
	if len(split) >= 3 {
		// Check if the keyword is include
		if split[1] == "include" {
			p.Include = true
			if split[2] == "all" {
			} else {
				p.IncludeTags = strings.Split(split[2], ",")
			}
			p.Exclude = false
			p.ExcludeTags = nil
			return
		}
		// Check if the keyword is exclude
		if split[1] == "exclude" {
			if split[2] == "all" {
				p.Include = false
				p.IncludeTags = nil
				p.Exclude = true
				p.ExcludeTags = nil
			} else {
				p.Include = true
				p.IncludeTags = nil
				p.Exclude = false
				p.ExcludeTags = strings.Split(split[2], ",")
			}
			return
		}
		return
	}

	return
}

// DoProcess returns true if the struct should be included and false if the struct should be excluded
func (p *Preprocessor) DoProcess() bool {
	// If the struct should be excluded, return false
	if p.Exclude {
		return false
	}
	// If the struct should be included, return true
	if p.Include {
		return true
	}
	// If the struct should not be included and not be excluded, return false
	return false
}

// ShouldProcess returns true if the tag should be included and false if the tag should be excluded
func (p *Preprocessor) ShouldProcess(tag string) bool {
	// If the struct should be excluded, return false
	if p.Exclude {
		return false
	}
	// If the struct should be included and the tags are in the exclude list, return false
	if p.Include && p.ExcludeTags != nil {
		for _, t := range p.ExcludeTags {
			if t == tag {
				return false
			}
		}
		return true
	}
	// If the struct should be included and the tags are nil, return true
	if p.Include && p.IncludeTags == nil {
		return true
	}
	// If the struct should be included and the tags are not nil, check if the tag is in the list
	if p.Include && p.IncludeTags != nil {
		for _, t := range p.IncludeTags {
			if t == tag {
				return true
			}
		}
		return false
	}
	// If the struct should not be included and not be excluded, return false
	return false
}
