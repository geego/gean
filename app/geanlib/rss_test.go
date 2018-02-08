package geanlib

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/geego/gean/app/deps"
)

func TestRSSOutput(t *testing.T) {
	t.Parallel()
	var (
		cfg, fs = newTestCfg()
		th      = testHelper{cfg, fs, t}
	)

	rssLimit := len(weightedSources) - 1

	rssURI := "customrss.xml"

	cfg.Set("baseURL", "http://auth/bub/")
	cfg.Set("rssURI", rssURI)
	cfg.Set("title", "RSSTest")
	cfg.Set("rssLimit", rssLimit)

	for _, src := range weightedSources {
		writeSource(t, fs, filepath.Join("content", "sect", src.Name), string(src.Content))
	}

	buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})

	// Home RSS
	th.assertFileContent(filepath.Join("public", rssURI), "<?xml", "rss version", "RSSTest")
	// Section RSS
	th.assertFileContent(filepath.Join("public", "sect", rssURI), "<?xml", "rss version", "Sects on RSSTest")
	// Taxonomy RSS
	th.assertFileContent(filepath.Join("public", "categories", "hugo", rssURI), "<?xml", "rss version", "Hugo on RSSTest")

	// RSS Item Limit
	content := readDestination(t, fs, filepath.Join("public", rssURI))
	c := strings.Count(content, "<item>")
	if c != rssLimit {
		t.Errorf("incorrect RSS item count: expected %d, got %d", rssLimit, c)
	}
}
