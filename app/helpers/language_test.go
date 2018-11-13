package helpers

import (
	"testing"

	"github.com/govenue/configurator"
	"github.com/govenue/require"
)

func TestGetGlobalOnlySetting(t *testing.T) {
	v := configurator.New()
	lang := NewDefaultLanguage(v)
	lang.SetParam("defaultContentLanguageInSubdir", false)
	lang.SetParam("paginatePath", "side")
	v.Set("defaultContentLanguageInSubdir", true)
	v.Set("paginatePath", "page")

	require.True(t, lang.GetBool("defaultContentLanguageInSubdir"))
	require.Equal(t, "side", lang.GetString("paginatePath"))
}
