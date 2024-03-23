package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type FilterComponent struct {
	choices   *[]*ConditionalWidget
	container *fyne.Container
	fileList  *[]string
}

func NewFilterComponent(fileList *[]string) *FilterComponent {
	return &FilterComponent{
		choices:   &[]*ConditionalWidget{},
		container: container.NewVBox(),
		fileList:  fileList,
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

	filterButton := widget.NewButton("Filter", func() {
		//TODO
	})

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
