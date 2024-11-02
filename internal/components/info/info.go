package info

import "fyne.io/fyne/v2/widget"

type Info interface {
	GetNodes(id string) []widget.TreeNodeID
	From(label string, id string)
}
