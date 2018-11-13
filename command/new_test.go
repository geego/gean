package command

import (
	"path/filepath"
	"testing"

	"github.com/govenue/assert"
	"github.com/govenue/configurator"
	"github.com/govenue/require"

	"github.com/geego/gean/app/geanfs"
)

// Issue #1133
func TestNewContentPathSectionWithForwardSlashes(t *testing.T) {
	p, s := newContentPathSection("/post/new.md")
	assert.Equal(t, filepath.FromSlash("/post/new.md"), p)
	assert.Equal(t, "post", s)
}

func checkNewSiteInited(fs *geanfs.Fs, basepath string, t *testing.T) {

	paths := []string{
		filepath.Join(basepath, "layouts"),
		filepath.Join(basepath, "content"),
		filepath.Join(basepath, "archetypes"),
		filepath.Join(basepath, "static"),
		filepath.Join(basepath, "data"),
		filepath.Join(basepath, "config.toml"),
	}

	for _, path := range paths {
		_, err := fs.Source.Stat(path)
		require.NoError(t, err)
	}
}

func TestDoNewSite(t *testing.T) {
	basepath := filepath.Join("base", "blog")
	_, fs := newTestCfg()

	require.NoError(t, doNewSite(fs, basepath, false))

	checkNewSiteInited(fs, basepath, t)
}

func TestDoNewSite_noerror_base_exists_but_empty(t *testing.T) {
	basepath := filepath.Join("base", "blog")
	_, fs := newTestCfg()

	require.NoError(t, fs.Source.MkdirAll(basepath, 777))

	require.NoError(t, doNewSite(fs, basepath, false))
}

func TestDoNewSite_error_base_exists(t *testing.T) {
	basepath := filepath.Join("base", "blog")
	_, fs := newTestCfg()

	require.NoError(t, fs.Source.MkdirAll(basepath, 777))
	_, err := fs.Source.Create(filepath.Join(basepath, "foo"))
	require.NoError(t, err)
	// Since the directory already exists and isn't empty, expect an error
	require.Error(t, doNewSite(fs, basepath, false))

}

func TestDoNewSite_force_empty_dir(t *testing.T) {
	basepath := filepath.Join("base", "blog")
	_, fs := newTestCfg()

	require.NoError(t, fs.Source.MkdirAll(basepath, 777))

	require.NoError(t, doNewSite(fs, basepath, true))

	checkNewSiteInited(fs, basepath, t)
}

func TestDoNewSite_error_force_dir_inside_exists(t *testing.T) {
	basepath := filepath.Join("base", "blog")
	_, fs := newTestCfg()

	contentPath := filepath.Join(basepath, "content")

	require.NoError(t, fs.Source.MkdirAll(contentPath, 777))
	require.Error(t, doNewSite(fs, basepath, true))
}

func TestDoNewSite_error_force_config_inside_exists(t *testing.T) {
	basepath := filepath.Join("base", "blog")
	_, fs := newTestCfg()

	configPath := filepath.Join(basepath, "config.toml")
	require.NoError(t, fs.Source.MkdirAll(basepath, 777))
	_, err := fs.Source.Create(configPath)
	require.NoError(t, err)

	require.Error(t, doNewSite(fs, basepath, true))
}

func newTestCfg() (*configurator.Configurator, *geanfs.Fs) {

	v := configurator.New()
	fs := geanfs.NewMem(v)

	v.SetFs(fs.Source)

	return v, fs

}
