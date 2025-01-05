package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/components/customs"
	jsonfilter "github.com/Developpeur-du-dimanche/MediaTools/pkg/filter"
)

type ConditionalWidget struct {
	widget.BaseWidget
	key       *widget.Select
	choice    jsonfilter.Filter
	container *fyne.Container
	condition string
	value     string
}

func NewConditionalWidget(filters *jsonfilter.Filters) *ConditionalWidget {
	choices := make([]string, len(filters.Filters))
	for i, filter := range filters.Filters {
		choices[i] = filter.Name
	}

	key := widget.NewSelect(choices, nil)
	choice := widget.NewSelect([]string{"Select a condition"}, nil)
	choice.Disable()
	value := widget.NewSelect([]string{"Select a value"}, nil)
	value.Disable()

	c := &ConditionalWidget{
		key:       key,
		container: container.NewGridWithColumns(3, key, choice, value),
	}

	c.key.OnChanged = func(s string) {
		c.updateChoiceAndValueWidgets(s, filters)
	}

	c.ExtendBaseWidget(c)
	return c
}

func (c *ConditionalWidget) updateChoiceAndValueWidgets(selectedKey string, filters *jsonfilter.Filters) {
	for _, filter := range filters.Filters {
		if filter.Name == selectedKey {
			c.choice = filter
			c.updateChoiceWidget(filter)
			c.updateValueWidget(filter)
			break
		}
	}
}

func (c *ConditionalWidget) updateChoiceWidget(filter jsonfilter.Filter) {
	choiceWidget := widget.NewSelect(filter.GetStringCondition(), func(s string) {
		c.condition = s
	})
	choiceWidget.SetSelectedIndex(0)
	c.container.Objects[1] = choiceWidget
}

func (c *ConditionalWidget) updateValueWidget(filter jsonfilter.Filter) {
	if filter.HasDefaultValues() {
		c.value = filter.GetDefaultValues()[0]
		valueWidget := widget.NewSelect(filter.GetDefaultValues(), func(s string) {
			c.value = s
		})
		valueWidget.SetSelectedIndex(0)
		c.container.Objects[2] = valueWidget
		valueWidget.Enable()
	} else {
		c.createCustomValueWidget(filter)
	}
}

func (c *ConditionalWidget) createCustomValueWidget(filter jsonfilter.Filter) {
	var valueWidget customs.CustomEntry
	switch filter.GetType() {
	case jsonfilter.Int:
		valueWidget = customs.NewNumericalEntry()
		valueWidget.(*customs.NumericalEntry).OnChanged = func(s string) {
			c.value = s
		}
	case jsonfilter.Bool:
		boolWidget := widget.NewSelect([]string{"true", "false"}, func(s string) {
			c.value = s
		})
		boolWidget.SetSelectedIndex(0)
		boolWidget.Enable()
		c.container.Objects[2] = boolWidget
		return
	default:
		valueWidget = widget.NewEntry()
		valueWidget.(*widget.Entry).OnChanged = func(s string) {
			c.value = s
		}
	}
	c.container.Objects[2] = valueWidget
	valueWidget.Enable()
}

func (c *ConditionalWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(c.container)
}

func (c *ConditionalWidget) MinSize() fyne.Size {
	return fyne.NewSize(30, 40)
}
