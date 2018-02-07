package geanlib

// Translations represent the other translations for a given page. The
// string here is the language code, as affected by the `post.LANG.md`
// filename.
type Translations map[string]*Page

func pagesToTranslationsMap(pages []*Page) map[string]Translations {
	out := make(map[string]Translations)

	for _, page := range pages {
		base := page.TranslationKey()

		pageTranslation, present := out[base]
		if !present {
			pageTranslation = make(Translations)
		}

		pageLang := page.Lang()
		if pageLang == "" {
			continue
		}

		pageTranslation[pageLang] = page
		out[base] = pageTranslation
	}

	return out
}

func assignTranslationsToPages(allTranslations map[string]Translations, pages []*Page) {
	for _, page := range pages {
		page.translations = page.translations[:0]
		base := page.TranslationKey()
		trans, exist := allTranslations[base]
		if !exist {
			continue
		}

		for _, translatedPage := range trans {
			page.translations = append(page.translations, translatedPage)
		}

		pageBy(languagePageSort).Sort(page.translations)
	}
}
