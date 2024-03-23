package application

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type LayoutInterface interface {
	SetContent(content fyne.CanvasObject) *fyne.Container
	GetFileList() *widget.List
}

type Layout struct {
	layout   *fyne.Container
	fileList *widget.List
}

func NewLayout(openFolder *widget.Button, openFile *widget.Button, fileList *widget.List) LayoutInterface {

	layout := container.NewAdaptiveGrid(3, container.NewVBox(
		container.NewHBox(
			openFolder,
			openFile,
		),
		fileList,
	), container.NewVScroll(widget.NewLabel("")))
	return &Layout{layout: layout, fileList: fileList}
}

func (l *Layout) SetContent(content fyne.CanvasObject) *fyne.Container {
	l.layout.Objects[1] = content
	return l.layout
}

func (l *Layout) GetFileList() *widget.List {
	return l.fileList
}
