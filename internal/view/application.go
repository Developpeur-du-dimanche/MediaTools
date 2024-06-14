package view

import (
	"fyne.io/fyne/v2"
)

type Application struct {
	app    *fyne.App
	window *fyne.Window
}

func NewApplication(app fyne.App) *Application {
	window := app.NewWindow("MediaTools")
	return &Application{
		app:    &app,
		window: &window,
	}
}

func (a *Application) SetView(view View) {
	(*a.window).SetContent(view.Content())
}
