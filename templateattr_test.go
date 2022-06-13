package aform_test

import (
	"github.com/roleupjobboard/aform"
	"github.com/stretchr/testify/assert"
	"html/template"
	"testing"
)

func TestExportTmplAttrs_HTMLAttributes(t *testing.T) {
	tests := []struct {
		name string
		a    aform.ExportTmplAttrs
		want []template.HTMLAttr
	}{
		{
			name: "nil list",
			a:    nil,
			want: []template.HTMLAttr{},
		},
		{
			name: "empty list",
			a:    aform.ExportTmplAttrs{},
			want: []template.HTMLAttr{},
		},
		{
			name: "one key element",
			a:    aform.ExportTmplAttrs{"required": ""},
			want: []template.HTMLAttr{"required"},
		},
		{
			name: "one key/value element",
			a:    aform.ExportTmplAttrs{"id": "test"},
			want: []template.HTMLAttr{"id=\"test\""},
		},
		{
			name: "unknown key/value element",
			a:    aform.ExportTmplAttrs{"unknown": "doe"},
			want: []template.HTMLAttr{"unknown=\"doe\""},
		},
		{
			name: "ordered complete list",
			a:    aform.ExportTmplAttrs{"class": "", "id": "", "name": "", "data-": "", "src": "", "for": "", "type": "", "href": "", "value": "", "title": "", "alt": "", "role": "", "aria-": ""},
			want: []template.HTMLAttr{"class", "id", "name", "data-", "src", "for", "type", "href", "value", "title", "alt", "role", "aria-"},
		},
		{
			name: "unordered list",
			a:    aform.ExportTmplAttrs{"alt": "", "class": "", "name": "", "src": "", "for": "", "type": "", "data-": "", "href": "", "value": "", "title": "", "id": "", "role": "", "aria-": ""},
			want: []template.HTMLAttr{"class", "id", "name", "data-", "src", "for", "type", "href", "value", "title", "alt", "role", "aria-"},
		},
		{
			name: "ordered data and aria list",
			a:    aform.ExportTmplAttrs{"data-a": "", "data-b": "", "data-c": "", "aria-a": "", "aria-b": "", "aria-c": ""},
			want: []template.HTMLAttr{"data-a", "data-b", "data-c", "aria-a", "aria-b", "aria-c"},
		},
		{
			name: "unordered data and aria list",
			a:    aform.ExportTmplAttrs{"aria-a": "", "data-b": "", "data-c": "", "aria-c": "", "aria-b": "", "data-a": ""},
			want: []template.HTMLAttr{"data-a", "data-b", "data-c", "aria-a", "aria-b", "aria-c"},
		},
		{
			name: "bug one key/value aria-describedby element and one unknown element",
			a:    aform.ExportTmplAttrs{"aria-describedby": "err_0_id_remember_me", "required": ""},
			want: []template.HTMLAttr{"aria-describedby=\"err_0_id_remember_me\"", "required"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.a.HTMLAttributes(), "HTMLAttributes()")
		})
	}
}
