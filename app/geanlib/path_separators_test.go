package geanlib

import (
	"path/filepath"
	"strings"
	"testing"
)

var simplePageYAML = `---
contenttype: ""
---
Sample Text
`

func TestDegenerateMissingFolderInPageFilename(t *testing.T) {
	t.Parallel()
	s := newTestSite(t)
	p, err := s.NewPageFrom(strings.NewReader(simplePageYAML), filepath.Join("foobar"))
	if err != nil {
		t.Fatalf("Error in NewPageFrom")
	}
	if p.Section() != "" {
		t.Fatalf("No section should be set for a file path: foobar")
	}
}
