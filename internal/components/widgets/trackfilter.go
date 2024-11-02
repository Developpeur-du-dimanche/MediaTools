package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/filter"
)

type TrackFilter struct {
	widget.BaseWidget
	condition filter.ConditionString
	value     string

	keyWidget       *widget.Label
	conditionWidget *widget.Select
	valueWidget     *widget.Entry

	filter *filter.Filter
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
		tf.condition = filter.FromString(s)
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

func (tf *TrackFilter) GetCondition() filter.ConditionString {
	return tf.condition
}

func (tf *TrackFilter) GetValue() string {
	return tf.value
}

func (tf *TrackFilter) Equals(other *TrackFilter) bool {
	return tf.keyWidget.Text == other.keyWidget.Text
}

func (tf *TrackFilter) SetFilter(filter filter.Filter) {
	tf.filter = &filter
}

func (tf *TrackFilter) GetFilter() filter.Filter {
	return *tf.filter
}
