package mediatools

import (
	"fyne.io/fyne/v2"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/configuration"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/view"
)

type Application struct {
	app           fyne.App
	configuration configuration.Configuration
}

func NewApplication(app fyne.App, configuration configuration.Configuration) *Application {
	return &Application{
		app:           app,
		configuration: configuration,
	}
}

func (a *Application) Run() {
	homeView := view.NewHomeView(a.app, a.configuration)
	homeView.ShowAndRun()

}
