package output

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/geego/gean/app/media"
	"github.com/govenue/mapstructure"
)

// Format represents an output representation, usually to a file on disk.
type Format struct {
	// The Name is used as an identifier. Internal output formats (i.e. HTML and RSS)
	// can be overridden by providing a new definition for those types.
	Name string `json:"name"`

	MediaType media.Type `json:"mediaType"`

	// Must be set to a value when there are two or more conflicting mediatype for the same resource.
	Path string `json:"path"`

	// The base output file name used when not using "ugly URLs", defaults to "index".
	BaseName string `json:"baseName"`

	// The value to use for rel links
	//
	// See https://www.w3schools.com/tags/att_link_rel.asp
	//
	// AMP has a special requirement in this department, see:
	// https://www.ampproject.org/docs/guides/deploy/discovery
	// I.e.:
	// <link rel="amphtml" href="https://www.example.com/url/to/amp/document.html">
	Rel string `json:"rel"`

	// The protocol to use, i.e. "webcal://". Defaults to the protocol of the baseURL.
	Protocol string `json:"protocol"`

	// IsPlainText decides whether to use text/template or html/template
	// as template parser.
	IsPlainText bool `json:"isPlainText"`

	// IsHTML returns whether this format is int the HTML family. This includes
	// HTML, AMP etc. This is used to decide when to create alias redirects etc.
	IsHTML bool `json:"isHTML"`

	// Enable to ignore the global uglyURLs setting.
	NoUgly bool `json:"noUgly"`

	// Enable if it doesn't make sense to include this format in an alternative
	// format listing, CSS being one good example.
	// Note that we use the term "alternative" and not "alternate" here, as it
	// does not necessarily replace the other format, it is an alternative representation.
	NotAlternative bool `json:"notAlternative"`
}

var (
	// An ordered list of built-in output formats
	//
	// See https://www.ampproject.org/learn/overview/
	AMPFormat = Format{
		Name:      "AMP",
		MediaType: media.HTMLType,
		BaseName:  "index",
		Path:      "amp",
		Rel:       "amphtml",
		IsHTML:    true,
	}

	CalendarFormat = Format{
		Name:        "Calendar",
		MediaType:   media.CalendarType,
		IsPlainText: true,
		Protocol:    "webcal://",
		BaseName:    "index",
		Rel:         "alternate",
	}

	CSSFormat = Format{
		Name:           "CSS",
		MediaType:      media.CSSType,
		BaseName:       "styles",
		IsPlainText:    true,
		Rel:            "stylesheet",
		NotAlternative: true,
	}
	CSVFormat = Format{
		Name:        "CSV",
		MediaType:   media.CSVType,
		BaseName:    "index",
		IsPlainText: true,
		Rel:         "alternate",
	}

	HTMLFormat = Format{
		Name:      "HTML",
		MediaType: media.HTMLType,
		BaseName:  "index",
		Rel:       "canonical",
		IsHTML:    true,
	}

	JSONFormat = Format{
		Name:        "JSON",
		MediaType:   media.JSONType,
		BaseName:    "index",
		IsPlainText: true,
		Rel:         "alternate",
	}

	RSSFormat = Format{
		Name:      "RSS",
		MediaType: media.RSSType,
		BaseName:  "index",
		NoUgly:    true,
		Rel:       "alternate",
	}
)

var DefaultFormats = Formats{
	AMPFormat,
	CalendarFormat,
	CSSFormat,
	CSVFormat,
	HTMLFormat,
	JSONFormat,
	RSSFormat,
}

func init() {
	sort.Sort(DefaultFormats)
}

type Formats []Format

func (formats Formats) Len() int           { return len(formats) }
func (formats Formats) Swap(i, j int)      { formats[i], formats[j] = formats[j], formats[i] }
func (formats Formats) Less(i, j int) bool { return formats[i].Name < formats[j].Name }

// GetBySuffix gets a output format given as suffix, e.g. "html".
// It will return false if no format could be found, or if the suffix given
// is ambiguous.
// The lookup is case insensitive.
func (formats Formats) GetBySuffix(suffix string) (f Format, found bool) {
	for _, ff := range formats {
		if strings.EqualFold(suffix, ff.MediaType.Suffix) {
			if found {
				// ambiguous
				found = false
				return
			}
			f = ff
			found = true
		}
	}
	return
}

