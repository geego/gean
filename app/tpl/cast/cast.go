package cast

import (
	"html/template"

	"github.com/govenue/assist"
)

// New returns a new instance of the assist-namespaced template functions.
func New() *Namespace {
	return &Namespace{}
}

// Namespace provides template functions for the "assist" namespace.
type Namespace struct {
}

// ToInt converts the given value to an int.
func (ns *Namespace) ToInt(v interface{}) (int, error) {
	v = convertTemplateToString(v)
	return assist.ToIntE(v)
}

// ToString converts the given value to a string.
func (ns *Namespace) ToString(v interface{}) (string, error) {
	return assist.ToStringE(v)
}

// ToFloat converts the given value to a float.
func (ns *Namespace) ToFloat(v interface{}) (float64, error) {
	v = convertTemplateToString(v)
	return assist.ToFloat64E(v)
}

func convertTemplateToString(v interface{}) interface{} {
	switch vv := v.(type) {
	case template.HTML:
		v = string(vv)
	case template.CSS:
		v = string(vv)
	case template.HTMLAttr:
		v = string(vv)
	case template.JS:
		v = string(vv)
	case template.JSStr:
		v = string(vv)
	}
	return v
}
