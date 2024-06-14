package components

import (
	"fyne.io/fyne/v2"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/application"
)

type MenuComponent struct {
	fyne.Menu
	app application.ApplicationInterface
}

func NewMenuComponent(app application.ApplicationInterface) *MenuComponent {
	return &MenuComponent{
		app: app,
	}
}
