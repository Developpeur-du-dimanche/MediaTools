package mediatools

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/view"
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

func (a *Application) SetView(view view.View) {
	(*a.window).SetContent(view.Content())
}

func Run() {
	/*metadata := &fyne.AppMetadata{
		ID:   "com.github.developpeur-du-dimanche.mediatools",
		Name: "MediaTools",
	}*/

	fmt.Printf("Starting MediaTools\n")

	a := app.New()

	application := NewApplication(a)

	homeView := view.NewHomeView(a)

	application.SetView(homeView)

	homeView.ShowAndRun()

}
