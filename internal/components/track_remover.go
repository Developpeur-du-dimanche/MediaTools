package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type TrackRemoverComponent struct {
}

func NewTrackRemoverComponent() *TrackRemoverComponent {
	return &TrackRemoverComponent{}
}

func (f *TrackRemoverComponent) Content() fyne.CanvasObject {
	return widget.NewLabel("TrackRemoverComponent")
}
