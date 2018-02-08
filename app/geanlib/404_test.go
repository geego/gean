package geanlib

import (
	"path/filepath"
	"testing"

	"github.com/geego/gean/app/deps"
)

func Test404(t *testing.T) {
	t.Parallel()
	var (
		cfg, fs = newTestCfg()
		th      = testHelper{cfg, fs, t}
	)

	cfg.Set("baseURL", "http://auth/bub/")

	writeSource(t, fs, filepath.Join("layouts", "404.html"), "<html><body>Not Found!</body></html>")
	writeSource(t, fs, filepath.Join("content", "page.md"), "A page")

	buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})

	// Note: We currently have only 1 404 page. One might think that we should have
	// multiple, to follow the Custom Output scheme, but I don't see how that wold work
	// right now.
	th.assertFileContent("public/404.html", "Not Found")

}
