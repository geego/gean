package geanlib

import (
	"testing"

	"github.com/gostores/assert"
	"github.com/gostores/fsintra"
	"github.com/gostores/require"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	// Add a random config variable for testing.
	// side = page in Norwegian.
	configContent := `
	PaginatePath = "side"
	`

	mm := fsintra.NewMemMapFs()

	writeToFs(t, mm, "hugo.toml", configContent)

	cfg, err := LoadConfig(mm, "", "hugo.toml")
	require.NoError(t, err)

	assert.Equal(t, "side", cfg.GetString("paginatePath"))
	// default
	assert.Equal(t, "layouts", cfg.GetString("layoutDir"))
}
func TestLoadMultiConfig(t *testing.T) {
	t.Parallel()

	// Add a random config variable for testing.
	// side = page in Norwegian.
	configContentBase := `
	DontChange = "same"
	PaginatePath = "side"
	`
	configContentSub := `
	PaginatePath = "top"
	`
	mm := fsintra.NewMemMapFs()

	writeToFs(t, mm, "base.toml", configContentBase)

	writeToFs(t, mm, "override.toml", configContentSub)

	cfg, err := LoadConfig(mm, "", "base.toml,override.toml")
	require.NoError(t, err)

	assert.Equal(t, "top", cfg.GetString("paginatePath"))
	assert.Equal(t, "same", cfg.GetString("DontChange"))
}
