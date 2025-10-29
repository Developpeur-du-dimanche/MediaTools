package components

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
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

	items       []*medias.FfprobeResult
	OnUpdate    chan bool
	OnRefresh   func()
	isSelected  map[int]bool
	maxSize     int
	currentSize int
}

func NewListView(onRefresh func(), window fyne.Window, ffmpegService *services.FFmpegService) *ListView {
	mutex := &sync.Mutex{}
	lv := &ListView{
		OnUpdate:      make(chan bool),
		items:         make([]*medias.FfprobeResult, 100),
		mutex:         mutex,
		OnRefresh:     onRefresh,
		list:          widget.NewList(nil, nil, nil),
		window:        window,
		ffmpegService: ffmpegService,
		isSelected:    make(map[int]bool),
		maxSize:       100,
		currentSize:   0,
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
	return NewListItem()
}

func (lv *ListView) onSelected(id int) {
	lv.list.Unselect(id)
}

func (lv *ListView) length() int {
	return lv.currentSize
}

func (lv *ListView) updateItem(i widget.ListItemID, o fyne.CanvasObject) {
	root := o.(*ListItem)
	item := lv.items[i]

	if item == nil {
		return
	}

	root.checkBox.Checked = lv.isSelected[int(i)]
	root.checkBox.OnChanged = func(checked bool) {
		lv.isSelected[int(i)] = checked
	}

	root.removeButton.OnTapped = func() {
		lv.RemoveItemAt(int(i))
		lv.Refresh()
	}

	root.infoButton.OnTapped = func() {
		fic := NewFileInfoComponent(item, lv.window)
		dialog := dialog.NewCustom("Media Info", "Close", fic, lv.window)
		dialog.Show()
	}

	root.filename.SetText(item.Format.Filename)
}

func (lv *ListView) RemoveItemAt(index int) {
	lv.mutex.Lock()
	defer lv.mutex.Unlock()
	if lv.currentSize == 0 || index >= lv.currentSize {
		return
	}

	// Shift all items after the removed index
	for i := index; i < lv.currentSize-1; i++ {
		lv.items[i] = lv.items[i+1]
	}
	lv.items[lv.currentSize-1] = nil
	lv.currentSize--

	// Update selection map
	newSelected := make(map[int]bool)
	for i, selected := range lv.isSelected {
		if i < index {
			newSelected[i] = selected
		} else if i > index {
			newSelected[i-1] = selected
		}
	}
	lv.isSelected = newSelected

	lv.list.Refresh()
}

func (lv *ListView) AddItem(item *medias.FfprobeResult) {
	if item == nil {
		return
	}
	lv.mutex.Lock()
	defer lv.mutex.Unlock()
	if lv.currentSize >= lv.maxSize-10 {
		newList := make([]*medias.FfprobeResult, lv.maxSize*2)
		copy(newList, lv.items)
		lv.items = newList
		lv.maxSize *= 2
	}
	lv.items[lv.currentSize] = item
	lv.currentSize++
	lv.list.Refresh()
}

func (lv *ListView) Clear() {
	lv.mutex.Lock()
	defer lv.mutex.Unlock()
	lv.items = make([]*medias.FfprobeResult, 100)
	lv.currentSize = 0
	lv.maxSize = 100
	lv.isSelected = make(map[int]bool)
	lv.list.Refresh()
}

func (lv *ListView) GetItems() []medias.FfprobeResult {
	lv.mutex.Lock()
	defer lv.mutex.Unlock()
	items := make([]medias.FfprobeResult, lv.currentSize)
	for i, item := range lv.items {
		if item != nil {
			items[i] = *item
		}
	}
	return items
}

func (lv *ListView) GetSelectedItems() []*medias.FfprobeResult {
	lv.mutex.Lock()
	defer lv.mutex.Unlock()

	selected := make([]*medias.FfprobeResult, 0)
	for i := 0; i < lv.currentSize; i++ {
		if lv.isSelected[i] && lv.items[i] != nil {
			selected = append(selected, lv.items[i])
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
