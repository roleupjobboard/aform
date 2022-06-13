package aform

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

// Attributable defines a common interface for name/value HTML attributes and
// boolean HTML attributes.
type Attributable interface {
	name() string
	attribute() tmplAttrs
}

// AttrName defines the HTML attribute name type
type AttrName interface {
	~string
}

// AttrValue defines the HTML attribute value type
type AttrValue interface {
	constraints.Integer | ~string
}

// Attr returns an HTML attribute named name and that is equal to value.
// Example to add attribute placeholder="****"
//	Attr("placeholder", "****")
// Example to add attribute maxlength="125"
//	Attr("maxlength", 125)
func Attr[N AttrName, V AttrValue](name N, value V) Attributable {
	return nameValueAttr[V]{n: string(name), v: value}
}

// BoolAttr returns a boolean HTML attribute named name.
// Example to add attribute readonly
//	BoolAttr("readonly")
func BoolAttr[N AttrName](name N) Attributable {
	return boolAttr(name)
}

func attributableListToAttrs(attrs []Attributable) tmplAttrs {
	panicOnForbiddenAttributeName(attrs)
	output := tmplAttrs{}
	for _, attr := range attrs {
		for name, value := range attr.attribute() {
			output[name] = value
		}
	}
	return output
}

func panicOnForbiddenAttributeName(attrs []Attributable) {
	l := []string{"type", "name", "value"}
	p := "You can't directly set %s with functions WithAttributes or SetAttributes."
	m := map[string]string{
		"type":  fmt.Sprintf(p, "type") + " To change the type, use a different Field and/or set a different Widget on an existing Field.",
		"name":  fmt.Sprintf(p, "name") + " To change the name attribute, set a different name to the Field when you create it.",
		"value": fmt.Sprintf(p, "value") + " To change the value attribute, set an initial value when you create the field or bind values with BindRequest or BindData.",
	}
	for _, attr := range attrs {
		if slices.Contains(l, attr.name()) {
			panic(m[attr.name()])
		}
	}
}

type nameValueAttr[T AttrValue] struct {
	n string
	v T
}

func (attr nameValueAttr[T]) name() string {
	return attr.n
}

func (attr nameValueAttr[T]) attribute() tmplAttrs {
	return map[string]string{attr.name(): fmt.Sprint(attr.v)}
}

type boolAttr string

func (attr boolAttr) name() string {
	return string(attr)
}

func (attr boolAttr) attribute() tmplAttrs {
	return map[string]string{attr.name(): ""}
}
