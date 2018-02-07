package fmt

import (
	_fmt "fmt"
	"yiqilai.tech/gean/app/helpers"
)

// New returns a new instance of the fmt-namespaced template functions.
func New() *Namespace {
	return &Namespace{helpers.NewDistinctErrorLogger()}
}

// Namespace provides template functions for the "fmt" namespace.
type Namespace struct {
	errorLogger *helpers.DistinctLogger
}

// Print returns string representation of the passed arguments.
func (ns *Namespace) Print(a ...interface{}) string {
	return _fmt.Sprint(a...)
}

// Printf returns a formatted string representation of the passed arguments.
func (ns *Namespace) Printf(format string, a ...interface{}) string {
	return _fmt.Sprintf(format, a...)

}

// Println returns string representation of the passed arguments ending with a newline.
func (ns *Namespace) Println(a ...interface{}) string {
	return _fmt.Sprintln(a...)
}

func (ns *Namespace) Errorf(format string, a ...interface{}) string {
	ns.errorLogger.Printf(format, a...)
	return _fmt.Sprintf(format, a...)
}
