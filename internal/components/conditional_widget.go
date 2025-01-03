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
		for _, filter := range filters.Filters {
			if filter.Name == s {
				c.choice = filter
				choiceWidget := widget.NewSelect(filter.GetStringCondition(), func(s string) {
					c.condition = s
				})
				choiceWidget.SetSelectedIndex(0)
				c.container.Objects[1] = choiceWidget
				switch filter.Type {
				case jsonfilter.Int:
					c.container.Objects[2] = customs.NewNumericalEntry()
				default:
					c.container.Objects[2] = widget.NewEntry()
				}
				value.Enable()
				break
			}
		}
	}

	c.ExtendBaseWidget(c)

	return c

}

func (c *ConditionalWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(c.container)
}

func (c *ConditionalWidget) MinSize() fyne.Size {
	return fyne.NewSize(30, 40)
}
