package application

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type ApplicationInterface interface {
	GetApp() fyne.App
	GetWindow() fyne.Window
	GetDisplaySize() (int, int)
}

type Application struct {
	App    fyne.App
	Window *Window
}

func NewApp(metadata *fyne.AppMetadata) ApplicationInterface {
	app.SetMetadata(*metadata)
	a := app.New()

	return &Application{
		App:    a,
		Window: newWindow(a.NewWindow("mediatools")),
	}

}

func (a *Application) GetApp() fyne.App {
	return a.App
}

func (a *Application) GetWindow() fyne.Window {
	return a.Window.Window
}

func (a *Application) GetDisplaySize() (int, int) {
	return 1, 1
}
