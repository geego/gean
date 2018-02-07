package geanlib

import (
	"bytes"

	"github.com/gostores/cssmin"
	"yiqilai.tech/gean/app/source"
)

func init() {
	RegisterHandler(new(cssHandler))
	RegisterHandler(new(defaultHandler))
}

type basicFileHandler Handle

func (h basicFileHandler) Read(f *source.File, s *Site) HandledResult {
	return HandledResult{file: f}
}

func (h basicFileHandler) PageConvert(*Page) HandledResult {
	return HandledResult{}
}

type defaultHandler struct{ basicFileHandler }

func (h defaultHandler) Extensions() []string { return []string{"*"} }
func (h defaultHandler) FileConvert(f *source.File, s *Site) HandledResult {
	err := s.publish(f.Path(), f.Contents)
	if err != nil {
		return HandledResult{err: err}
	}
	return HandledResult{file: f}
}

type cssHandler struct{ basicFileHandler }

func (h cssHandler) Extensions() []string { return []string{"css"} }
func (h cssHandler) FileConvert(f *source.File, s *Site) HandledResult {
	x := cssmin.Minify(f.Bytes())
	err := s.publish(f.Path(), bytes.NewReader(x))
	if err != nil {
		return HandledResult{err: err}
	}
	return HandledResult{file: f}
}
