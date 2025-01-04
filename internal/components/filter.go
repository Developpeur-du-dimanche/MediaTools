package components

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"
	mediatools_embed "github.com/Developpeur-du-dimanche/MediaTools"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/helper"
	jsonfilter "github.com/Developpeur-du-dimanche/MediaTools/pkg/filter"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/list"
)

type FilterComponent struct {
	choices            []*ConditionalWidget
	container          *fyne.Container
	fileList           *list.List[*helper.FileMetadata]
	window             *fyne.Window
	filterButton       *widget.Button
	addFilterButton    *widget.Button
	removeFilterButton *widget.Button
	filters            *jsonfilter.Filters
}

func NewFilterComponent(window *fyne.Window, fileList *list.List[*helper.FileMetadata]) Component {
	jf := jsonfilter.NewParser(mediatools_embed.Filters)
	p, err := jf.Parse()

	if err != nil {
		dialog.ShowError(err, *window)
	}

	return &FilterComponent{
		choices:            []*ConditionalWidget{},
		container:          container.NewVBox(),
		fileList:           fileList,
		window:             window,
		filterButton:       widget.NewButton(lang.L("filter"), nil),
		addFilterButton:    widget.NewButton(lang.L("add_filter"), nil),
		removeFilterButton: widget.NewButton(lang.L("remove_filter"), nil),
		filters:            p,
	}
}

func (f *FilterComponent) Content() fyne.CanvasObject {

	f.addFilterButton.OnTapped = f.addFilter
	f.removeFilterButton.OnTapped = f.removeFilter

	f.filterButton.OnTapped = func() {
		f.filterButton.Disable()
		f.Filter()
		f.filterButton.Enable()
	}

	return container.NewBorder(
		container.NewHBox(
			f.addFilterButton,
			f.removeFilterButton,
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

	output := list.NewList[*helper.FileMetadata]()
	treatmentOf := widget.NewLabel("file is currently being treated, please wait...")
	cd := dialog.NewCustomWithoutButtons("Please wait", treatmentOf, *f.window)
	cd.Show()
	for _, file := range f.fileList.GetItems() {
		treatmentOf.SetText("file is currently being treated: " + file.FileName + " please wait...")

		isValid := 0
		for _, c := range f.choices {

			if c.choice.Check(file, jsonfilter.Condition(c.condition), c.value) {
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

func (f *FilterComponent) addFilter() {
	c := NewConditionalWidget(f.filters)
	f.choices = append(f.choices, c)
	f.container.Add(c.container)
}

func (f *FilterComponent) removeFilter() {
	if len(f.choices) > 0 {
		f.choices = f.choices[:len(f.choices)-1]
		f.container.Remove(f.container.Objects[len(f.container.Objects)-1])
	}
}
