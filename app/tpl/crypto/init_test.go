package crypto

import (
	"testing"

	"github.com/gostores/require"

	"yiqilai.tech/gean/app/deps"
	"yiqilai.tech/gean/app/tpl/internal"
)

func TestInit(t *testing.T) {
	var found bool
	var ns *internal.TemplateFuncsNamespace

	for _, nsf := range internal.TemplateFuncsNamespaceRegistry {
		ns = nsf(&deps.Deps{})
		if ns.Name == name {
			found = true
			break
		}
	}

	require.True(t, found)
	require.IsType(t, &Namespace{}, ns.Context())
}
