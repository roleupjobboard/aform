package aform

import (
	"github.com/stretchr/testify/assert"
	"html/template"
	"testing"
)

func Test_labelWithSuffix(t *testing.T) {
	type fields struct {
		Label  string
		Suffix string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"without suffix",
			fields{
				Label:  "Name",
				Suffix: "",
			},
			"Name",
		},
		{
			"with suffix",
			fields{
				Label:  "Name",
				Suffix: ":",
			},
			"Name:",
		},
		{
			"with suffix and end with .",
			fields{
				Label:  "Name.",
				Suffix: ">",
			},
			"Name.",
		},
		{
			"with suffix and end with !",
			fields{
				Label:  "Name!",
				Suffix: ">",
			},
			"Name!",
		},
		{
			"with suffix and end with ?",
			fields{
				Label:  "Name?",
				Suffix: ">",
			},
			"Name?",
		},
		{
			"with suffix and end with :",
			fields{
				Label:  "Name:",
				Suffix: ">",
			},
			"Name:",
		},
		{
			"with suffix and end with a space",
			fields{
				Label:  "Name: ",
				Suffix: ">",
			},
			"Name: >",
		},
		{
			"with suffix, contains : and end with :",
			fields{
				Label:  "N:ame:",
				Suffix: ">",
			},
			"N:ame:",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, labelWithSuffix(tt.fields.Label, tt.fields.Suffix), "LabelWithSuffix()")
		})
	}
}

