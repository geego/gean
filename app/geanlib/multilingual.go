package geanlib

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/geego/gean/app/config"
	"github.com/geego/gean/app/helpers"
	"github.com/gostores/assist"
)

// Multilingual manages the all languages used in a multilingual site.
type Multilingual struct {
	Languages helpers.Languages

	DefaultLang *helpers.Language

	langMap     map[string]*helpers.Language
	langMapInit sync.Once
}

// Language returns the Language associated with the given string.
func (ml *Multilingual) Language(lang string) *helpers.Language {
	ml.langMapInit.Do(func() {
		ml.langMap = make(map[string]*helpers.Language)
		for _, l := range ml.Languages {
			ml.langMap[l.Lang] = l
		}
	})
	return ml.langMap[lang]
}

func getLanguages(cfg config.Provider) helpers.Languages {
	if cfg.IsSet("languagesSorted") {
		return cfg.Get("languagesSorted").(helpers.Languages)
	}

	return helpers.Languages{helpers.NewDefaultLanguage(cfg)}
}

func newMultiLingualFromSites(cfg config.Provider, sites ...*Site) (*Multilingual, error) {
	languages := make(helpers.Languages, len(sites))

	for i, s := range sites {
		if s.Language == nil {
			return nil, errors.New("Missing language for site")
		}
		languages[i] = s.Language
	}

	defaultLang := cfg.GetString("defaultContentLanguage")

	if defaultLang == "" {
		defaultLang = "en"
	}

	return &Multilingual{Languages: languages, DefaultLang: helpers.NewLanguage(defaultLang, cfg)}, nil

}

func newMultiLingualForLanguage(language *helpers.Language) *Multilingual {
	languages := helpers.Languages{language}
	return &Multilingual{Languages: languages, DefaultLang: language}
}
func (ml *Multilingual) enabled() bool {
	return len(ml.Languages) > 1
}

func (s *Site) multilingualEnabled() bool {
	if s.owner == nil {
		return false
	}
	return s.owner.multilingual != nil && s.owner.multilingual.enabled()
}

func toSortedLanguages(cfg config.Provider, l map[string]interface{}) (helpers.Languages, error) {
	langs := make(helpers.Languages, len(l))
	i := 0

	for lang, langConf := range l {
		langsMap, err := assist.ToStringMapE(langConf)

		if err != nil {
			return nil, fmt.Errorf("Language config is not a map: %T", langConf)
		}

		language := helpers.NewLanguage(lang, cfg)

		for loki, v := range langsMap {
			switch loki {
			case "title":
				language.Title = assist.ToString(v)
			case "languagename":
				language.LanguageName = assist.ToString(v)
			case "weight":
				language.Weight = assist.ToInt(v)
			}

			// Put all into the Params map
			language.SetParam(loki, v)
		}

		langs[i] = language
		i++
	}

	sort.Sort(langs)

	return langs, nil
}
