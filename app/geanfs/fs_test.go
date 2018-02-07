package geanfs

import (
	"testing"

	"github.com/gostores/assert"
	"github.com/gostores/configurator"
	"github.com/gostores/fsintra"
)

func TestNewDefault(t *testing.T) {
	v := configurator.New()
	f := NewDefault(v)

	assert.NotNil(t, f.Source)
	assert.IsType(t, new(fsintra.OsFs), f.Source)
	assert.NotNil(t, f.Destination)
	assert.IsType(t, new(fsintra.OsFs), f.Destination)
	assert.NotNil(t, f.Os)
	assert.IsType(t, new(fsintra.OsFs), f.Os)
	assert.Nil(t, f.WorkingDir)

	assert.IsType(t, new(fsintra.OsFs), Os)
}

func TestNewMem(t *testing.T) {
	v := configurator.New()
	f := NewMem(v)

	assert.NotNil(t, f.Source)
	assert.IsType(t, new(fsintra.MemMapFs), f.Source)
	assert.NotNil(t, f.Destination)
	assert.IsType(t, new(fsintra.MemMapFs), f.Destination)
	assert.IsType(t, new(fsintra.OsFs), f.Os)
	assert.Nil(t, f.WorkingDir)
}

func TestWorkingDir(t *testing.T) {
	v := configurator.New()

	v.Set("workingDir", "/a/b/")

	f := NewMem(v)

	assert.NotNil(t, f.WorkingDir)
	assert.IsType(t, new(fsintra.BasePathFs), f.WorkingDir)
}
