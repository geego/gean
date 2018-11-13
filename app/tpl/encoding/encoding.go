package encoding

import (
	"encoding/base64"
	"encoding/json"
	"html/template"

	"github.com/govenue/assist"
)

// New returns a new instance of the encoding-namespaced template functions.
func New() *Namespace {
	return &Namespace{}
}

// Namespace provides template functions for the "encoding" namespace.
type Namespace struct{}

// Base64Decode returns the base64 decoding of the given content.
func (ns *Namespace) Base64Decode(content interface{}) (string, error) {
	conv, err := assist.ToStringE(content)
	if err != nil {
		return "", err
	}

	dec, err := base64.StdEncoding.DecodeString(conv)
	return string(dec), err
}

// Base64Encode returns the base64 encoding of the given content.
func (ns *Namespace) Base64Encode(content interface{}) (string, error) {
	conv, err := assist.ToStringE(content)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString([]byte(conv)), nil
}

// Jsonify encodes a given object to JSON.
func (ns *Namespace) Jsonify(v interface{}) (template.HTML, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return template.HTML(b), nil
}
