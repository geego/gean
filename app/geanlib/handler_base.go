package geanlib

import (
	"github.com/geego/gean/app/source"
)

// Handler is used for processing files of a specific type.
type Handler interface {
	FileConvert(*source.File, *Site) HandledResult
	PageConvert(*Page) HandledResult
	Read(*source.File, *Site) HandledResult
	Extensions() []string
}

// Handle identifies functionality associated with certain file extensions.
type Handle struct {
	extensions []string
}

// Extensions returns a list of extensions.
func (h Handle) Extensions() []string {
	return h.extensions
}

// HandledResult describes the results of a file handling operation.
type HandledResult struct {
	page *Page
	file *source.File
	err  error
}

// HandledResult is an error
func (h HandledResult) Error() string {
	if h.err != nil {
		if h.page != nil {
			return "Error: " + h.err.Error() + " for " + h.page.File.LogicalName()
		}
		if h.file != nil {
			return "Error: " + h.err.Error() + " for " + h.file.LogicalName()
		}
	}
	return h.err.Error()
}

func (h HandledResult) String() string {
	return h.Error()
}

// Page returns the affected page.
func (h HandledResult) Page() *Page {
	return h.page
}
