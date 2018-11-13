package create_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/geego/gean/app/create"
	"github.com/geego/gean/app/deps"
	"github.com/geego/gean/app/geanfs"
	"github.com/geego/gean/app/geanlib"
	"github.com/geego/gean/app/helpers"
	"github.com/govenue/configurator"
	"github.com/govenue/fsintra"
	"github.com/govenue/require"
)

func TestNewContent(t *testing.T) {
	v := configurator.New()
	initConfigurator(v)

	cases := []struct {
		kind     string
		path     string
		expected []string
	}{
		{"post", "post/sample-1.md", []string{`title = "Post Arch title"`, `test = "test1"`, "date = \"2015-01-12T19:20:04-07:00\""}},
		{"post", "post/org-1.org", []string{`#+title: ORG-1`}},
		{"emptydate", "post/sample-ed.md", []string{`title = "Empty Date Arch title"`, `test = "test1"`}},
		{"stump", "stump/sample-2.md", []string{`title: "Sample 2"`}},      // no archetype file
		{"", "sample-3.md", []string{`title: "Sample 3"`}},                 // no archetype
		{"product", "product/sample-4.md", []string{`title = "SAMPLE-4"`}}, // empty archetype front matter
		{"shortcodes", "shortcodes/go.md", []string{
			`title = "GO"`,
			"{{< myshortcode >}}",
			"{{% myshortcode %}}",
			"{{</* comment */>}}\n{{%/* comment */%}}"}}, // shortcodes
	}

	for _, c := range cases {
		cfg, fs := newTestCfg()
		ps, err := helpers.NewPathSpec(fs, cfg)
		require.NoError(t, err)
		h, err := geanlib.NewHugoSites(deps.DepsCfg{Cfg: cfg, Fs: fs})
		require.NoError(t, err)
		require.NoError(t, initFs(fs))

		siteFactory := func(filename string, siteUsed bool) (*geanlib.Site, error) {
			return h.Sites[0], nil
		}

		require.NoError(t, create.NewContent(ps, siteFactory, c.kind, c.path))

		fname := filepath.Join("content", filepath.FromSlash(c.path))
		content := readFileFromFs(t, fs.Source, fname)
		for i, v := range c.expected {
			found := strings.Contains(content, v)
			if !found {
				t.Errorf("[%d] %q missing from output:\n%q", i, v, content)
			}
		}
	}
}

func initConfigurator(v *configurator.Configurator) {
	v.Set("metaDataFormat", "toml")
	v.Set("archetypeDir", "archetypes")
	v.Set("contentDir", "content")
	v.Set("themesDir", "themes")
	v.Set("layoutDir", "layouts")
	v.Set("i18nDir", "i18n")
	v.Set("theme", "sample")
}

func initFs(fs *geanfs.Fs) error {
	perm := os.FileMode(0755)
	var err error

	// create directories
	dirs := []string{
		"archetypes",
		"content",
		filepath.Join("themes", "sample", "archetypes"),
	}
	for _, dir := range dirs {
		err = fs.Source.Mkdir(dir, perm)
		if err != nil {
			return err
		}
	}

	// create files
	for _, v := range []struct {
		path    string
		content string
	}{
		{
			path:    filepath.Join("archetypes", "post.md"),
			content: "+++\ndate = \"2015-01-12T19:20:04-07:00\"\ntitle = \"Post Arch title\"\ntest = \"test1\"\n+++\n",
		},
		{
			path:    filepath.Join("archetypes", "post.org"),
			content: "#+title: {{ .BaseFileName  | upper }}",
		},
		{
			path: filepath.Join("archetypes", "product.md"),
			content: `+++
title = "{{ .BaseFileName  | upper }}"
+++`,
		},
		{
			path:    filepath.Join("archetypes", "emptydate.md"),
			content: "+++\ndate =\"\"\ntitle = \"Empty Date Arch title\"\ntest = \"test1\"\n+++\n",
		},
		// #3623x
		{
			path: filepath.Join("archetypes", "shortcodes.md"),
			content: `+++
title = "{{ .BaseFileName  | upper }}"
+++

{{< myshortcode >}}

Some text.

{{% myshortcode %}}
{{</* comment */>}}
{{%/* comment */%}}


`,
		},
	} {
		f, err := fs.Source.Create(v.path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = f.Write([]byte(v.content))
		if err != nil {
			return err
		}
	}

	return nil
}

// TODO(bep) extract common testing package with this and some others
func readFileFromFs(t *testing.T, fs fsintra.Fs, filename string) string {
	filename = filepath.FromSlash(filename)
	b, err := fsintra.ReadFile(fs, filename)
	if err != nil {
		// Print some debug info
		root := strings.Split(filename, helpers.FilePathSeparator)[0]
		fsintra.Walk(fs, root, func(path string, info os.FileInfo, err error) error {
			if info != nil && !info.IsDir() {
				fmt.Println("    ", path)
			}

			return nil
		})
		t.Fatalf("Failed to read file: %s", err)
	}
	return string(b)
}

func newTestCfg() (*configurator.Configurator, *geanfs.Fs) {

	v := configurator.New()
	fs := geanfs.NewMem(v)

	v.SetFs(fs.Source)

	initConfigurator(v)

	return v, fs

}
