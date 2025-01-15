package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type OpenFile struct {
	widget.BaseWidget
	button *widget.Button
	window *fyne.Window

	onFileOpen func(path string)
}

func NewOpenFile(parent *fyne.Window, onFileOpened func(path string)) *OpenFile {
	of := &OpenFile{
		window:     parent,
		onFileOpen: onFileOpened,
	}
	of.button = widget.NewButtonWithIcon("Open File", theme.FileIcon(), of.openFileDialog)
	of.ExtendBaseWidget(of)
	return of
}

func (of *OpenFile) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(of.button)
}

func (of *OpenFile) openFileDialog() {
	dialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, *of.window)
			return
		}

		if reader == nil {
			return
		}

		if of.onFileOpen != nil {
			of.onFileOpen(reader.URI().String())
		}
	}, *of.window)
	size := (*of.window).Canvas().Size()
	dialog.Resize(fyne.NewSize(size.Width-150, size.Height-150))
	dialog.Show()
}
