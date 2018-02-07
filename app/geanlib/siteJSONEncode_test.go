package geanlib

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"yiqilai.tech/gean/app/deps"
)

// Issue #1123
// Testing prevention of cyclic refs in JSON encoding
// May be smart to run with: -timeout 4000ms
func TestEncodePage(t *testing.T) {
	t.Parallel()
	cfg, fs := newTestCfg()

	// borrowed from menu_test.go
	for _, src := range menuPageSources {
		writeSource(t, fs, filepath.Join("content", src.Name), string(src.Content))

	}

	s := buildSingleSite(t, deps.DepsCfg{Fs: fs, Cfg: cfg}, BuildCfg{})

	_, err := json.Marshal(s)
	check(t, err)

	_, err = json.Marshal(s.RegularPages[0])
	check(t, err)
}

func check(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Failed %s", err)
	}
}
