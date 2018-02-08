package cast

import (
	"testing"

	"github.com/geego/gean/app/deps"
	"github.com/geego/gean/app/tpl/internal"
	"github.com/gostores/require"
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
