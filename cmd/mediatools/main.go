package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()

	a.Settings().SetTheme(newMediaToolsTheme())

	w := a.NewWindow("Media Tools")
	w.Resize(fyne.NewSize(800, 600))

	w.SetContent(container.NewVBox(
		widget.NewButton("Button", func() {}),
		widget.NewButton("Disabled", func() {}),
	))

	w.ShowAndRun()
}
