package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/list"
)

type ListView struct {
	widget.BaseWidget

	list *widget.List

	items     *list.List[string]
	OnUpdate  chan bool
	OnRefresh func()
}

func NewListView(items *list.List[string], onRefresh func()) *ListView {
	lv := &ListView{OnUpdate: make(chan bool),
		items:     items,
		OnRefresh: onRefresh,
	}
	lv.list = widget.NewList(
		func() int {
			return items.GetLength()
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewButtonWithIcon("", theme.DeleteIcon(), nil),
				widget.NewLabel(""),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			root := o.(*fyne.Container)
			root.Objects[0].(*widget.Button).OnTapped = func() {
				items.RemoveItemAt(int(i))
				lv.Refresh()
			}
			root.Objects[1].(*widget.Label).SetText(items.GetItem(int(i)))
		},
	)
	return lv
}

func (lv *ListView) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(lv.list)
}

func (lv *ListView) AddItem(item string) {
	lv.items.AddItem(item)
	lv.Refresh()
}

func (lv *ListView) Refresh() {
	lv.list.Refresh()
	if lv.OnRefresh != nil {
		lv.OnRefresh()
	}

}
