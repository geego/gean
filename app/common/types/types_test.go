package types

import (
	"testing"

	"github.com/govenue/require"
)

func TestKeyValues(t *testing.T) {
	assert := require.New(t)

	kv := NewKeyValuesStrings("key", "a1", "a2")

	assert.Equal("key", kv.KeyString())
	assert.Equal([]interface{}{"a1", "a2"}, kv.Values)
}
