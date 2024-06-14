package application

import "fyne.io/fyne/v2"

type ListInterface interface {
	GetCanvas() fyne.CanvasObject
	container() fyne.Container
}
