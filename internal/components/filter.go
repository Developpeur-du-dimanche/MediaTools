package components

import (
	"context"
	"errors"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/list"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type FilterComponent struct {
	choices   *[]*ConditionalWidget
	container *fyne.Container
	fileList  *list.List
	window    *fyne.Window
}

func NewFilterComponent(window *fyne.Window, fileList *list.List) *FilterComponent {
	return &FilterComponent{
		choices:   &[]*ConditionalWidget{},
		container: container.NewVBox(),
		fileList:  fileList,
		window:    window,
	}
}

func (f *FilterComponent) Content() fyne.CanvasObject {

	addFilterButton := widget.NewButton("Add filter", func() {
		nc := NewConditionalWidget()
		*f.choices = append(*f.choices, nc)
		f.container.Add(nc)
	})

	removeFilterButton := widget.NewButton("Remove filter", func() {
		if len(*f.choices) > 0 {
			*f.choices = (*f.choices)[:len(*f.choices)-1]
			f.container.Remove(f.container.Objects[len(f.container.Objects)-1])
		}
	})

	filterButton := widget.NewButton("Filter", f.Filter)

	return container.NewBorder(
		container.NewHBox(
			addFilterButton,
			removeFilterButton,
		),
		filterButton,
		nil,
		nil,
		container.NewVScroll(
			f.container,
		),
	)
}

func (f *FilterComponent) Filter() {

	if f.fileList.GetLength() == 0 || len(*f.choices) == 0 {
		fmt.Printf("no file selected or no filter added, len %d\n", f.fileList.GetLength())
		dialog.ShowError(errors.New("no file selected or no filter added"), *f.window)
		return
	}

	output := list.NewList()
	for _, file := range f.fileList.GetItems() {
		ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelFn()

		data, err := ffprobe.ProbeURL(ctx, file)
		for _, c := range *f.choices {

			if err != nil {
				fmt.Printf("error while probing file %s: %s\n", file, err)
				dialog.ShowError(err, *f.window)
				return
			}

			if c.choice.Check(data) {
				output.AddItem(file)
				break
			}
		}
	}

	// create new window to display filtered files
	w := fyne.CurrentApp().NewWindow("Filtered files")
	screen := fyne.CurrentApp().Driver().AllWindows()[0].Canvas().Size()
	w.Resize(fyne.NewSize(float32(screen.Width/2), float32(screen.Height/2)))

	list := widget.NewList(
		func() int {
			return output.GetLength()
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(output.GetItem(i))
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

	w.Show()

}
