package command

import (
	"errors"

	"github.com/govenue/goman"

	"github.com/geego/gean/app/releaser"
)

func init() {
	HugoCmd.AddCommand(createReleaser().cmd)
}

type releaseCommandeer struct {
	cmd *goman.Command

	version string

	skipPublish bool
	try         bool
}

func createReleaser() *releaseCommandeer {
	// Note: This is a command only meant for internal use and must be run
	// via "go run -tags release main.go release" on the actual code base that is in the release.
	r := &releaseCommandeer{
		cmd: &goman.Command{
			Use:    "release",
			Short:  "Release a new version of gean.",
			Hidden: true,
		},
	}

	r.cmd.RunE = func(cmd *goman.Command, args []string) error {
		return r.release()
	}

	r.cmd.PersistentFlags().StringVarP(&r.version, "rel", "r", "", "new release version, i.e. 0.25.1")
	r.cmd.PersistentFlags().BoolVarP(&r.skipPublish, "skip-publish", "", false, "skip all publishing pipes of the release")
	r.cmd.PersistentFlags().BoolVarP(&r.try, "try", "", false, "simulate a release, i.e. no changes")

	return r
}

func (r *releaseCommandeer) release() error {
	if r.version == "" {
		return errors.New("must set the --rel flag to the relevant version number")
	}
	return releaser.New(r.version, r.skipPublish, r.try).Run()
}
