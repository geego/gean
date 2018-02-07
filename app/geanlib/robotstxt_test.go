package geanlib

import (
	"path/filepath"
	"testing"

	"yiqilai.tech/gean/app/deps"
)

const robotTxtTemplate = `User-agent: Googlebot
  {{ range .Data.Pages }}
	Disallow: {{.RelPermalink}}
	{{ end }}
`

func TestRobotsTXTOutput(t *testing.T) {
	t.Parallel()
	var (
		cfg, fs = newTestCfg()
		th      = testHelper{cfg, fs, t}
	)

	cfg.Set("baseURL", "http://auth/bub/")
	cfg.Set("enableRobotsTXT", true)

	writeSource(t, fs, filepath.Join("layouts", "robots.txt"), robotTxtTemplate)
	writeSourcesToSource(t, "content", fs, weightedSources...)

	buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})

	th.assertFileContent("public/robots.txt", "User-agent: Googlebot")

}
