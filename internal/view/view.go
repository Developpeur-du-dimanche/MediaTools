package view

import "fyne.io/fyne/v2"

type View interface {
	Content() fyne.CanvasObject
	GetWindow() *fyne.Window
	ShowAndRun()
}
