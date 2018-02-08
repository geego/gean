package geanfs

import (
	"github.com/geego/gean/app/config"
	"github.com/gostores/fsintra"
)

// Os points to an Os fsintra file system.
var Os = &fsintra.OsFs{}

// Fs abstracts the file system to separate source and destination file systems
// and allows both to be mocked for testing.
type Fs struct {
	// Source is Hugo's source file system.
	Source fsintra.Fs

	// Destination is Hugo's destination file system.
	Destination fsintra.Fs

	// Os is an OS file system.
	// NOTE: Field is currently unused.
	Os fsintra.Fs

	// WorkingDir is a read-only file system
	// restricted to the project working dir.
	WorkingDir *fsintra.BasePathFs
}

// NewDefault creates a new Fs with the OS file system
// as source and destination file systems.
func NewDefault(cfg config.Provider) *Fs {
	fs := &fsintra.OsFs{}
	return newFs(fs, cfg)
}

// NewMem creates a new Fs with the MemMapFs
// as source and destination file systems.
// Useful for testing.
func NewMem(cfg config.Provider) *Fs {
	fs := &fsintra.MemMapFs{}
	return newFs(fs, cfg)
}

// NewFrom creates a new Fs based on the provided fsintra Fs
// as source and destination file systems.
// Useful for testing.
func NewFrom(fs fsintra.Fs, cfg config.Provider) *Fs {
	return newFs(fs, cfg)
}

func newFs(base fsintra.Fs, cfg config.Provider) *Fs {
	return &Fs{
		Source:      base,
		Destination: base,
		Os:          &fsintra.OsFs{},
		WorkingDir:  getWorkingDirFs(base, cfg),
	}
}

func getWorkingDirFs(base fsintra.Fs, cfg config.Provider) *fsintra.BasePathFs {
	workingDir := cfg.GetString("workingDir")

	if workingDir != "" {
		return fsintra.NewBasePathFs(fsintra.NewReadOnlyFs(base), workingDir).(*fsintra.BasePathFs)
	}

	return nil
}
