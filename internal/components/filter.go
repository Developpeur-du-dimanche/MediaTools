package components

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type FilterComponent struct {
	choices   *[]*ConditionalWidget
	container *fyne.Container
	fileList  *[]string
	window    *fyne.Window
}

func NewFilterComponent(fileList *[]string, window *fyne.Window) *FilterComponent {
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
	if len(*f.fileList) == 0 || len(*f.choices) == 0 {
		dialog.ShowError(errors.New("no file selected or no filter added"), *f.window)
		return
	}
}
