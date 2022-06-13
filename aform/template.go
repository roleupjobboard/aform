package aform

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

// Highly inspired from https://github.com/django/django/tree/main/django/forms/jinja2/django/forms/widgets
//
// To indent templates eventually https://stackoverflow.com/a/59401696/21052

type labelContent interface {
	~string
}

type label[C labelContent] struct {
	UseTag          bool
	Label           C
	Suffix          C
	LabelWithSuffix C
	Attrs           tmplAttrs
}

func newLabel(useTag bool, lbl, suffix string, attrs tmplAttrs) label[string] {
	return label[string]{
		UseTag:          useTag,
		Label:           lbl,
		Suffix:          suffix,
		LabelWithSuffix: labelWithSuffix(lbl, suffix),
		Attrs:           attrs,
	}
}

func newSafeLabel(useTag bool, lbl, suffix string, attrs tmplAttrs) label[template.HTML] {
	return label[template.HTML]{
		UseTag:          useTag,
		Label:           template.HTML(lbl),
		Suffix:          template.HTML(suffix),
		LabelWithSuffix: template.HTML(labelWithSuffix(lbl, suffix)),
		Attrs:           attrs,
	}
}

func labelWithSuffix(lbl, suffix string) string {
	if len(suffix) == 0 {
		return lbl
	}
	if strings.LastIndexAny(lbl, ".!?:") == len(lbl)-1 {
		return lbl
	}
	return lbl + suffix
}

type widgetInput struct {
	Type  Widget
	Name  string
	Value string
	Attrs tmplAttrs
}

func (w widgetInput) HTMLType() string {
	return w.Type.htmlType()
}

func (w widgetInput) HTMLNameAttribute() template.HTMLAttr {
	return keyValueToHTMLAttr("name", w.Name)
}

type widgetChoice struct {
	Type   Widget
	Name   string
	Values []string
	Groups []map[string][]widgetOption
	Attrs  tmplAttrs
}

func (w widgetChoice) HTMLType() string {
	return w.Type.htmlType()
}

type widgetOption struct {
	Label     string
	WrapLabel bool
	widgetInput
}

type tmplErrors struct {
	List  []tmplError
	Attrs tmplAttrs
}

type tmplError struct {
	Text  string
	Attrs tmplAttrs
}

type tmplHelpText struct {
	Text  template.HTML
	Attrs tmplAttrs
}

