package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"
)

type ListItem struct {
	widget.BaseWidget

	checkBox     *widget.Check
	filename     *widget.Label
	removeButton *widget.Button
	infoButton   *widget.Button

	ffprobeData *medias.FfprobeResult
}

func NewListItem() *ListItem {
	li := &ListItem{}
	li.initUI()
	li.ExtendBaseWidget(li)
	return li
}

func (li *ListItem) initUI() {
	li.checkBox = widget.NewCheck("", nil)
	li.filename = widget.NewLabel("filename")
	li.removeButton = widget.NewButtonWithIcon("", theme.DeleteIcon(), nil)
	li.infoButton = widget.NewButtonWithIcon("", theme.InfoIcon(), nil)
}

func (mi *ListItem) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewHBox(
		mi.checkBox,
		mi.removeButton,
		mi.infoButton,
		mi.filename,
	))
}

func (li *ListItem) SetFfprobeData(data *medias.FfprobeResult) {
	li.ffprobeData = data
}
