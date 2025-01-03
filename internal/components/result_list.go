package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/helper"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/list"
)

type ResultList struct {
	result *list.List[*helper.FileMetadata]
	window fyne.Window
}

func NewResultList(result *list.List[*helper.FileMetadata]) *ResultList {
	w := fyne.CurrentApp().NewWindow("Filtered files")
	screen := fyne.CurrentApp().Driver().AllWindows()[0].Canvas().Size()
	w.Resize(fyne.NewSize(float32(screen.Width/2), float32(screen.Height/2)))

	list := widget.NewList(
		func() int {
			return result.GetLength()
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(result.GetItem(i).FileName)
		},
	)

	list.Resize(fyne.NewSize(float32(screen.Width/4), float32(screen.Height/2)))

	w.SetContent(container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		list,
	))

	return &ResultList{
		result: result,
		window: w,
	}
}

func (r *ResultList) Show() {
	r.window.Show()
}
