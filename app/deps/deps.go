package deps

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/geego/gean/app/config"
	"github.com/geego/gean/app/geanfs"
	"github.com/geego/gean/app/helpers"
	"github.com/geego/gean/app/metrics"
	"github.com/geego/gean/app/output"
	"github.com/geego/gean/app/tpl"
	"github.com/gostores/notepad"
)

// Deps holds dependencies used by many.
// There will be normally only one instance of deps in play
// at a given time, i.e. one per Site built.
type Deps struct {
	// The logger to use.
	Log *notepad.Notepad `json:"-"`

	// The templates to use. This will usually implement the full tpl.TemplateHandler.
	Tmpl tpl.TemplateFinder `json:"-"`

	// The file systems to use.
	Fs *geanfs.Fs `json:"-"`

	// The PathSpec to use
	*helpers.PathSpec `json:"-"`

	// The ContentSpec to use
	*helpers.ContentSpec `json:"-"`

	// The configuration to use
	Cfg config.Provider `json:"-"`

	// The translation func to use
	Translate func(translationID string, args ...interface{}) string `json:"-"`

	Language *helpers.Language

	// All the output formats available for the current site.
	OutputFormatsConfig output.Formats

	templateProvider ResourceProvider
	WithTemplate     func(templ tpl.TemplateHandler) error `json:"-"`

	translationProvider ResourceProvider

	Metrics metrics.Provider
}

// ResourceProvider is used to create and refresh, and clone resources needed.
type ResourceProvider interface {
	Update(deps *Deps) error
	Clone(deps *Deps) error
}

// TemplateHandler returns the used tpl.TemplateFinder as tpl.TemplateHandler.
func (d *Deps) TemplateHandler() tpl.TemplateHandler {
	return d.Tmpl.(tpl.TemplateHandler)
}

// LoadResources loads translations and templates.
func (d *Deps) LoadResources() error {
	// Note that the translations need to be loaded before the templates.
	if err := d.translationProvider.Update(d); err != nil {
		return err
	}

	if err := d.templateProvider.Update(d); err != nil {
		return err
	}

	if th, ok := d.Tmpl.(tpl.TemplateHandler); ok {
		th.PrintErrors()
	}

	return nil
}

// New initializes a Dep struct.
// Defaults are set for nil values,
// but TemplateProvider, TranslationProvider and Language are always required.
func New(cfg DepsCfg) (*Deps, error) {
	var (
		logger = cfg.Logger
		fs     = cfg.Fs
	)

	if cfg.TemplateProvider == nil {
		panic("Must have a TemplateProvider")
	}

	if cfg.TranslationProvider == nil {
		panic("Must have a TranslationProvider")
	}

	if cfg.Language == nil {
		panic("Must have a Language")
	}

	if logger == nil {
		logger = notepad.NewNotepad(notepad.LevelError, notepad.LevelError, os.Stdout, ioutil.Discard, "", log.Ldate|log.Ltime)
	}

	if fs == nil {
		// Default to the production file system.
		fs = geanfs.NewDefault(cfg.Language)
	}

	ps, err := helpers.NewPathSpec(fs, cfg.Language)

	if err != nil {
		return nil, err
	}

	contentSpec, err := helpers.NewContentSpec(cfg.Language)
	if err != nil {
		return nil, err
	}

	d := &Deps{
		Fs:                  fs,
		Log:                 logger,
		templateProvider:    cfg.TemplateProvider,
		translationProvider: cfg.TranslationProvider,
		WithTemplate:        cfg.WithTemplate,
		PathSpec:            ps,
		ContentSpec:         contentSpec,
		Cfg:                 cfg.Language,
		Language:            cfg.Language,
	}

	if cfg.Cfg.GetBool("templateMetrics") {
		d.Metrics = metrics.NewProvider(cfg.Cfg.GetBool("templateMetricsHints"))
	}

	return d, nil
}

// ForLanguage creates a copy of the Deps with the language dependent
// parts switched out.
func (d Deps) ForLanguage(l *helpers.Language) (*Deps, error) {
	var err error

	d.PathSpec, err = helpers.NewPathSpec(d.Fs, l)
	if err != nil {
		return nil, err
	}

	d.ContentSpec, err = helpers.NewContentSpec(l)
	if err != nil {
		return nil, err
	}

	d.Cfg = l
	d.Language = l

	if err := d.translationProvider.Clone(&d); err != nil {
		return nil, err
	}

	if err := d.templateProvider.Clone(&d); err != nil {
		return nil, err
	}

	return &d, nil

}

// DepsCfg contains configuration options that can be used to configure Hugo
// on a global level, i.e. logging etc.
// Nil values will be given default values.
type DepsCfg struct {

	// The Logger to use.
	Logger *notepad.Notepad

	// The file systems to use
	Fs *geanfs.Fs

	// The language to use.
	Language *helpers.Language

	// The configuration to use.
	Cfg config.Provider

	// Template handling.
	TemplateProvider ResourceProvider
	WithTemplate     func(templ tpl.TemplateHandler) error

	// i18n handling.
	TranslationProvider ResourceProvider
}