var buildingBlocksTemplateDefinitions = []map[string]string{
	{"attrs": `{{with .Attrs}}{{ range $attr := .HTMLAttributes }} {{ $attr }}{{end}}{{end}}`},
	{"errors": `{{if .Errors.List}}<ul{{ template "attrs" .Errors }}>{{ range $error := .Errors.List }}<li{{ template "attrs" $error }}>{{$error.Text}}</li>{{end}}</ul>{{end}}`},
	{"help_text": `{{with .HelpText}}<span{{ template "attrs" . }}>{{.Text}}</span>{{end}}`},
	{"select_option": `<option value="{{ .Value }}"{{ template "attrs" . }}>{{ .Label }}</option>`},
}
var labelTemplateDefinitions = []map[string]string{
	{"label": `{{if .Label.UseTag}}<label{{ template "attrs" .Label }}>{{ .Label.LabelWithSuffix }}</label>{{else}}{{ .Label.LabelWithSuffix }}{{end}}`},
	{"legend": `{{if .Label.UseTag}}<legend{{ template "attrs" .Label }}>{{ .Label.LabelWithSuffix }}</legend>{{else}}{{ .Label.LabelWithSuffix }}{{end}}`},
}
var widgetTemplateDefinitions = []map[string]string{
	{"input": `<input type="{{ .HTMLType }}" {{ .HTMLNameAttribute }}{{with .Value}} value="{{ . }}"{{end}}{{ template "attrs" . }}>`},
	{"input_option": `{{if .WrapLabel}}<label{{with .Attrs.Value "id"}} for="{{ . }}"{{end}}>{{end}}{{ template "input" . }}{{if .WrapLabel}}{{.Label}}</label>{{end}}`},
	{"text": `{{ template "input" .Widget }}`},
	{"email": `{{ template "input" .Widget }}`},
	{"url": `{{ template "input" .Widget }}`},
	{"password": `{{ template "input" .Widget }}`},
	{"hidden": `{{ template "input" .Widget }}`},
	{"checkbox": `{{ template "input" .Widget }}`},
	{"textarea": `<textarea name="{{ .Widget.Name }}"{{ template "attrs" .Widget }}>
{{with .Widget.Value }}{{ . }}{{ end }}</textarea>`},
	{"select": `{{$group_tab := ""}}{{$option_tab := "  "}}<select name="{{ .Widget.Name }}"{{ template "attrs" .Widget }}>{{range $group := .Widget.Groups}}{{range $group_name, $group_options := $group}}{{with $group_name}}{{$group_tab = "  "}}
{{$group_tab}}<optgroup label="{{ . }}">{{else}}{{$group_tab = ""}}{{end}}{{range $options := $group_options}}
{{$group_tab}}{{$option_tab}}{{ template "select_option" $options}}{{end}}{{with $group_name}}
{{$group_tab}}</optgroup>{{end}}{{end}}{{ end }}
</select>`},
	{"multiple_input": `<div{{with .Widget.Attrs.Value "id"}} id="{{ . }}"{{end}}{{with .Widget.Attrs.Attr "class"}} {{ . }}{{end}}>{{range $group := .Widget.Groups}}{{range $group_name, $group_options := $group}}{{with $group_name}}
<div><label>{{ . }}</label>{{end}}{{range $options := $group_options}}
{{template "input_option" $options}}{{end}}{{with $group_name}}
</div>{{end}}{{end}}{{end}}
</div>`},
	{"checkbox_select": `{{ template "multiple_input" . }}`},
	{"radio": `{{ template "multiple_input" . }}`},
}
var formTemplateDefinitions = []map[string]string{
	{"field_as_div": `<div{{ with .CSSClasses }} class="{{.}}"{{end}}>{{if .UseFieldset}}
<fieldset>{{ .LegendTag }}
{{if .HasErrors}}{{ .Errors }}
{{end}}{{else}}{{ .LabelTag }}{{if .HasErrors}}
{{ .Errors }}
{{end}}{{end}}{{ .Widget }}{{if .HasHelpText}}
{{.HelpText}}{{end}}{{if .UseFieldset}}
</fieldset>
{{end}}</div>`},
	{"form_as_div": `{{- range .Form.Fields}}
{{ .AsDiv }}
{{- end}}`},
}

func mustFormAsDivTemplate(f *Form) template.HTML {
	tmpl, err := formAsDivTemplate(f)
	if err != nil {
		panic(fmt.Sprintf("mustFormAsDivTemplate: %s", err.Error()))
	}
	return tmpl
}

