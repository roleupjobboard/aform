package aform_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBoolAttr_withString(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want aform.Attributable
	}{
		{
			name: "",
			args: args{
				name: "required",
			},
			want: aform.BoolAttr("required"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, aform.BoolAttr(tt.args.name), "BoolAttr(%v)", tt.args.name)
		})
	}
}

type testString string

func TestBoolAttr_withCustomType(t *testing.T) {
	type args struct {
		name testString
	}
	tests := []struct {
		name string
		args args
		want aform.Attributable
	}{
		{
			name: "",
			args: args{
				name: testString("disabled"),
			},
			want: aform.BoolAttr("disabled"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, aform.BoolAttr(tt.args.name), "BoolAttr(%v)", tt.args.name)
		})
	}
}

func TestAttr_withInt(t *testing.T) {
	type args struct {
		name  string
		value int
	}
	tests := []struct {
		name string
		args args
		want aform.ExportTmplAttrs
	}{
		{
			name: "with a positive",
			args: args{
				name:  "max",
				value: 3,
			},
			want: aform.ExportTmplAttrs{"max": "3"},
		},
		{
			name: "with a negative",
			args: args{
				name:  "minneg",
				value: -7,
			},
			want: aform.ExportTmplAttrs{"minneg": "-7"},
		},
		{
			name: "with 0",
			args: args{
				name:  "zero",
				value: 0,
			},
			want: aform.ExportTmplAttrs{"zero": "0"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, aform.ExportAttributableAttribute(aform.Attr(tt.args.name, tt.args.value)), "Attr(%v, %v)", tt.args.name, tt.args.value)
		})
	}
}

func TestAttr_withString(t *testing.T) {
	type args struct {
		name  string
		value string
	}
	tests := []struct {
		name string
		args args
		want aform.ExportTmplAttrs
	}{
		{
			name: "with a word",
			args: args{
				name:  "class",
				value: "big",
			},
			want: aform.ExportTmplAttrs{"class": "big"},
		},
		{
			name: "with a space",
			args: args{
				name:  "class",
				value: "big blue",
			},
			want: aform.ExportTmplAttrs{"class": "big blue"},
		},
		{
			name: "with an empty string",
			args: args{
				name:  "class",
				value: "",
			},
			want: aform.ExportTmplAttrs{"class": ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, aform.ExportAttributableAttribute(aform.Attr(tt.args.name, tt.args.value)), "Attr(%v, %v)", tt.args.name, tt.args.value)
		})
	}
}

func TestAttr_withCustomType(t *testing.T) {
	type args struct {
		name  string
		value testString
	}
	tests := []struct {
		name string
		args args
		want aform.ExportTmplAttrs
	}{
		{
			name: "with a word",
			args: args{
				name:  "class",
				value: testString("big"),
			},
			want: aform.ExportTmplAttrs{"class": "big"},
		},
		{
			name: "with a space",
			args: args{
				name:  "class",
				value: testString("big blue"),
			},
			want: aform.ExportTmplAttrs{"class": "big blue"},
		},
		{
			name: "with an empty string",
			args: args{
				name:  "class",
				value: testString(""),
			},
			want: aform.ExportTmplAttrs{"class": ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, aform.ExportAttributableAttribute(aform.Attr(tt.args.name, tt.args.value)), "Attr(%v, %v)", tt.args.name, tt.args.value)
		})
	}
}

func Test_attributableListToAttrs(t *testing.T) {
	type args struct {
		attrs []aform.Attributable
	}
	tests := []struct {
		name string
		args args
		want aform.ExportTmplAttrs
	}{
		{
			name: "with a name/value attribute",
			args: args{
				attrs: []aform.Attributable{aform.Attr("class", "big")},
			},
			want: aform.ExportTmplAttrs{"class": "big"},
		},
		{
			name: "with a boolean attribute",
			args: args{
				attrs: []aform.Attributable{aform.BoolAttr("required")},
			},
			want: aform.ExportTmplAttrs{"required": ""},
		},
		{
			name: "with a multiple different attributes",
			args: args{
				attrs: []aform.Attributable{
					aform.BoolAttr("required"),
					aform.Attr("id", "me"),
					aform.Attr("min", 7),
					aform.BoolAttr("disabled"),
				},
			},
			want: aform.ExportTmplAttrs{"disabled": "", "id": "me", "min": "7", "required": ""},
		},
		{
			name: "with a multiple times the same boolean attribute",
			args: args{
				attrs: []aform.Attributable{
					aform.BoolAttr("required"),
					aform.BoolAttr("required"),
					aform.BoolAttr("required"),
				},
			},
			want: aform.ExportTmplAttrs{"required": ""},
		},
		{
			name: "with a multiple times the same attribute name",
			args: args{
				attrs: []aform.Attributable{
					aform.Attr("id", "first"),
					aform.Attr("id", "second"),
					aform.Attr("id", "third"),
					aform.Attr("id", "last"),
				},
			},
			want: aform.ExportTmplAttrs{"id": "last"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, aform.ExportAttributableListToAttrs(tt.args.attrs), "attributableListToAttrs(%v)", tt.args.attrs)
		})
	}
}

func TestExportAttributableListToAttrs_panicWithType(t *testing.T) {
	a := assert.New(t)
	a.PanicsWithValue(
		"You can't directly set type with functions WithAttributes or SetAttributes. To change the type, use a different Field and/or set a different Widget on an existing Field.",
		func() {
			aform.ExportAttributableListToAttrs([]aform.Attributable{aform.BoolAttr("type")})
		})
	a.PanicsWithValue(
		"You can't directly set type with functions WithAttributes or SetAttributes. To change the type, use a different Field and/or set a different Widget on an existing Field.",
		func() {
			aform.ExportAttributableListToAttrs([]aform.Attributable{aform.Attr("type", "anything")})
		})
}

func TestExportAttributableListToAttrs_panicWithName(t *testing.T) {
	a := assert.New(t)
	a.PanicsWithValue(
		"You can't directly set name with functions WithAttributes or SetAttributes. To change the name attribute, set a different name to the Field when you create it.",
		func() {
			aform.ExportAttributableListToAttrs([]aform.Attributable{aform.BoolAttr("name")})
		})
	a.PanicsWithValue(
		"You can't directly set name with functions WithAttributes or SetAttributes. To change the name attribute, set a different name to the Field when you create it.",
		func() {
			aform.ExportAttributableListToAttrs([]aform.Attributable{aform.Attr("name", "anything")})
		})
}

func TestExportAttributableListToAttrs_panicWithValue(t *testing.T) {
	a := assert.New(t)
	a.PanicsWithValue(
		"You can't directly set value with functions WithAttributes or SetAttributes. To change the value attribute, set an initial value when you create the field or bind values with BindRequest or BindData.",
		func() {
			aform.ExportAttributableListToAttrs([]aform.Attributable{aform.BoolAttr("value")})
		})
	a.PanicsWithValue(
		"You can't directly set value with functions WithAttributes or SetAttributes. To change the value attribute, set an initial value when you create the field or bind values with BindRequest or BindData.",
		func() {
			aform.ExportAttributableListToAttrs([]aform.Attributable{aform.Attr("value", "anything")})
		})
}
