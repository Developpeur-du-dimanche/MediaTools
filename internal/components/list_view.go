package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/list"
)

type ListView struct {
	widget.BaseWidget

	list *widget.List

	items    *list.List[string]
	OnUpdate chan bool
}

func NewListView(items *list.List[string]) *ListView {
	lv := &ListView{OnUpdate: make(chan bool), list: widget.NewList(
		func() int {
			return items.GetLength()
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Item")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(items.GetItem(int(i)))
		},
	), items: items}
	return lv
}

func (lv *ListView) CreateRenderer() fyne.WidgetRenderer {

	return lv.list.CreateRenderer()
}
