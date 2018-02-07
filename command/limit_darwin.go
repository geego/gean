package command

import (
	"syscall"

	"github.com/gostores/goman"
	"github.com/gostores/notepad"
)

func init() {
	commandCheck.AddCommand(limit)
}

var limit = &goman.Command{
	Use:   "ulimit",
	Short: "Check system ulimit settings",
	Long: `Gean will inspect the current ulimit settings on the system.
This is primarily to ensure that Gean can watch enough files on some OSs`,
	RunE: func(cmd *goman.Command, args []string) error {
		var rLimit syscall.Rlimit
		err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
		if err != nil {
			return newSystemError("Error Getting Rlimit ", err)
		}

		notepad.FEEDBACK.Println("Current rLimit:", rLimit)

		notepad.FEEDBACK.Println("Attempting to increase limit")
		rLimit.Max = 999999
		rLimit.Cur = 999999
		err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
		if err != nil {
			return newSystemError("Error Setting rLimit ", err)
		}
		err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
		if err != nil {
			return newSystemError("Error Getting rLimit ", err)
		}
		notepad.FEEDBACK.Println("rLimit after change:", rLimit)

		return nil
	},
}

func tweakLimit() {
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		notepad.ERROR.Println("Unable to obtain rLimit", err)
	}
	if rLimit.Cur < rLimit.Max {
		rLimit.Max = 64000
		rLimit.Cur = 64000
		err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
		if err != nil {
			notepad.WARN.Println("Unable to increase number of open files limit", err)
		}
	}
}
