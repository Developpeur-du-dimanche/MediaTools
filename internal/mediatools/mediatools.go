package mediatools

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/components"
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

	//listView := components.NewListView(items)

	history := components.NewLastScanSelector(func(path string) {
		fmt.Printf("Folder selected: %s\n", path)
	})

	openFile := components.NewOpenFile(&w, func(path string) {
		fmt.Printf("File opened: %s\n", path)
		history.AddFolder(path)
	})

	openFolder := components.NewOpenFolder(&w, func(path string) {
		history.AddFolder(path)

	}, func(path string) {
		fmt.Printf("File detected: %s\n", path)
	})

	w.SetContent(container.NewBorder(
		container.NewHBox(
			history,
			openFile,
			openFolder,
		), nil, nil, nil,
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
