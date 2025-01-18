package components

import (
	"io/fs"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type OpenFolder struct {
	widget.BaseWidget
	button *widget.Button
	window *fyne.Window

	infoDialog  *dialog.CustomDialog
	infoMessage *widget.Label

	onFolderOpen     func(path string)
	onFileDetected   func(path string)
	OnScanTerminated func()
}

func NewOpenFolder(parent *fyne.Window, onFolderOpened func(path string), onFileDetected func(path string)) *OpenFolder {
	of := &OpenFolder{
		window:         parent,
		onFolderOpen:   onFolderOpened,
		onFileDetected: onFileDetected,
		infoMessage:    widget.NewLabel(""),
	}

	of.infoDialog = dialog.NewCustomWithoutButtons("Media Info", of.infoMessage, *of.window)
	of.infoDialog.Hide()
	of.button = widget.NewButtonWithIcon("Open Folder", theme.FolderIcon(), of.openFolderDialog)
	of.ExtendBaseWidget(of)
	return of
}

func (of *OpenFolder) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(of.button)
}

func (of *OpenFolder) openFolderDialog() {
	dialog := dialog.NewFolderOpen(func(lu fyne.ListableURI, err error) {
		if err != nil {
			dialog.ShowError(err, *of.window)
			return
		}

		if lu == nil {
			return
		}

		list, err := lu.List()
		if err != nil {
			dialog.ShowError(err, *of.window)
			return
		}

		if of.onFolderOpen != nil {
			of.onFolderOpen(lu.Path())
		}

		of.infoDialog.Show()

		defer of.infoDialog.Hide()

		var path string

		if of.onFolderOpen != nil {
			for _, item := range list {
				path = item.Path()
				if i, err := os.Stat(path); err == nil && i.IsDir() {
					filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
						if err != nil {
							return err
						}

						if filepath.Ext(path) == "" {
							return nil
						}

						if info.IsDir() {
							return nil
						}

						if of.onFileDetected != nil {
							of.infoMessage.SetText("Scanning " + path)
							of.onFileDetected(path)
						}

						return nil
					})
				} else {
					if of.onFileDetected != nil {
						of.infoMessage.SetText("Scanning " + path)
						of.onFileDetected(path)
					}
				}
			}
		}
		if of.OnScanTerminated != nil {
			of.OnScanTerminated()
		}
	}, *of.window)
	size := (*of.window).Canvas().Size()
	dialog.Resize(fyne.NewSize(size.Width-150, size.Height-150))
	dialog.Show()
}
