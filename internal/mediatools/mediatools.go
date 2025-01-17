package mediatools

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/components"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/list"
	"github.com/kbinani/screenshot"
)

type MediaTools struct {
	app    *fyne.App
	window *fyne.Window
}

func NewMediaTools(app fyne.App) *MediaTools {

	w := app.NewWindow("MediaTools")

	screen := screenshot.GetDisplayBounds(0)
	w.Resize(fyne.NewSize(
		float32(screen.Dx()/2), float32(screen.Dy()/2),
	))

	items := list.NewList[string]()

	listView := components.NewListView(items, nil)

	history := components.NewLastScanSelector(func(path string) {
		fmt.Printf("Folder selected: %s\n", path)
	})

	burgerMenu := components.NewBurgerMenu(
		container.NewHBox(
			components.NewOpenFolder(&w, func(path string) {
				history.AddFolder(path)
			}, func(path string) {
				listView.AddItem(path)
				listView.Refresh()
			}),
			components.NewOpenFile(&w, func(path string) {
				listView.AddItem(path)
				listView.Refresh()

			}),
			history,
		),
		nil, nil, nil, listView, w, func() {
			listView.Refresh()

		})

	w.SetContent(container.NewBorder(
		burgerMenu, nil, nil, nil,
		nil,
	))

	return &MediaTools{
		app:    &app,
		window: &w,
	}
}

func (mt *MediaTools) Run() {
	(*mt.window).ShowAndRun()
}
