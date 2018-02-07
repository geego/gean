package main

import (
	"os"
	"runtime"

	"github.com/gostores/notepad"

	"yiqilai.tech/gean/command"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	command.Execute()

	if notepad.LogCountForLevelsGreaterThanorEqualTo(notepad.LevelError) > 0 {
		os.Exit(-1)
	}

	if command.Hugo != nil {
		if command.Hugo.Log.LogCountForLevelsGreaterThanorEqualTo(notepad.LevelError) > 0 {
			os.Exit(-1)
		}
	}
}
