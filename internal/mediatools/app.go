package mediatools

import (
	"fyne.io/fyne/v2"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/configuration"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/view"
)

type Application struct {
	app           fyne.App
	window        fyne.Window
	configuration configuration.Configuration
}

func NewApplication(app fyne.App, configuration configuration.Configuration, window *fyne.Window) *Application {
	return &Application{
		app:           app,
		window:        *window,
		configuration: configuration,
	}
}

func (a *Application) SetView(view view.View) {
	a.window.SetContent(view.Content())
}

func (a *Application) Run() {
	homeView := view.NewHomeView(&a.window, a.configuration)
	a.SetView(homeView)
	homeView.ShowAndRun()

}
