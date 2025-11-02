package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/i18n"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/mediatools"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/theme"
)

func main() {

	a := app.NewWithID("com.TOomaAh.mediatools")
	app.SetMetadata(fyne.AppMetadata{
		Name:    "MediaTools",
		Version: "0.1",
	})

	// Initialize translations
	if err := i18n.Init(); err != nil {
		panic(err)
	}

	// Get saved language preference and load it
	savedLang := a.Preferences().StringWithFallback("language", "")
	if savedLang != "" {
		i18n.AddLocale(savedLang)
	}

	a.Settings().SetTheme(theme.NewMediaToolsTheme())

	mt := mediatools.NewMediaTools(a)
	mt.Run()

}
