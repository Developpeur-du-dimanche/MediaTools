package application

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

type WindowInterface interface {
	SetContent(content fyne.CanvasObject)
	SetMainMenu(menu *fyne.MainMenu)
	NewFileOpen(callback func(fyne.URIReadCloser, error), location fyne.ListableURI) *dialog.FileDialog
	NewFolderOpen(callback func(fyne.ListableURI, error), location fyne.ListableURI) *dialog.FileDialog
	ShowAndRun()
}

type Window struct {
	Window fyne.Window
}

func newWindow(window fyne.Window) *Window {
	return &Window{Window: window}
}
