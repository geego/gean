package types

import (
	"fmt"

	"github.com/govenue/assist"
)

// KeyValues holds an key and a slice of values.
type KeyValues struct {
	Key    interface{}
	Values []interface{}
}

// KeyString returns the key as a string, an empty string if conversion fails.
func (k KeyValues) KeyString() string {
	return assist.ToString(k.Key)
}

func (k KeyValues) String() string {
	return fmt.Sprintf("%v: %v", k.Key, k.Values)
}

func NewKeyValuesStrings(key string, values ...string) KeyValues {
	iv := make([]interface{}, len(values))
	for i := 0; i < len(values); i++ {
		iv[i] = values[i]
	}
	return KeyValues{Key: key, Values: iv}
}
