package command

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/govenue/goman"
	"github.com/govenue/notepad"
	"github.com/govenue/osext"

	"github.com/geego/gean/app/geanlib"
	"github.com/geego/gean/app/helpers"
)

var versionCmd = &goman.Command{
	Use:   "version",
	Short: "Print the version number of Gean",
	Long:  `All software has versions. This is Gean's.`,
	RunE: func(cmd *goman.Command, args []string) error {
		printHugoVersion()
		return nil
	},
}

func printHugoVersion() {
	if geanlib.BuildDate == "" {
		setBuildDate() // set the build date from executable's mdate
	} else {
		formatBuildDate() // format the compile time
	}
	if geanlib.CommitHash == "" {
		notepad.FEEDBACK.Printf("Gean Static Site Generator v%s %s/%s BuildDate: %s\n", helpers.CurrentHugoVersion, runtime.GOOS, runtime.GOARCH, geanlib.BuildDate)
	} else {
		notepad.FEEDBACK.Printf("Gean Static Site Generator v%s-%s %s/%s BuildDate: %s\n", helpers.CurrentHugoVersion, strings.ToUpper(geanlib.CommitHash), runtime.GOOS, runtime.GOARCH, geanlib.BuildDate)
	}
}

// setBuildDate checks the ModTime of the Gean executable and returns it as a
// formatted string.  This assumes that the executable name is Gean, if it does
// not exist, an empty string will be returned.  This is only called if the
// Geanlib.BuildDate wasn't set during compile time.
//
// osext is used for cross-platform.
func setBuildDate() {
	fname, _ := osext.Executable()
	dir, err := filepath.Abs(filepath.Dir(fname))
	if err != nil {
		notepad.ERROR.Println(err)
		return
	}
	fi, err := os.Lstat(filepath.Join(dir, filepath.Base(fname)))
	if err != nil {
		notepad.ERROR.Println(err)
		return
	}
	t := fi.ModTime()
	geanlib.BuildDate = t.Format(time.RFC3339)
}

// formatBuildDate formats the geanlib.BuildDate according to the value in
// .Params.DateFormat, if it's set.
func formatBuildDate() {
	t, _ := time.Parse("2006-01-02T15:04:05-0700", geanlib.BuildDate)
	geanlib.BuildDate = t.Format(time.RFC3339)
}
