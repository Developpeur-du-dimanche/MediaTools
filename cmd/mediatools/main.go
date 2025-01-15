package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/mediatools"
)

func main() {

	a := app.NewWithID("com.TOomaAh.mediatools")
	app.SetMetadata(fyne.AppMetadata{
		Name:    "MediaTools",
		Version: "0.1",
	})

	a.Settings().SetTheme(newMediaToolsTheme())

	mt := mediatools.NewMediaTools(a)
	mt.Run()

}
