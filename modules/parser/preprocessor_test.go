package parser

import (
	"testing"
)

var preprocessor = "#tagsvar"

func TestPreprocessor_Init(t *testing.T) {

	var tests = []struct {
		comment string
		want    *Preprocessor
	}{
		{
			comment: `#tagsvar`,
			want: &Preprocessor{
				preprocessor: preprocessor,
				Include:      true,
				IncludeTags:  nil,
				Exclude:      false,
				ExcludeTags:  nil,
			},
		},
		{
			comment: `#tagsvar:include`,
			want: &Preprocessor{
				preprocessor: preprocessor,
				Include:      true,
				IncludeTags:  nil,
				Exclude:      false,
				ExcludeTags:  nil,
			},
		},
		{
			comment: `#tagsvar:include:all`,
			want: &Preprocessor{
				preprocessor: preprocessor,
				Include:      true,
				IncludeTags:  nil,
				Exclude:      false,
				ExcludeTags:  nil,
			},
		},
		{
			comment: `#tagsvar:include:json`,
			want: &Preprocessor{
				preprocessor: preprocessor,
				Include:      true,
				IncludeTags:  []string{"json"},
				Exclude:      false,
				ExcludeTags:  nil,
			},
		},
		{
			comment: `#tagsvar:include:json,xml`,
			want: &Preprocessor{
				preprocessor: preprocessor,
				Include:      true,
				IncludeTags:  []string{"json", "xml"},
				Exclude:      false,
				ExcludeTags:  nil,
			},
		},
		{
			comment: `#tagsvar:exclude`,
			want: &Preprocessor{
				preprocessor: preprocessor,
				Include:      false,
				IncludeTags:  nil,
				Exclude:      true,
				ExcludeTags:  nil,
			},
		},
		{
			comment: `#tagsvar:exclude:all`,
			want: &Preprocessor{
				preprocessor: preprocessor,
				Include:      false,
				IncludeTags:  nil,
				Exclude:      true,
				ExcludeTags:  nil,
			},
		},
		{
			comment: `#tagsvar:exclude:json`,
			want: &Preprocessor{
				preprocessor: preprocessor,
				Include:      true,
				IncludeTags:  nil,
				Exclude:      false,
				ExcludeTags:  []string{"json"},
			},
		},
		{
			comment: `#tagsvar:exclude:json,xml`,
			want: &Preprocessor{
				preprocessor: preprocessor,
				Include:      true,
				IncludeTags:  nil,
				Exclude:      false,
				ExcludeTags:  []string{"json", "xml"},
			},
		},
	}

	// Prepare the preprocessor
	p := NewPreprocessor()

	for _, test := range tests {
		p.Parse(test.comment)
		t.Log(test.comment)

		if p.preprocessor != preprocessor {
			t.Errorf("Init() got = %v, want %v", p.preprocessor, preprocessor)
		}

		if p.Include != test.want.Include {
			t.Errorf("Init() got = %v, want %v", p.Include, test.want.Include)
		}
		if p.Exclude != test.want.Exclude {
			t.Errorf("Init() got = %v, want %v", p.Exclude, test.want.Exclude)
		}
		if len(p.IncludeTags) != len(test.want.IncludeTags) {
			t.Errorf("Init() got = %v, want %v", len(p.IncludeTags), len(test.want.IncludeTags))
		}
		if len(p.ExcludeTags) != len(test.want.ExcludeTags) {
			t.Errorf("Init() got = %v, want %v", len(p.ExcludeTags), len(test.want.ExcludeTags))
		}
	}
}

func TestPreprocessor_ParseResetsTags(t *testing.T) {
	tests := []struct {
		name     string
		previous string
		comment  string
	}{
		{name: "include all after include list", previous: "#tagsvar:include:json", comment: "#tagsvar:include:all"},
		{name: "include after include list", previous: "#tagsvar:include:json", comment: "#tagsvar:include"},
		{name: "directive after exclude list", previous: "#tagsvar:exclude:xml", comment: "#tagsvar"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPreprocessor()
			p.Parse(tt.previous)
			p.Parse(tt.comment)
			if p.IncludeTags != nil || p.ExcludeTags != nil {
				t.Errorf("Parse(%q) after %q: IncludeTags = %v, ExcludeTags = %v, want nil",
					tt.comment, tt.previous, p.IncludeTags, p.ExcludeTags)
			}
			for _, tag := range []string{"json", "xml", "gorm"} {
				if !p.ShouldProcess(tag) {
					t.Errorf("ShouldProcess(%q) = false, want true", tag)
				}
			}
		})
	}
}
