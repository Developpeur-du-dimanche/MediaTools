package mediatools

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/view"
)

type Application struct {
	app    fyne.App
	window fyne.Window
}

func NewApplication(app fyne.App) *Application {
	return &Application{
		app:    app,
		window: app.NewWindow("MediaTools"),
	}
}

func (a *Application) SetView(view view.View) {
	a.window.SetContent(view.Content())
}

func (a *Application) Run() {
	app.SetMetadata(fyne.AppMetadata{
		ID:   "com.github.developpeur-du-dimanche.mediatools",
		Name: "MediaTools",
		Icon: fyne.NewStaticResource("Icon.svg", []byte(logo)),
	})
	homeView := view.NewHomeView(a.window)
	a.SetView(homeView)
	homeView.ShowAndRun()

}
