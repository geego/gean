package geanlib

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/geego/gean/app/config"
	"github.com/geego/gean/app/deps"
	"github.com/geego/gean/app/geanfs"
	"github.com/geego/gean/app/helpers"
	"github.com/geego/gean/app/source"
	"github.com/geego/gean/app/tpl"
	"github.com/govenue/configurator"
	"github.com/govenue/fsintra"
	"github.com/govenue/notepad"
	"github.com/govenue/require"
)

type testHelper struct {
	Cfg config.Provider
	Fs  *geanfs.Fs
	T   testing.TB
}

func (th testHelper) assertFileContent(filename string, matches ...string) {
	filename = th.replaceDefaultContentLanguageValue(filename)
	content := readDestination(th.T, th.Fs, filename)
	for _, match := range matches {
		match = th.replaceDefaultContentLanguageValue(match)
		require.True(th.T, strings.Contains(content, match), fmt.Sprintf("File no match for\n%q in\n%q:\n%s", strings.Replace(match, "%", "%%", -1), filename, strings.Replace(content, "%", "%%", -1)))
	}
}

// TODO(bep) better name for this. It does no magic replacements depending on defaultontentLanguageInSubDir.
func (th testHelper) assertFileContentStraight(filename string, matches ...string) {
	content := readDestination(th.T, th.Fs, filename)
	for _, match := range matches {
		require.True(th.T, strings.Contains(content, match), fmt.Sprintf("File no match for\n%q in\n%q:\n%s", strings.Replace(match, "%", "%%", -1), filename, strings.Replace(content, "%", "%%", -1)))
	}
}

func (th testHelper) assertFileContentRegexp(filename string, matches ...string) {
	filename = th.replaceDefaultContentLanguageValue(filename)
	content := readDestination(th.T, th.Fs, filename)
	for _, match := range matches {
		match = th.replaceDefaultContentLanguageValue(match)
		r := regexp.MustCompile(match)
		require.True(th.T, r.MatchString(content), fmt.Sprintf("File no match for\n%q in\n%q:\n%s", strings.Replace(match, "%", "%%", -1), filename, strings.Replace(content, "%", "%%", -1)))
	}
}

func (th testHelper) assertFileNotExist(filename string) {
	exists, err := helpers.Exists(filename, th.Fs.Destination)
	require.NoError(th.T, err)
	require.False(th.T, exists)
}

func (th testHelper) replaceDefaultContentLanguageValue(value string) string {
	defaultInSubDir := th.Cfg.GetBool("defaultContentLanguageInSubDir")
	replace := th.Cfg.GetString("defaultContentLanguage") + "/"

	if !defaultInSubDir {
		value = strings.Replace(value, replace, "", 1)

	}
	return value
}

func newTestPathSpec(fs *geanfs.Fs, v *configurator.Configurator) *helpers.PathSpec {
	l := helpers.NewDefaultLanguage(v)
	ps, _ := helpers.NewPathSpec(fs, l)
	return ps
}

func newTestDefaultPathSpec() *helpers.PathSpec {
	v := configurator.New()
	// Easier to reason about in tests.
	v.Set("disablePathToLower", true)
	fs := geanfs.NewDefault(v)
	ps, _ := helpers.NewPathSpec(fs, v)
	return ps
}

func newTestCfg() (*configurator.Configurator, *geanfs.Fs) {

	v := configurator.New()
	fs := geanfs.NewMem(v)

	v.SetFs(fs.Source)

	loadDefaultSettingsFor(v)

	// Default is false, but true is easier to use as default in tests
	v.Set("defaultContentLanguageInSubdir", true)

	return v, fs

}

// newTestSite creates a new site in the  English language with in-memory Fs.
// The site will have a template system loaded and ready to use.
// Note: This is only used in single site tests.
func newTestSite(t testing.TB, configKeyValues ...interface{}) *Site {

	cfg, fs := newTestCfg()

	for i := 0; i < len(configKeyValues); i += 2 {
		cfg.Set(configKeyValues[i].(string), configKeyValues[i+1])
	}

	d := deps.DepsCfg{Language: helpers.NewLanguage("en", cfg), Fs: fs, Cfg: cfg}

	s, err := NewSiteForCfg(d)

	if err != nil {
		t.Fatalf("Failed to create Site: %s", err)
	}
	return s
}

func newTestSitesFromConfig(t testing.TB, afs fsintra.Fs, tomlConfig string, layoutPathContentPairs ...string) (testHelper, *HugoSites) {
	if len(layoutPathContentPairs)%2 != 0 {
		t.Fatalf("Layouts must be provided in pairs")
	}

	writeToFs(t, afs, "config.toml", tomlConfig)

	cfg, err := LoadConfig(afs, "", "config.toml")
	require.NoError(t, err)

	fs := geanfs.NewFrom(afs, cfg)
	th := testHelper{cfg, fs, t}

	for i := 0; i < len(layoutPathContentPairs); i += 2 {
		writeSource(t, fs, layoutPathContentPairs[i], layoutPathContentPairs[i+1])
	}

	h, err := NewHugoSites(deps.DepsCfg{Fs: fs, Cfg: cfg})

	require.NoError(t, err)

	return th, h
}

func newTestSitesFromConfigWithDefaultTemplates(t testing.TB, tomlConfig string) (testHelper, *HugoSites) {
	return newTestSitesFromConfig(t, fsintra.NewMemMapFs(), tomlConfig,
		"layouts/_default/single.html", "Single|{{ .Title }}|{{ .Content }}",
		"layouts/_default/list.html", "List|{{ .Title }}|{{ .Content }}",
		"layouts/_default/terms.html", "Terms List|{{ .Title }}|{{ .Content }}",
	)
}

func newDebugLogger() *notepad.Notepad {
	return notepad.NewNotepad(notepad.LevelDebug, notepad.LevelError, os.Stdout, ioutil.Discard, "", log.Ldate|log.Ltime)
}

func newErrorLogger() *notepad.Notepad {
	return notepad.NewNotepad(notepad.LevelError, notepad.LevelError, os.Stdout, ioutil.Discard, "", log.Ldate|log.Ltime)
}
func createWithTemplateFromNameValues(additionalTemplates ...string) func(templ tpl.TemplateHandler) error {

	return func(templ tpl.TemplateHandler) error {
		for i := 0; i < len(additionalTemplates); i += 2 {
			err := templ.AddTemplate(additionalTemplates[i], additionalTemplates[i+1])
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func buildSingleSite(t testing.TB, depsCfg deps.DepsCfg, buildCfg BuildCfg) *Site {
	return buildSingleSiteExpected(t, false, depsCfg, buildCfg)
}

func buildSingleSiteExpected(t testing.TB, expectBuildError bool, depsCfg deps.DepsCfg, buildCfg BuildCfg) *Site {
	h, err := NewHugoSites(depsCfg)

	require.NoError(t, err)
	require.Len(t, h.Sites, 1)

	if expectBuildError {
		require.Error(t, h.Build(buildCfg))
		return nil

	}

	require.NoError(t, h.Build(buildCfg))

	return h.Sites[0]
}

func writeSourcesToSource(t *testing.T, base string, fs *geanfs.Fs, sources ...source.ByteSource) {
	for _, src := range sources {
		writeSource(t, fs, filepath.Join(base, src.Name), string(src.Content))
	}
}

func isCI() bool {
	return os.Getenv("CI") != ""
}