func TestFieldTemplate(t *testing.T) {
	type args struct {
		fld *Field
	}
	tests := []struct {
		name    string
		args    args
		want    template.HTML
		wantErr bool
	}{
		{
			name: "happy char field",
			args: args{
				fld: func() *Field {
					fld, _ := NewCharField("Happy Field", "yeah!", "", 0, 0)
					return fld.field()
				}(),
			},
			want:    `<div><label for="id_happy_field">Happy Field</label><input type="text" name="happy_field" value="yeah!" id="id_happy_field" required></div>`,
			wantErr: false,
		},
		{
			name: "happy char field with option",
			args: args{
				fld: func() *Field {
					fld, _ := NewCharField("Happy Field", "yeah!", "", 0, 0, WithLabel("Joyful Field"), IsNotRequired(), IsDisabled())
					return fld.field()
				}(),
			},
			want:    `<div><label for="id_happy_field">Joyful Field</label><input type="text" name="happy_field" value="yeah!" id="id_happy_field" disabled></div>`,
			wantErr: false,
		},
		{
			name: "Boolean field checked",
			args: args{
				fld: func() *Field {
					fld, _ := NewBooleanField("Boolean Field", true, WithLabel("Boolean Field Checked"), IsDisabled())
					return fld.field()
				}(),
			},
			want:    `<div><label for="id_boolean_field">Boolean Field Checked</label><input type="checkbox" name="boolean_field" id="id_boolean_field" checked disabled required></div>`,
			wantErr: false,
		},
		{
			name: "Boolean field not checked",
			args: args{
				fld: func() *Field {
					fld, _ := DefaultBooleanField("Boolean Field", WithLabel("Boolean Field Checked"), IsDisabled())
					fld.Clean("off")
					return fld.field()
				}(),
			},
			want:    `<div><label for="id_boolean_field">Boolean Field Checked</label><input type="checkbox" name="boolean_field" id="id_boolean_field" disabled required></div>`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fieldAsDivTemplate(tt.args.fld)
			if (err != nil) != tt.wantErr {
				t.Errorf("FieldTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FieldTemplate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLabelTemplate(t *testing.T) {
	type args struct {
		label label[string]
	}
	tests := []struct {
		name    string
		args    args
		want    template.HTML
		wantErr bool
	}{
		{
			name: "just the label",
			args: args{
				label: newLabel(false, "Hello", "", nil),
			},
			want:    "Hello",
			wantErr: false,
		},
		{
			name: "with label tag without attributes",
			args: args{
				label: newLabel(true, "Hello", "", nil),
			},
			want:    "<label>Hello</label>",
			wantErr: false,
		},
		{
			name: "with label tag and attributes",
			args: args{
				label: newLabel(true, "Hello", "", map[string]string{"class": "required big blue"}),
			},
			want:    `<label class="required big blue">Hello</label>`,
			wantErr: false,
		},
		{
			name: "HTML-escaped",
			args: args{
				label: newLabel(false, "Hello <strong>good</strong> friend", "", nil),
			},
			want:    "Hello &lt;strong&gt;good&lt;/strong&gt; friend",
			wantErr: false,
		},
		{
			name: "HTML-escaped with suffix",
			args: args{
				label: newLabel(false, "Hello <strong>good</strong> friend", "<strong>*</strong>", nil),
			},
			want:    "Hello &lt;strong&gt;good&lt;/strong&gt; friend&lt;strong&gt;*&lt;/strong&gt;",
			wantErr: false,
		},
		{
			name: "HTML-escaped with label tag without attributes",
			args: args{
				label: newLabel(true, "Hello <strong>good</strong> friend", "", nil),
			},
			want:    "<label>Hello &lt;strong&gt;good&lt;/strong&gt; friend</label>",
			wantErr: false,
		},
		{
			name: "HTML-escaped with label tag and attributes",
			args: args{
				label: newLabel(true, "Hello <strong>good</strong> friend", "", map[string]string{"class": "required big blue"}),
			},
			want:    `<label class="required big blue">Hello &lt;strong&gt;good&lt;/strong&gt; friend</label>`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := labelTemplate(&tt.args.label)
			if (err != nil) != tt.wantErr {
				t.Errorf("labelTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equalf(t, tt.want, got, "labelTemplate(%v)", tt.args.label)
		})
	}
}

func TestSafeLabelTemplate(t *testing.T) {
	type args struct {
		label label[template.HTML]
	}
	tests := []struct {
		name    string
		args    args
		want    template.HTML
		wantErr bool
	}{
		{
			name: "just the label",
			args: args{
				label: newSafeLabel(false, "Hello", "", nil),
			},
			want:    "Hello",
			wantErr: false,
		},
		{
			name: "with label tag without attributes",
			args: args{
				label: newSafeLabel(true, "Hello", "", nil),
			},
			want:    "<label>Hello</label>",
			wantErr: false,
		},
		{
			name: "with label tag and attributes",
			args: args{
				label: newSafeLabel(true, "Hello", "", map[string]string{"class": "required big blue"}),
			},
			want:    `<label class="required big blue">Hello</label>`,
			wantErr: false,
		},
		{
			name: "not HTML-escaped",
			args: args{
				label: newSafeLabel(false, "Hello <strong>good</strong> friend", "", nil),
			},
			want:    "Hello <strong>good</strong> friend",
			wantErr: false,
		},
		{
			name: "not HTML-escaped with suffix",
			args: args{
				label: newSafeLabel(false, "Hello <strong>good</strong> friend", "<strong>*</strong>", nil),
			},
			want:    "Hello <strong>good</strong> friend<strong>*</strong>",
			wantErr: false,
		},
		{
			name: "not HTML-escaped with label tag without attributes",
			args: args{
				label: newSafeLabel(true, "Hello <strong>good</strong> friend", "", nil),
			},
			want:    "<label>Hello <strong>good</strong> friend</label>",
			wantErr: false,
		},
		{
			name: "not HTML-escaped with label tag and attributes",
			args: args{
				label: newSafeLabel(true, "Hello <strong>good</strong> friend", "", map[string]string{"class": "required big blue"}),
			},
			want:    `<label class="required big blue">Hello <strong>good</strong> friend</label>`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := labelTemplate(&tt.args.label)
			if (err != nil) != tt.wantErr {
				t.Errorf("labelTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equalf(t, tt.want, got, "labelTemplate(%v)", tt.args.label)
		})
	}
}

func TestInputTemplate(t *testing.T) {
	type args struct {
		widget *widgetInput
	}
	tests := []struct {
		name    string
		args    args
		want    template.HTML
		wantErr bool
	}{
		{
			name: "happy path without value",
			args: args{
				widget: &widgetInput{
					Type:  TextInput,
					Name:  "happypath",
					Value: "",
				},
			},
			want:    `<input type="text" name="happypath">`,
			wantErr: false,
		},
		{
			name: "happy path with value",
			args: args{
				widget: &widgetInput{
					Type:  TextInput,
					Name:  "happypath",
					Value: "yeah",
				},
			},
			want:    `<input type="text" name="happypath" value="yeah">`,
			wantErr: false,
		},
		{
			name: "happy path without value but attributes",
			args: args{
				widget: &widgetInput{
					Type:  TextInput,
					Name:  "happypath",
					Value: "",
					Attrs: map[string]string{"class": "magic", "required": ""},
				},
			},
			want:    `<input type="text" name="happypath" class="magic" required>`,
			wantErr: false,
		},
		{
			name: "happy path with value and attributes",
			args: args{
				widget: &widgetInput{
					Type:  TextInput,
					Name:  "happypath",
					Value: "very happy",
					Attrs: map[string]string{"selected": "", "class": "magic hidden", "required": ""},
				},
			},
			want:    `<input type="text" name="happypath" value="very happy" class="magic hidden" required selected>`,
			wantErr: false,
		},
		{
			name: "check box checked",
			args: args{
				widget: &widgetInput{
					Type:  CheckboxInput,
					Name:  "check_box",
					Value: "",
					Attrs: map[string]string{"checked": "", "class": "magic hidden", "required": ""},
				},
			},
			want:    `<input type="checkbox" name="check_box" class="magic hidden" checked required>`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := inputTemplate(tt.args.widget)
			if (err != nil) != tt.wantErr {
				t.Errorf("inputTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("inputTemplate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChoiceTemplate_select(t *testing.T) {
	type args struct {
		widget *widgetChoice
	}
	tests := []struct {
		name    string
		args    args
		want    template.HTML
		wantErr bool
	}{
		{
			name: "no options",
			args: args{
				widget: &widgetChoice{
					Type:   Select,
					Name:   "color",
					Values: nil,
					Groups: nil,
					Attrs:  nil,
				},
			},
			want: `<select name="color">
</select>`,
			wantErr: false,
		},
		{
			name: "no options with attributes",
			args: args{
				widget: &widgetChoice{
					Type:   Select,
					Name:   "color",
					Values: nil,
					Groups: nil,
					Attrs:  map[string]string{"id": "color_id", "class": "required"},
				},
			},
			want: `<select name="color" class="required" id="color_id">
</select>`,
			wantErr: false,
		},
		{
			name: "with options",
			args: args{
				widget: &widgetChoice{
					Type:   Select,
					Name:   "color",
					Values: nil,
					Groups: []map[string][]widgetOption{
						{"": []widgetOption{
							{
								Label:       "Orange Color",
								widgetInput: widgetInput{Value: "orange"},
							},
							{
								Label:       "Pink Color",
								widgetInput: widgetInput{Value: "pink"},
							},
						}},
					},
					Attrs: nil,
				},
			},
			want: `<select name="color">
  <option value="orange">Orange Color</option>
  <option value="pink">Pink Color</option>
</select>`,
			wantErr: false,
		},
		{
			name: "with one option selected",
			args: args{
				widget: &widgetChoice{
					Type:   Select,
					Name:   "color",
					Values: nil,
					Groups: []map[string][]widgetOption{
						{"": []widgetOption{
							{
								Label:       "Orange Color",
								widgetInput: widgetInput{Value: "orange", Attrs: map[string]string{"selected": ""}},
							},
							{
								Label:       "Pink Color",
								widgetInput: widgetInput{Value: "pink"},
							},
						}},
					},
					Attrs: nil,
				},
			},
			want: `<select name="color">
  <option value="orange" selected>Orange Color</option>
  <option value="pink">Pink Color</option>
</select>`,
			wantErr: false,
		},
		{
			name: "with group of options",
			args: args{
				widget: &widgetChoice{
					Type:   Select,
					Name:   "color",
					Values: nil,
					Groups: []map[string][]widgetOption{
						{"RGB": []widgetOption{
							{
								Label:       "Red Color",
								widgetInput: widgetInput{Value: "red"},
							},
							{
								Label:       "Green Color",
								widgetInput: widgetInput{Value: "green"},
							},
							{
								Label:       "Blue Color",
								widgetInput: widgetInput{Value: "blue"},
							},
						}},
					},
					Attrs: nil,
				},
			},
			want: `<select name="color">
  <optgroup label="RGB">
    <option value="red">Red Color</option>
    <option value="green">Green Color</option>
    <option value="blue">Blue Color</option>
  </optgroup>
</select>`,
			wantErr: false,
		},
		{
			name: "with one option in a group selected",
			args: args{
				widget: &widgetChoice{
					Type:   Select,
					Name:   "color",
					Values: nil,
					Groups: []map[string][]widgetOption{
						{"RGB": []widgetOption{
							{
								Label:       "Red Color",
								widgetInput: widgetInput{Value: "red"},
							},
							{
								Label:       "Green Color",
								widgetInput: widgetInput{Value: "green"},
							},
							{
								Label:       "Blue Color",
								widgetInput: widgetInput{Value: "blue", Attrs: map[string]string{"selected": ""}},
							},
						}},
					},
					Attrs: nil,
				},
			},
			want: `<select name="color">
  <optgroup label="RGB">
    <option value="red">Red Color</option>
    <option value="green">Green Color</option>
    <option value="blue" selected>Blue Color</option>
  </optgroup>
</select>`,
			wantErr: false,
		},
		{
			name: "with naked options and group of options",
			args: args{
				widget: &widgetChoice{
					Type:   Select,
					Name:   "color",
					Values: nil,
					Groups: []map[string][]widgetOption{
						{"": []widgetOption{
							{
								Label:       "Orange Color",
								widgetInput: widgetInput{Value: "orange"},
							},
							{
								Label:       "Pink Color",
								widgetInput: widgetInput{Value: "pink"},
							},
						}},
						{"RGB": []widgetOption{
							{
								Label:       "Red Color",
								widgetInput: widgetInput{Value: "red"},
							},
							{
								Label:       "Green Color",
								widgetInput: widgetInput{Value: "green"},
							},
							{
								Label:       "Blue Color",
								widgetInput: widgetInput{Value: "blue"},
							},
						}},
					},
					Attrs: nil,
				},
			},
			want: `<select name="color">
  <option value="orange">Orange Color</option>
  <option value="pink">Pink Color</option>
  <optgroup label="RGB">
    <option value="red">Red Color</option>
    <option value="green">Green Color</option>
    <option value="blue">Blue Color</option>
  </optgroup>
</select>`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := choiceTemplate(tt.args.widget)
			if (err != nil) != tt.wantErr {
				t.Errorf("choiceTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equalf(t, tt.want, got, "choiceTemplate(%v)", tt.args.widget)
		})
	}
}
func TestChoiceTemplate_radio(t *testing.T) {
	type args struct {
		widget *widgetChoice
	}
	tests := []struct {
		name    string
		args    args
		want    template.HTML
		wantErr bool
	}{
		{
			name: "no options",
			args: args{
				widget: &widgetChoice{
					Type:   RadioSelect,
					Name:   "color",
					Values: nil,
					Groups: nil,
					Attrs:  nil,
				},
			},
			want: `<div>
</div>`,
			wantErr: false,
		},
		{
			name: "no options with attributes",
			args: args{
				widget: &widgetChoice{
					Type:   RadioSelect,
					Name:   "color",
					Values: nil,
					Groups: nil,
					Attrs:  map[string]string{"id": "color_id", "class": "required"},
				},
			},
			want: `<div id="color_id" class="required">
</div>`,
			wantErr: false,
		},
		{
			name: "with options",
			args: args{
				widget: &widgetChoice{
					Type:   RadioSelect,
					Name:   "color",
					Values: nil,
					Groups: []map[string][]widgetOption{
						{"": []widgetOption{
							{
								Label:       "Orange Color",
								WrapLabel:   true,
								widgetInput: widgetInput{Type: RadioSelect, Name: "color", Value: "orange"},
							},
							{
								Label:       "Pink Color",
								WrapLabel:   true,
								widgetInput: widgetInput{Type: RadioSelect, Name: "color", Value: "pink"},
							},
						}},
					},
					Attrs: nil,
				},
			},
			want: `<div>
<label><input type="radio" name="color" value="orange">Orange Color</label>
<label><input type="radio" name="color" value="pink">Pink Color</label>
</div>`,
			wantErr: false,
		},
		{
			name: "with one option selected",
			args: args{
				widget: &widgetChoice{
					Type:   RadioSelect,
					Name:   "color",
					Values: nil,
					Groups: []map[string][]widgetOption{
						{"": []widgetOption{
							{
								Label:       "Orange Color",
								WrapLabel:   true,
								widgetInput: widgetInput{Type: RadioSelect, Name: "color", Value: "orange", Attrs: map[string]string{"checked": ""}},
							},
							{
								Label:       "Pink Color",
								WrapLabel:   true,
								widgetInput: widgetInput{Type: RadioSelect, Name: "color", Value: "pink"},
							},
						}},
					},
					Attrs: nil,
				},
			},
			want: `<div>
<label><input type="radio" name="color" value="orange" checked>Orange Color</label>
<label><input type="radio" name="color" value="pink">Pink Color</label>
</div>`,
			wantErr: false,
		},
		{
			name: "with group of options",
			args: args{
				widget: &widgetChoice{
					Type:   RadioSelect,
					Name:   "color",
					Values: nil,
					Groups: []map[string][]widgetOption{
						{"RGB": []widgetOption{
							{
								Label:       "Red Color",
								WrapLabel:   true,
								widgetInput: widgetInput{Type: RadioSelect, Name: "color", Value: "red"},
							},
							{
								Label:       "Green Color",
								WrapLabel:   true,
								widgetInput: widgetInput{Type: RadioSelect, Name: "color", Value: "green"},
							},
							{
								Label:       "Blue Color",
								WrapLabel:   true,
								widgetInput: widgetInput{Type: RadioSelect, Name: "color", Value: "blue"},
							},
						}},
					},
					Attrs: nil,
				},
			},
			want: `<div>
<div><label>RGB</label>
<label><input type="radio" name="color" value="red">Red Color</label>
<label><input type="radio" name="color" value="green">Green Color</label>
<label><input type="radio" name="color" value="blue">Blue Color</label>
</div>
</div>`,
			wantErr: false,
		},
		{
			name: "with one option in a group selected",
			args: args{
				widget: &widgetChoice{
					Type:   RadioSelect,
					Name:   "color",
					Values: nil,
					Groups: []map[string][]widgetOption{
						{"RGB": []widgetOption{
							{
								Label:       "Red Color",
								WrapLabel:   true,
								widgetInput: widgetInput{Type: RadioSelect, Name: "color", Value: "red"},
							},
							{
								Label:       "Green Color",
								WrapLabel:   true,
								widgetInput: widgetInput{Type: RadioSelect, Name: "color", Value: "green"},
							},
							{
								Label:       "Blue Color",
								WrapLabel:   true,
								widgetInput: widgetInput{Type: RadioSelect, Name: "color", Value: "blue", Attrs: map[string]string{"checked": ""}},
							},
						}},
					},
					Attrs: nil,
				},
			},
			want: `<div>
<div><label>RGB</label>
<label><input type="radio" name="color" value="red">Red Color</label>
<label><input type="radio" name="color" value="green">Green Color</label>
<label><input type="radio" name="color" value="blue" checked>Blue Color</label>
</div>
</div>`,
			wantErr: false,
		},
		{
			name: "with naked options and group of options",
			args: args{
				widget: &widgetChoice{
					Type:   RadioSelect,
					Name:   "color",
					Values: nil,
					Groups: []map[string][]widgetOption{
						{"": []widgetOption{
							{
								Label:       "Orange Color",
								WrapLabel:   true,
								widgetInput: widgetInput{Type: RadioSelect, Name: "color", Value: "orange"},
							},
							{
								Label:       "Pink Color",
								WrapLabel:   true,
								widgetInput: widgetInput{Type: RadioSelect, Name: "color", Value: "pink"},
							},
						}},
						{"RGB": []widgetOption{
							{
								Label:       "Red Color",
								WrapLabel:   true,
								widgetInput: widgetInput{Type: RadioSelect, Name: "color", Value: "red"},
							},
							{
								Label:       "Green Color",
								WrapLabel:   true,
								widgetInput: widgetInput{Type: RadioSelect, Name: "color", Value: "green"},
							},
							{
								Label:       "Blue Color",
								WrapLabel:   true,
								widgetInput: widgetInput{Type: RadioSelect, Name: "color", Value: "blue"},
							},
						}},
					},
					Attrs: nil,
				},
			},
			want: `<div>
<label><input type="radio" name="color" value="orange">Orange Color</label>
<label><input type="radio" name="color" value="pink">Pink Color</label>
<div><label>RGB</label>
<label><input type="radio" name="color" value="red">Red Color</label>
<label><input type="radio" name="color" value="green">Green Color</label>
<label><input type="radio" name="color" value="blue">Blue Color</label>
</div>
</div>`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := choiceTemplate(tt.args.widget)
			if (err != nil) != tt.wantErr {
				t.Errorf("choiceTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equalf(t, tt.want, got, "choiceTemplate(%v)", tt.args.widget)
		})
	}
}
