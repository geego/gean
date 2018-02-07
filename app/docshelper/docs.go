package docshelper

import (
	"encoding/json"
)

// DocProviders contains all DocProviders added to the system.
var DocProviders = make(map[string]DocProvider)

// AddDocProvider adds or updates the DocProvider for a given name.
func AddDocProvider(name string, provider DocProvider) {
	DocProviders[name] = provider
}

// DocProvider is used to save arbitrary JSON data
// used for the generation of the documentation.
type DocProvider func() map[string]interface{}

// MarshalJSON returns a JSON representation of the DocProvider.
func (d DocProvider) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(d(), "", "  ")
}
