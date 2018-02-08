package geanlib

import (
	"fmt"
	"html/template"

	"github.com/geego/gean/app/helpers"
)

var (
	// CommitHash contains the current Git revision. Use make to build to make
	// sure this gets set.
	CommitHash string

	// BuildDate contains the date of the current build.
	BuildDate string
)

var hugoInfo *HugoInfo

// HugoInfo contains information about the current Hugo environment
type HugoInfo struct {
	Version    string
	Generator  template.HTML
	CommitHash string
	BuildDate  string
}

func init() {
	hugoInfo = &HugoInfo{
		Version:    helpers.CurrentHugoVersion.String(),
		CommitHash: CommitHash,
		BuildDate:  BuildDate,
		Generator:  template.HTML(fmt.Sprintf(`<meta name="generator" content="Hugo %s" />`, helpers.CurrentHugoVersion.String())),
	}
}
