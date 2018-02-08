package geanlib

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/geego/gean/app/deps"
)

func TestAllTemplateEngines(t *testing.T) {
	t.Parallel()
	noOp := func(s string) string {
		return s
	}

	amberFixer := func(s string) string {
		fixed := strings.Replace(s, "{{ .Title", "{{ Title", -1)
		fixed = strings.Replace(fixed, ".Content", "Content", -1)
		fixed = strings.Replace(fixed, ".IsNamedParams", "IsNamedParams", -1)
		fixed = strings.Replace(fixed, "{{", "#{", -1)
		fixed = strings.Replace(fixed, "}}", "}", -1)
		fixed = strings.Replace(fixed, `title "hello world"`, `title("hello world")`, -1)

		return fixed
	}

	for _, config := range []struct {
		suffix        string
		templateFixer func(s string) string
	}{
		{"amber", amberFixer},
		{"html", noOp},
		{"ace", noOp},
	} {
		t.Run(config.suffix,
			func(t *testing.T) {
				doTestTemplateEngine(t, config.suffix, config.templateFixer)
			})
	}

}

func doTestTemplateEngine(t *testing.T, suffix string, templateFixer func(s string) string) {

	cfg, fs := newTestCfg()

	t.Log("Testing", suffix)

	templTemplate := `
p
	|
	| Page Title: {{ .Title }}
	br
	| Page Content: {{ .Content }}
	br
	| {{ title "hello world" }}

`

	templShortcodeTemplate := `
p
	|
	| Shortcode: {{ .IsNamedParams }}
`

	templ := templateFixer(templTemplate)
	shortcodeTempl := templateFixer(templShortcodeTemplate)

	writeSource(t, fs, filepath.Join("content", "p.md"), `
---
title: My Title 
---
My Content

Shortcode: {{< myShort >}}

`)

	writeSource(t, fs, filepath.Join("layouts", "_default", fmt.Sprintf("single.%s", suffix)), templ)
	writeSource(t, fs, filepath.Join("layouts", "shortcodes", fmt.Sprintf("myShort.%s", suffix)), shortcodeTempl)

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})
	th := testHelper{s.Cfg, s.Fs, t}

	th.assertFileContent(filepath.Join("public", "p", "index.html"),
		"Page Title: My Title",
		"My Content",
		"Hello World",
		"Shortcode: false",
	)

}