func formAsDivTemplate(f *Form) (template.HTML, error) {
	t := loadTemplates()
	buf := &bytes.Buffer{}
	err := t.ExecuteTemplate(buf, "form_as_div", map[string]interface{}{"Form": f})
	if err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

func mustFieldAsDivTemplate(fld *Field) template.HTML {
	tmpl, err := fieldAsDivTemplate(fld)
	if err != nil {
		panic(fmt.Sprintf("mustFieldAsDivTemplate: %s", err.Error()))
	}
	return tmpl
}

func fieldAsDivTemplate(fld *Field) (template.HTML, error) {
	t := loadTemplates()
	buf := &bytes.Buffer{}
	err := t.ExecuteTemplate(buf, "field_as_div", fld)
	if err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

func mustLabelTemplate[C labelContent](label *label[C]) template.HTML {
	tmpl, err := labelTemplate(label)
	if err != nil {
		panic(fmt.Sprintf("mustLabelTemplate: %s", err.Error()))
	}
	return tmpl
}

func labelTemplate[C labelContent](label *label[C]) (template.HTML, error) {
	return labelOrLegendTemplate("label", label)
}

func mustLegendTemplate[C labelContent](label *label[C]) template.HTML {
	tmpl, err := legendTemplate(label)
	if err != nil {
		panic(fmt.Sprintf("mustLegendTemplate: %s", err.Error()))
	}
	return tmpl
}

func legendTemplate[C labelContent](label *label[C]) (template.HTML, error) {
	return labelOrLegendTemplate("legend", label)
}

func labelOrLegendTemplate[C labelContent](tag string, label *label[C]) (template.HTML, error) {
	t := loadTemplates()
	buf := &bytes.Buffer{}
	err := t.ExecuteTemplate(buf, tag, map[string]interface{}{"Label": label})
	if err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

func mustInputTemplate(widget *widgetInput) template.HTML {
	tmpl, err := inputTemplate(widget)
	if err != nil {
		panic(fmt.Sprintf("mustInputTemplate: %s", err.Error()))
	}
	return tmpl
}

func inputTemplate(widget *widgetInput) (template.HTML, error) {
	t := loadTemplates()
	buf := &bytes.Buffer{}
	err := t.ExecuteTemplate(buf, widget.HTMLType(), map[string]interface{}{"Widget": widget})
	if err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

func mustChoiceTemplate(widget *widgetChoice) template.HTML {
	tmpl, err := choiceTemplate(widget)
	if err != nil {
		panic(fmt.Sprintf("mustChoiceTemplate: %s", err.Error()))
	}
	return tmpl
}

func choiceTemplate(widget *widgetChoice) (template.HTML, error) {
	t := loadTemplates()
	buf := &bytes.Buffer{}
	err := t.ExecuteTemplate(buf, widget.HTMLType(), map[string]interface{}{"Widget": widget})
	if err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

func mustErrorsTemplate(errors *tmplErrors) template.HTML {
	tmpl, err := errorsTemplate(errors)
	if err != nil {
		panic(fmt.Sprintf("mustErrorsTemplate: %s", err.Error()))
	}
	return tmpl
}

func errorsTemplate(errors *tmplErrors) (template.HTML, error) {
	t := loadTemplates()
	buf := &bytes.Buffer{}
	err := t.ExecuteTemplate(buf, "errors", map[string]interface{}{"Errors": errors})
	if err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

func mustHelpTextTemplate(helpText *tmplHelpText) template.HTML {
	tmpl, err := helpTextTemplate(helpText)
	if err != nil {
		panic(fmt.Sprintf("mustHelpTextTemplate: %s", err.Error()))
	}
	return tmpl
}

func helpTextTemplate(helpText *tmplHelpText) (template.HTML, error) {
	t := loadTemplates()
	buf := &bytes.Buffer{}
	err := t.ExecuteTemplate(buf, "help_text", map[string]interface{}{"HelpText": helpText})
	if err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

var _templates *template.Template = nil

func loadTemplates() *template.Template {
	if _templates != nil {
		return _templates
	}
	var all []map[string]string
	all = append(all, buildingBlocksTemplateDefinitions...)
	all = append(all, labelTemplateDefinitions...)
	all = append(all, widgetTemplateDefinitions...)
	all = append(all, formTemplateDefinitions...)
	_templates = loadTemplateList(all)
	return _templates
}

func loadTemplateList(list []map[string]string) *template.Template {
	var tmpl *template.Template = nil
	for _, t := range list {
		for k, v := range t {
			if tmpl == nil {
				tmpl = template.Must(template.New(k).Parse(v))
			} else {
				template.Must(tmpl.New(k).Parse(v))
			}
		}
	}
	return tmpl
}
