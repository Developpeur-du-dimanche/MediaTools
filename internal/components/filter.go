package components

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"
	jsonembed "github.com/Developpeur-du-dimanche/MediaTools"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/fileinfo"
	jsonfilter "github.com/Developpeur-du-dimanche/MediaTools/pkg/filter"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/list"
)

type FilterComponent struct {
	choices      []*ConditionalWidget
	container    *fyne.Container
	fileList     *list.List[fileinfo.FileInfo]
	window       *fyne.Window
	filterButton *widget.Button
	filters      *jsonfilter.Filters
}

func NewFilterComponent(window *fyne.Window, fileList *list.List[fileinfo.FileInfo]) *FilterComponent {
	jf := jsonfilter.NewParser(jsonembed.Filters)
	p, err := jf.Parse()

	if err != nil {
		dialog.ShowError(err, *window)
	}

	return &FilterComponent{
		choices:      []*ConditionalWidget{},
		container:    container.NewVBox(),
		fileList:     fileList,
		window:       window,
		filterButton: widget.NewButton(lang.L("filter"), nil),
		filters:      p,
	}
}

func (f *FilterComponent) Content() fyne.CanvasObject {

	addFilterButton := widget.NewButton("Add filter", func() {
		nc := NewConditionalWidget(f.filters)
		f.choices = append(f.choices, nc)
		f.container.Add(nc)
	})

	removeFilterButton := widget.NewButton("Remove filter", func() {
		if len(f.choices) > 0 {
			f.choices = f.choices[:len(f.choices)-1]
			f.container.Remove(f.container.Objects[len(f.container.Objects)-1])
		}
	})

	f.filterButton.OnTapped = func() {
		f.filterButton.Disable()
		f.Filter()
		f.filterButton.Enable()
	}

	return container.NewBorder(
		container.NewHBox(
			addFilterButton,
			removeFilterButton,
		),
		f.filterButton,
		nil,
		nil,
		container.NewVScroll(
			f.container,
		),
	)
}

func (f *FilterComponent) Filter() {

	if f.fileList.GetLength() == 0 || len(f.choices) == 0 {
		dialog.ShowError(errors.New("no file selected or no filter added"), *f.window)
		return
	}

	output := list.NewList[fileinfo.FileInfo]()
	treatmentOf := widget.NewLabel("file is currently being treated, please wait...")
	cd := dialog.NewCustomWithoutButtons("Please wait", treatmentOf, *f.window)
	cd.Show()
	for _, file := range f.fileList.GetItems() {
		treatmentOf.SetText("file is currently being treated: " + file.GetPath() + " please wait...")

		data := file.GetInfo()
		isValid := 0
		for _, c := range f.choices {

			if c.choice.Check(data, c.condition) {
				isValid++
			}
		}
		if isValid == len(f.choices) {
			output.AddItem(file)
		}
	}

	cd.Hide()

	// create new window to display filtered files
	NewResultList(output).Show()

}
