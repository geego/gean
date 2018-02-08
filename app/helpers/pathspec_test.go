package helpers

import (
	"testing"

	"github.com/geego/gean/app/geanfs"
	"github.com/gostores/configurator"
	"github.com/gostores/require"
)

func TestNewPathSpecFromConfig(t *testing.T) {
	v := configurator.New()
	l := NewLanguage("no", v)
	v.Set("disablePathToLower", true)
	v.Set("removePathAccents", true)
	v.Set("uglyURLs", true)
	v.Set("multilingual", true)
	v.Set("defaultContentLanguageInSubdir", true)
	v.Set("defaultContentLanguage", "no")
	v.Set("canonifyURLs", true)
	v.Set("paginatePath", "side")
	v.Set("baseURL", "http://base.com")
	v.Set("themesDir", "thethemes")
	v.Set("layoutDir", "thelayouts")
	v.Set("workingDir", "thework")
	v.Set("staticDir", "thestatic")
	v.Set("theme", "thetheme")

	p, err := NewPathSpec(geanfs.NewMem(v), l)

	require.NoError(t, err)
	require.True(t, p.canonifyURLs)
	require.True(t, p.defaultContentLanguageInSubdir)
	require.True(t, p.disablePathToLower)
	require.True(t, p.multilingual)
	require.True(t, p.removePathAccents)
	require.True(t, p.uglyURLs)
	require.Equal(t, "no", p.defaultContentLanguage)
	require.Equal(t, "no", p.language.Lang)
	require.Equal(t, "side", p.paginatePath)

	require.Equal(t, "http://base.com", p.BaseURL.String())
	require.Equal(t, "thethemes", p.themesDir)
	require.Equal(t, "thelayouts", p.layoutDir)
	require.Equal(t, "thework", p.workingDir)
	require.Equal(t, "thetheme", p.theme)
}
