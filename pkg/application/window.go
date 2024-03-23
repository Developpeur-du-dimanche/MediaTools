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
	Window        *fyne.Window
	mainContainer LayoutInterface
	size          *fyne.Size
}

func newWindow(window *fyne.Window, mainContainer LayoutInterface, size *fyne.Size) WindowInterface {
	return &Window{Window: window, mainContainer: mainContainer, size: size}
}

func (w *Window) SetContent(content fyne.CanvasObject) {
	(*w.Window).SetContent(w.mainContainer.SetContent(content))
}

func (w *Window) SetMainMenu(menu *fyne.MainMenu) {
	(*w.Window).SetMainMenu(menu)
}

func (w *Window) NewFileOpen(callback func(fyne.URIReadCloser, error), location fyne.ListableURI) *dialog.FileDialog {
	dialog := dialog.NewFileOpen(callback, *w.Window)
	size := (*w.Window).Canvas().Size()
	dialog.Resize(fyne.NewSize(size.Width-150, size.Height-150))
	return dialog
}

func (w *Window) NewFolderOpen(callback func(fyne.ListableURI, error), location fyne.ListableURI) *dialog.FileDialog {
	dialog := dialog.NewFolderOpen(callback, *w.Window)
	size := (*w.Window).Canvas().Size()
	dialog.Resize(fyne.NewSize(size.Width-150, size.Height-150))
	return dialog
}

func (w *Window) ShowAndRun() {
	(*w.Window).ShowAndRun()
}
