package command

import (
	"runtime"

	"github.com/govenue/goman"
	"github.com/govenue/notepad"
)

var envCmd = &goman.Command{
	Use:   "env",
	Short: "Print Gean version and environment info",
	Long:  `Print Gean version and environment info. This is useful in Gean bug reports.`,
	RunE: func(cmd *goman.Command, args []string) error {
		printHugoVersion()
		notepad.FEEDBACK.Printf("GOOS=%q\n", runtime.GOOS)
		notepad.FEEDBACK.Printf("GOARCH=%q\n", runtime.GOARCH)
		notepad.FEEDBACK.Printf("GOVERSION=%q\n", runtime.Version())

		return nil
	},
}
