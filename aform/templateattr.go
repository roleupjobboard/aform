package aform

import (
	"fmt"
	"golang.org/x/exp/slices"
	"html/template"
	"strings"
)

type tmplAttrs map[string]string

// HTMLAttributes transforms the attributes in HTML attributes and orders them.
//
// The order followed is:class, id, name, data-*, src, for, type, href, value,
// minlength, maxlength, title, alt, role, aria-* and finally others attributes.
//
// In each group, attributes are alphabetically ordered.
func (a tmplAttrs) HTMLAttributes() []template.HTMLAttr {
	output := make([]template.HTMLAttr, 0, len(a))
	output = append(output, htmlAttributeFromKey("class", a)...)
	output = append(output, htmlAttributeFromKey("id", a)...)
	output = append(output, htmlAttributeFromKey("name", a)...)
	output = append(output, htmlAttributeFromKeyPrefix("data-", a)...)
	output = append(output, htmlAttributeFromKey("src", a)...)
	output = append(output, htmlAttributeFromKey("for", a)...)
	output = append(output, htmlAttributeFromKey("type", a)...)
	output = append(output, htmlAttributeFromKey("href", a)...)
	output = append(output, htmlAttributeFromKey("value", a)...)
	output = append(output, htmlAttributeFromKey("minlength", a)...)
	output = append(output, htmlAttributeFromKey("maxlength", a)...)
	output = append(output, htmlAttributeFromKey("title", a)...)
	output = append(output, htmlAttributeFromKey("alt", a)...)
	output = append(output, htmlAttributeFromKey("role", a)...)
	output = append(output, htmlAttributeFromKeyPrefix("aria-", a)...)
	if len(output) < len(a) {
		output = append(output, remainingAttributes(a)...)
	}
	return output
}

// Value returns the value of the attribute named name or an empty string if it's not present.
func (a tmplAttrs) Value(name string) template.HTMLAttr {
	if value, ok := a[name]; ok {
		return template.HTMLAttr(value)
	}
	return ""
}

// Attr returns the attribute named name or an empty string if it's not present.
func (a tmplAttrs) Attr(name string) template.HTMLAttr {
	if value, ok := a[name]; ok {
		return keyValueToHTMLAttr(name, value)
	}
	return ""
}

func htmlAttributeFromKey(key string, attrs tmplAttrs) []template.HTMLAttr {
	if value, ok := attrs[key]; ok {
		return []template.HTMLAttr{keyValueToHTMLAttr(key, value)}
	}
	return nil
}

func htmlAttributeFromKeyPrefix(prefix string, attrs tmplAttrs) []template.HTMLAttr {
	var output []template.HTMLAttr
	for key, value := range attrs {
		if strings.HasPrefix(key, prefix) {
			output = append(output, keyValueToHTMLAttr(key, value))
		}
	}
	if len(output) == 0 {
		return nil
	}
	slices.SortFunc(output, func(a, b template.HTMLAttr) bool {
		return a < b
	})
	return output
}

func keyValueToHTMLAttr(key, value string) template.HTMLAttr {
	if len(value) > 0 {
		return template.HTMLAttr(fmt.Sprintf("%s=\"%s\"", key, value))
	} else {
		return template.HTMLAttr(key)
	}
}

var attrNameList = []string{"class", "id", "name", "src", "for", "type", "href", "value", "minlength", "maxlength", "title", "alt", "role"}
var attrPrefixList = []string{"data-", "aria-"}

func remainingAttributes(attrs tmplAttrs) []template.HTMLAttr {
	var output []template.HTMLAttr
	isPrefixed := func(key string) bool {
		for _, prefix := range attrPrefixList {
			if strings.HasPrefix(key, prefix) {
				return true
			}
		}
		return false
	}
	for key, value := range attrs {
		if slices.Contains(attrNameList, key) {
			continue
		}
		if isPrefixed(key) {
			continue
		}
		output = append(output, keyValueToHTMLAttr(key, value))
	}
	if len(output) == 0 {
		return nil
	}
	slices.SortFunc(output, func(a, b template.HTMLAttr) bool {
		return a < b
	})
	return output
}
