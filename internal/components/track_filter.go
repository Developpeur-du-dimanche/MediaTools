package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type StreamType string

const (
	Video      StreamType = "Video"
	Audio      StreamType = "Audio"
	Subtitle   StreamType = "Subtitle"
	Attachment StreamType = "Attachment"
)

type TrackFilter struct {
	widget.BaseWidget
	condition string
	value     string

	keyWidget       *widget.Label
	conditionWidget *widget.Select
	valueWidget     *widget.Entry

	StreamType StreamType
}

func NewTrackFilter(text string) *TrackFilter {
	tf := &TrackFilter{
		value:     "",
		keyWidget: widget.NewLabel(text),
		conditionWidget: widget.NewSelect([]string{
			"ignore",
			"equals",
			"not equals",
			"contains",
		}, nil),
		valueWidget: widget.NewEntry(),
	}

	tf.conditionWidget.OnChanged = func(s string) {
		tf.condition = s
	}

	tf.valueWidget.OnChanged = func(s string) {
		tf.value = s
	}

	return tf
}

func (tf *TrackFilter) CreateRenderer() fyne.WidgetRenderer {

	return widget.NewSimpleRenderer(container.NewAdaptiveGrid(3,
		tf.keyWidget,
		tf.conditionWidget,
		tf.valueWidget,
	))
}

func (tf *TrackFilter) MinSize() fyne.Size {
	return tf.keyWidget.MinSize()
}

func (tf *TrackFilter) SetText(text string) {
	tf.keyWidget.SetText(text)
}

func (tf *TrackFilter) GetValue() string {
	return tf.value
}

func (tf *TrackFilter) Equals(other *TrackFilter) bool {
	return tf.keyWidget.Text == other.keyWidget.Text
}
