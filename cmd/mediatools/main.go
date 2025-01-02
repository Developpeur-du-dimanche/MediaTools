package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/lang"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/mediatools"
)

func main() {
	a := app.New()
	application := mediatools.NewApplication(a)
	lang.AddTranslationsFS(mediatools.Translations, "localize")
	application.Run()
}
