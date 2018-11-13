package geanlib

import (
	"html/template"
	"path/filepath"
	"testing"

	"github.com/geego/gean/app/deps"
	"github.com/geego/gean/app/source"
	"github.com/govenue/require"
)

const slugDoc1 = "---\ntitle: slug doc 1\nslug: slug-doc-1\naliases:\n - sd1/foo/\n - sd2\n - sd3/\n - sd4.html\n---\nslug doc 1 content\n"

const slugDoc2 = `---
title: slug doc 2
slug: slug-doc-2
---
slug doc 2 content
`

var urlFakeSource = []source.ByteSource{
	{Name: filepath.FromSlash("content/blue/doc1.md"), Content: []byte(slugDoc1)},
	{Name: filepath.FromSlash("content/blue/doc2.md"), Content: []byte(slugDoc2)},
}

// Issue #1105
func TestShouldNotAddTrailingSlashToBaseURL(t *testing.T) {
	t.Parallel()
	for i, this := range []struct {
		in       string
		expected string
	}{
		{"http://base.com/", "http://base.com/"},
		{"http://base.com/sub/", "http://base.com/sub/"},
		{"http://base.com/sub", "http://base.com/sub"},
		{"http://base.com", "http://base.com"}} {

		cfg, fs := newTestCfg()
		cfg.Set("baseURL", this.in)
		d := deps.DepsCfg{Cfg: cfg, Fs: fs}
		s, err := NewSiteForCfg(d)
		require.NoError(t, err)
		s.initializeSiteInfo()

		if s.Info.BaseURL() != template.URL(this.expected) {
			t.Errorf("[%d] got %s expected %s", i, s.Info.BaseURL(), this.expected)
		}
	}
}

func TestPageCount(t *testing.T) {
	t.Parallel()
	cfg, fs := newTestCfg()
	cfg.Set("uglyURLs", false)
	cfg.Set("paginate", 10)

	writeSourcesToSource(t, "", fs, urlFakeSource...)
	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})

	_, err := s.Fs.Destination.Open("public/blue")
	if err != nil {
		t.Errorf("No indexed rendered.")
	}

	for _, pth := range []string{
		"public/sd1/foo/index.html",
		"public/sd2/index.html",
		"public/sd3/index.html",
		"public/sd4.html",
	} {
		if _, err := s.Fs.Destination.Open(filepath.FromSlash(pth)); err != nil {
			t.Errorf("No alias rendered: %s", pth)
		}
	}
}
