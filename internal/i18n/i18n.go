package i18n

import (
	"embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/lang"
)

//go:embed locales/*.json
var Locales embed.FS

// Init initializes the translations
func Init() error {
	return lang.AddTranslationsFS(Locales, "locales")
}

// AddLocale loads translations for a specific locale
func AddLocale(locale string) error {
	data, err := Locales.ReadFile("locales/" + locale + ".json")
	if err != nil {
		return err
	}
	return lang.AddTranslationsForLocale(data, fyne.Locale(locale))
}
