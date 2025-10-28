package components

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/services"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"
)

type ListView struct {
	widget.BaseWidget

	list          *widget.List
	mutex         *sync.Mutex
	window        fyne.Window
	ffmpegService *services.FFmpegService

	items      []*medias.FfprobeResult
	OnUpdate   chan bool
	OnRefresh  func()
	isSelected map[int]bool
}

func NewListView(onRefresh func(), window fyne.Window, ffmpegService *services.FFmpegService) *ListView {
	mutex := &sync.Mutex{}
	lv := &ListView{
		OnUpdate:      make(chan bool),
		items:         []*medias.FfprobeResult{},
		mutex:         mutex,
		OnRefresh:     onRefresh,
		list:          widget.NewList(nil, nil, nil),
		window:        window,
		ffmpegService: ffmpegService,
		isSelected:    make(map[int]bool),
	}

	return lv
}

func (lv *ListView) CreateRenderer() fyne.WidgetRenderer {
	lv.list.Length = lv.length
	lv.list.CreateItem = lv.CreateItem
	lv.list.UpdateItem = lv.updateItem
	lv.list.OnSelected = lv.onSelected
	return widget.NewSimpleRenderer(lv.list)
}

func (lv *ListView) CreateItem() fyne.CanvasObject {
	return container.NewHBox(
		widget.NewCheck("", nil),
		widget.NewButtonWithIcon("", theme.DeleteIcon(), nil),
		widget.NewButtonWithIcon("", theme.InfoIcon(), nil),
		widget.NewLabel(""),
	)
}

func (lv *ListView) onSelected(id int) {
	lv.list.Unselect(id)
}

func (lv *ListView) length() int {
	return len(lv.items)
}

func (lv *ListView) updateItem(i widget.ListItemID, o fyne.CanvasObject) {
	root := o.(*fyne.Container)
	item := lv.items[i]

	checkbox := root.Objects[0].(*widget.Check)
	checkbox.Checked = lv.isSelected[int(i)]
	checkbox.OnChanged = func(checked bool) {
		lv.isSelected[int(i)] = checked
	}

	root.Objects[1].(*widget.Button).OnTapped = func() {
		lv.RemoveItemAt(int(i))
		lv.Refresh()
	}

	root.Objects[2].(*widget.Button).OnTapped = func() {
		fic := NewFileInfoComponent(*item, lv.window)
		dialog := dialog.NewCustom("Media Info", "Close", fic, lv.window)
		dialog.Show()
	}

	root.Objects[3].(*widget.Label).SetText(item.Format.Filename)
}

func (lv *ListView) RemoveItemAt(index int) {
	lv.mutex.Lock()
	defer lv.mutex.Unlock()
	if len(lv.items) == 0 {
		return
	}
	lv.items = append(lv.items[:index], lv.items[index+1:]...)
	lv.list.Refresh()
}

func (lv *ListView) AddItem(item *medias.FfprobeResult) {
	lv.mutex.Lock()
	lv.items = append(lv.items, item)
	lv.list.Refresh()
	lv.mutex.Unlock()
}

func (lv *ListView) Clear() {
	lv.mutex.Lock()
	lv.items = make([]*medias.FfprobeResult, 0)
	lv.list.Refresh()
	lv.mutex.Unlock()
}

func (lv *ListView) GetItems() []*medias.FfprobeResult {
	return lv.items
}

func (lv *ListView) GetSelectedItems() []*medias.FfprobeResult {
	lv.mutex.Lock()
	defer lv.mutex.Unlock()

	selected := make([]*medias.FfprobeResult, 0)
	for i, item := range lv.items {
		if lv.isSelected[i] {
			selected = append(selected, item)
		}
	}
	return selected
}

func (lv *ListView) SelectAll() {
	lv.mutex.Lock()
	defer lv.mutex.Unlock()

	for i := range lv.items {
		lv.isSelected[i] = true
	}
	lv.list.Refresh()
}

func (lv *ListView) UnselectAll() {
	lv.mutex.Lock()
	defer lv.mutex.Unlock()

	lv.isSelected = make(map[int]bool)
	lv.list.Refresh()
}