// GetByName gets a format by its identifier name.
func (formats Formats) GetByName(name string) (f Format, found bool) {
	for _, ff := range formats {
		if strings.EqualFold(name, ff.Name) {
			f = ff
			found = true
			return
		}
	}
	return
}

// GetByNames gets a list of formats given a list of identifiers.
func (formats Formats) GetByNames(names ...string) (Formats, error) {
	var types []Format

	for _, name := range names {
		tpe, ok := formats.GetByName(name)
		if !ok {
			return types, fmt.Errorf("OutputFormat with key %q not found", name)
		}
		types = append(types, tpe)
	}
	return types, nil
}

// FromFilename gets a Format given a filename.
func (formats Formats) FromFilename(filename string) (f Format, found bool) {
	// mytemplate.amp.html
	// mytemplate.html
	// mytemplate
	var ext, outFormat string

	parts := strings.Split(filename, ".")
	if len(parts) > 2 {
		outFormat = parts[1]
		ext = parts[2]
	} else if len(parts) > 1 {
		ext = parts[1]
	}

	if outFormat != "" {
		return formats.GetByName(outFormat)
	}

	if ext != "" {
		f, found = formats.GetBySuffix(ext)
		if !found && len(parts) == 2 {
			// For extensionless output formats (e.g. Netlify's _redirects)
			// we must fall back to using the extension as format lookup.
			f, found = formats.GetByName(ext)
		}
	}
	return
}

// DecodeFormats takes a list of output format configurations and merges those,
// in the order given, with the Hugo defaults as the last resort.
func DecodeFormats(mediaTypes media.Types, maps ...map[string]interface{}) (Formats, error) {
	f := make(Formats, len(DefaultFormats))
	copy(f, DefaultFormats)

	for _, m := range maps {
		for k, v := range m {
			found := false
			for i, vv := range f {
				if strings.EqualFold(k, vv.Name) {
					// Merge it with the existing
					if err := decode(mediaTypes, v, &f[i]); err != nil {
						return f, err
					}
					found = true
				}
			}
			if !found {
				var newOutFormat Format
				newOutFormat.Name = k
				if err := decode(mediaTypes, v, &newOutFormat); err != nil {
					return f, err
				}

				// We need values for these
				if newOutFormat.BaseName == "" {
					newOutFormat.BaseName = "index"
				}
				if newOutFormat.Rel == "" {
					newOutFormat.Rel = "alternate"
				}

				f = append(f, newOutFormat)
			}
		}
	}

	sort.Sort(f)

	return f, nil
}

func decode(mediaTypes media.Types, input, output interface{}) error {
	config := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           output,
		WeaklyTypedInput: true,
		DecodeHook: func(a reflect.Type, b reflect.Type, c interface{}) (interface{}, error) {
			if a.Kind() == reflect.Map {
				dataVal := reflect.Indirect(reflect.ValueOf(c))
				for _, key := range dataVal.MapKeys() {
					keyStr, ok := key.Interface().(string)
					if !ok {
						// Not a string key
						continue
					}
					if strings.EqualFold(keyStr, "mediaType") {
						// If mediaType is a string, look it up and replace it
						// in the map.
						vv := dataVal.MapIndex(key)
						if mediaTypeStr, ok := vv.Interface().(string); ok {
							mediaType, found := mediaTypes.GetByType(mediaTypeStr)
							if !found {
								return c, fmt.Errorf("media type %q not found", mediaTypeStr)
							}
							dataVal.SetMapIndex(key, reflect.ValueOf(mediaType))
						}
					}
				}
			}
			return c, nil
		},
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}

func (formats Format) BaseFilename() string {
	return formats.BaseName + "." + formats.MediaType.Suffix
}

func (formats Format) MarshalJSON() ([]byte, error) {
	type Alias Format
	return json.Marshal(&struct {
		MediaType string
		Alias
	}{
		MediaType: formats.MediaType.String(),
		Alias:     (Alias)(formats),
	})
}
