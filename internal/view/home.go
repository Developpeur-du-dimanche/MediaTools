package view

import (
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/application"
)

type HomeView struct {
	app application.ApplicationInterface
}

func NewHomeView(app application.ApplicationInterface) *HomeView {

	return &HomeView{
		app: app,
	}
}

func (v *HomeView) LogGui() {
	v.app.GetWindow().SetContent(widget.NewLabel("Hello, World!"))
	v.app.GetWindow().ShowAndRun()
}
