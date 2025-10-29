package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type MergeItem struct {
	widget.BaseWidget

	upButton     *widget.Button
	downButton   *widget.Button
	removeButton *widget.Button
	fileLabel    *widget.Label
}

func NewMergeItem() *MergeItem {
	mi := &MergeItem{}

	mi.initUI()
	mi.ExtendBaseWidget(mi)
	return mi
}

func (mi *MergeItem) initUI() {
	mi.upButton = widget.NewButtonWithIcon("", theme.MoveUpIcon(), nil)
	mi.downButton = widget.NewButtonWithIcon("", theme.MoveDownIcon(), nil)
	mi.removeButton = widget.NewButtonWithIcon("", theme.DeleteIcon(), nil)
	mi.fileLabel = widget.NewLabel("filename")
}

func (mi *MergeItem) SetLabel(text string) {
	mi.fileLabel.SetText(text)
}

func (mi *MergeItem) DisableUpButton() {
	mi.upButton.Disable()
}

func (mi *MergeItem) EnableUpButton() {
	mi.upButton.Enable()
}

func (mi *MergeItem) DisableDownButton() {
	mi.downButton.Disable()
}

func (mi *MergeItem) EnableDownButton() {
	mi.downButton.Enable()
}

func (mi *MergeItem) SetUpButtonOnTapped(onTapped func()) {
	mi.upButton.OnTapped = onTapped
}

func (mi *MergeItem) SetDownButtonOnTapped(onTapped func()) {
	mi.downButton.OnTapped = onTapped
}

func (mi *MergeItem) SetRemoveButtonOnTapped(onTapped func()) {
	mi.removeButton.OnTapped = onTapped
}

func (mi *MergeItem) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewBorder(
		nil, nil,
		container.NewHBox(
			mi.upButton,
			mi.downButton,
		),
		mi.removeButton,
		mi.fileLabel,
	))
}
