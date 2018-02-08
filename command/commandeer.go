package command

import (
	"github.com/geego/gean/app/common/types"
	"github.com/geego/gean/app/deps"
	"github.com/geego/gean/app/geanfs"
	"github.com/geego/gean/app/helpers"
)

type commandeer struct {
	*deps.DepsCfg
	pathSpec    *helpers.PathSpec
	visitedURLs *types.EvictingStringQueue

	serverPorts []int

	configured bool
}

func (c *commandeer) Set(key string, value interface{}) {
	if c.configured {
		panic("commandeer cannot be changed")
	}
	c.Cfg.Set(key, value)
}

// PathSpec lazily creates a new PathSpec, as all the paths must
// be configured before it is created.
func (c *commandeer) PathSpec() *helpers.PathSpec {
	c.configured = true
	return c.pathSpec
}

func (c *commandeer) languages() helpers.Languages {
	return c.Cfg.Get("languagesSorted").(helpers.Languages)
}

func (c *commandeer) initFs(fs *geanfs.Fs) error {
	c.DepsCfg.Fs = fs
	ps, err := helpers.NewPathSpec(fs, c.Cfg)
	if err != nil {
		return err
	}
	c.pathSpec = ps
	return nil
}

func newCommandeer(cfg *deps.DepsCfg) (*commandeer, error) {
	l := cfg.Language
	if l == nil {
		l = helpers.NewDefaultLanguage(cfg.Cfg)
	}
	ps, err := helpers.NewPathSpec(cfg.Fs, l)
	if err != nil {
		return nil, err
	}

	return &commandeer{DepsCfg: cfg, pathSpec: ps, visitedURLs: types.NewEvictingStringQueue(10)}, nil
}
