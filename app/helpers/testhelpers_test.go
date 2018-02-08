package helpers

import (
	"github.com/geego/gean/app/geanfs"
	"github.com/gostores/configurator"
)

func newTestPathSpec(fs *geanfs.Fs, v *configurator.Configurator) *PathSpec {
	l := NewDefaultLanguage(v)
	ps, _ := NewPathSpec(fs, l)
	return ps
}

func newTestDefaultPathSpec(configKeyValues ...interface{}) *PathSpec {
	v := configurator.New()
	fs := geanfs.NewMem(v)
	cfg := newTestCfg(fs)

	for i := 0; i < len(configKeyValues); i += 2 {
		cfg.Set(configKeyValues[i].(string), configKeyValues[i+1])
	}
	return newTestPathSpec(fs, cfg)
}

func newTestCfg(fs *geanfs.Fs) *configurator.Configurator {
	v := configurator.New()

	v.SetFs(fs.Source)

	return v

}

func newTestContentSpec() *ContentSpec {
	v := configurator.New()
	spec, err := NewContentSpec(v)
	if err != nil {
		panic(err)
	}
	return spec
}
