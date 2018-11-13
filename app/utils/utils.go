package utils

import (
	"os"

	"github.com/govenue/notepad"
)

// CheckErr logs the messages given and then the error.
// TODO(bep) Remove this package.
func CheckErr(logger *notepad.Notepad, err error, s ...string) {
	if err == nil {
		return
	}
	if len(s) == 0 {
		logger.CRITICAL.Println(err)
		return
	}
	for _, message := range s {
		logger.ERROR.Println(message)
	}
	logger.ERROR.Println(err)
}

// StopOnErr exits on any error after logging it.
func StopOnErr(logger *notepad.Notepad, err error, s ...string) {
	if err == nil {
		return
	}

	defer os.Exit(-1)

	if len(s) == 0 {
		newMessage := err.Error()
		// Printing an empty string results in a error with
		// no message, no bueno.
		if newMessage != "" {
			logger.CRITICAL.Println(newMessage)
		}
	}
	for _, message := range s {
		if message != "" {
			logger.CRITICAL.Println(message)
		}
	}
}
