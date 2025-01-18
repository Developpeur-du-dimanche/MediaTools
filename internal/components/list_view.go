package components

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"
)

type ListView struct {
	widget.BaseWidget

	list  *widget.List
	mutex *sync.Mutex

	items     []*medias.FfprobeResult
	OnUpdate  chan bool
	OnRefresh func()
}

func NewListView(onRefresh func()) *ListView {
	mutex := &sync.Mutex{}
	lv := &ListView{
		OnUpdate:  make(chan bool),
		items:     []*medias.FfprobeResult{},
		mutex:     mutex,
		OnRefresh: onRefresh,
	}
	lv.list = widget.NewList(
		func() int {
			return len(lv.items)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewButtonWithIcon("", theme.DeleteIcon(), nil),
				widget.NewButtonWithIcon("", theme.InfoIcon(), nil),
				widget.NewLabel(""),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			root := o.(*fyne.Container)
			item := lv.items[i]
			root.Objects[0].(*widget.Button).OnTapped = func() {
				lv.RemoveItemAt(int(i))
				lv.Refresh()
			}

			root.Objects[1].(*widget.Button).OnTapped = func() {
				fic := NewFileInfoComponent(*item)
				dialog := dialog.NewCustom("Media Info", "Close", fic, fyne.CurrentApp().Driver().AllWindows()[0])
				dialog.Show()
			}

			root.Objects[2].(*widget.Label).SetText(item.Format.Filename)
		},
	)

	lv.list.OnSelected = func(id widget.ListItemID) {
		lv.list.Unselect(id)
	}

	return lv
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

func (lv *ListView) AddItems(items []*medias.FfprobeResult) {
	lv.mutex.Lock()
	lv.items = append(lv.items, items...)
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

func (lv *ListView) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(lv.list)
}
