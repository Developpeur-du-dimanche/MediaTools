package components

import (
	"errors"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/filter"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/fileinfo"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/list"
)

var conditions = []filter.ConditionContract{
	filter.NewContainerFilter(),
	filter.NewAudioLanguageFilter(),
	filter.NewBitrateFilter(),
	filter.NewSubtitleForcedFilter(),
	filter.NewSubtitleLanguageFilter(),
	filter.NewSubtitleTitleFilter(),
	filter.NewSubtitleCodecFilter(),
	filter.NewVideoTitleFilter(),
}

type FilterComponent struct {
	choices      []*ConditionalWidget
	container    *fyne.Container
	fileList     *list.List[fileinfo.FileInfo]
	window       *fyne.Window
	filterButton *widget.Button
}

func NewFilterComponent(window *fyne.Window, fileList *list.List[fileinfo.FileInfo]) *FilterComponent {
	return &FilterComponent{
		choices:      []*ConditionalWidget{},
		container:    container.NewVBox(),
		fileList:     fileList,
		window:       window,
		filterButton: widget.NewButton("Filter", nil),
	}
}

func (f *FilterComponent) Content() fyne.CanvasObject {

	addFilterButton := widget.NewButton("Add filter", func() {
		nc := NewConditionalWidget(conditions)
		f.choices = append(f.choices, nc)
		f.container.Add(nc)
	})

	removeFilterButton := widget.NewButton("Remove filter", func() {
		if len(f.choices) > 0 {
			f.choices = f.choices[:len(f.choices)-1]
			f.container.Remove(f.container.Objects[len(f.container.Objects)-1])
		}
	})

	f.filterButton.OnTapped = f.Filter

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
	f.filterButton.Disable()

	if f.fileList.GetLength() == 0 || len(f.choices) == 0 {
		dialog.ShowError(errors.New("no file selected or no filter added"), *f.window)
		f.filterButton.Enable()
		return
	}

	output := list.NewList[fileinfo.FileInfo]()
	treatmentOf := widget.NewLabel("file is currently being treated, please wait...")
	cd := dialog.NewCustomWithoutButtons("Please wait", treatmentOf, *f.window)
	cd.Show()
	for _, file := range f.fileList.GetItems() {
		treatmentOf.SetText("file is currently being treated: " + file.GetPath() + " please wait...")

		data := file.GetInfo()
		isValid := false
		for _, c := range f.choices {

			if c.choice.CheckGlobal(data) {
				isValid = true
			}
		}
		if isValid {
			output.AddItem(file)
		}
	}

	cd.Hide()

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
			item.(*widget.Label).SetText(output.GetItem(i).GetPath())
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

	f.filterButton.Enable()

}
