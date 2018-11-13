package geanlib

import (
	"github.com/govenue/assist"
	"github.com/govenue/notepad"
)

// Sitemap configures the sitemap to be generated.
type Sitemap struct {
	ChangeFreq string
	Priority   float64
	Filename   string
}

func parseSitemap(input map[string]interface{}) Sitemap {
	sitemap := Sitemap{Priority: -1, Filename: "sitemap.xml"}

	for key, value := range input {
		switch key {
		case "changefreq":
			sitemap.ChangeFreq = assist.ToString(value)
		case "priority":
			sitemap.Priority = assist.ToFloat64(value)
		case "filename":
			sitemap.Filename = assist.ToString(value)
		default:
			notepad.WARN.Printf("Unknown Sitemap field: %s\n", key)
		}
	}

	return sitemap
}
