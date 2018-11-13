package helpers

import (
	"sort"
	"strings"
	"sync"

	"github.com/geego/gean/app/config"
	"github.com/govenue/assist"
)

// These are the settings that should only be looked up in the global Viper
// config and not per language.
// This list may not be complete, but contains only settings that we know
// will be looked up in both.
// This isn't perfect, but it is ultimately the user who shoots him/herself in
// the foot.
// See the pathSpec.
var globalOnlySettings = map[string]bool{
	strings.ToLower("defaultContentLanguageInSubdir"): true,
	strings.ToLower("defaultContentLanguage"):         true,
	strings.ToLower("multilingual"):                   true,
}

// Language manages specific-language configuration.
type Language struct {
	Lang         string
	LanguageName string
	Title        string
	Weight       int

	Cfg        config.Provider
	params     map[string]interface{}
	paramsInit sync.Once
}

func (l *Language) String() string {
	return l.Lang
}

// NewLanguage creates a new language.
func NewLanguage(lang string, cfg config.Provider) *Language {
	return &Language{Lang: lang, Cfg: cfg, params: make(map[string]interface{})}
}

// NewDefaultLanguage creates the default language for a config.Provider.
// If not otherwise specified the default is "en".
func NewDefaultLanguage(cfg config.Provider) *Language {
	defaultLang := cfg.GetString("defaultContentLanguage")

	if defaultLang == "" {
		defaultLang = "en"
	}

	return NewLanguage(defaultLang, cfg)
}

// Languages is a sortable list of languages.
type Languages []*Language

// NewLanguages creates a sorted list of languages.
// NOTE: function is currently unused.
func NewLanguages(l ...*Language) Languages {
	languages := make(Languages, len(l))
	for i := 0; i < len(l); i++ {
		languages[i] = l[i]
	}
	sort.Sort(languages)
	return languages
}

func (l Languages) Len() int           { return len(l) }
func (l Languages) Less(i, j int) bool { return l[i].Weight < l[j].Weight }
func (l Languages) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

// Params retunrs language-specific params merged with the global params.
func (l *Language) Params() map[string]interface{} {
	l.paramsInit.Do(func() {
		// Merge with global config.
		// TODO(bep) consider making this part of a constructor func.

		globalParams := l.Cfg.GetStringMap("params")
		for k, v := range globalParams {
			if _, ok := l.params[k]; !ok {
				l.params[k] = v
			}
		}
	})
	return l.params
}

// IsMultihost returns whether the languages has baseURL specificed on the
// language level.
func (l Languages) IsMultihost() bool {
	for _, lang := range l {
		if lang.GetLocal("baseURL") != nil {
			return true
		}
	}
	return false
}

// SetParam sets param with the given key and value.
// SetParam is case-insensitive.
func (l *Language) SetParam(k string, v interface{}) {
	l.params[strings.ToLower(k)] = v
}

// GetBool returns the value associated with the key as a boolean.
func (l *Language) GetBool(key string) bool { return assist.ToBool(l.Get(key)) }

// GetString returns the value associated with the key as a string.
func (l *Language) GetString(key string) string { return assist.ToString(l.Get(key)) }

// GetInt returns the value associated with the key as an int.
func (l *Language) GetInt(key string) int { return assist.ToInt(l.Get(key)) }

// GetStringMap returns the value associated with the key as a map of interfaces.
func (l *Language) GetStringMap(key string) map[string]interface{} {
	return assist.ToStringMap(l.Get(key))
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (l *Language) GetStringMapString(key string) map[string]string {
	return assist.ToStringMapString(l.Get(key))
}

// Get returns a value associated with the key relying on specified language.
// Get is case-insensitive for a key.
//
// Get returns an interface. For a specific value use one of the Get____ methods.
func (l *Language) Get(key string) interface{} {
	local := l.GetLocal(key)
	if local != nil {
		return local
	}
	return l.Cfg.Get(key)
}

// GetLocal gets a configuration value set on language level. It will
// not fall back to any global value.
// It will return nil if a value with the given key cannot be found.
func (l *Language) GetLocal(key string) interface{} {
	if l == nil {
		panic("language not set")
	}
	key = strings.ToLower(key)
	if !globalOnlySettings[key] {
		if v, ok := l.params[key]; ok {
			return v
		}
	}
	return nil
}

// Set sets the value for the key in the language's params.
func (l *Language) Set(key string, value interface{}) {
	if l == nil {
		panic("language not set")
	}
	key = strings.ToLower(key)
	l.params[key] = value
}

// IsSet checks whether the key is set in the language or the related config store.
func (l *Language) IsSet(key string) bool {
	key = strings.ToLower(key)

	key = strings.ToLower(key)
	if !globalOnlySettings[key] {
		if _, ok := l.params[key]; ok {
			return true
		}
	}
	return l.Cfg.IsSet(key)

}
