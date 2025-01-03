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
		for _, filter := range filters.Filters {
			if filter.Name == s {
				c.choice = filter
				choiceWidget := widget.NewSelect(filter.GetStringCondition(), func(s string) {
					c.condition = s
				})
				choiceWidget.SetSelectedIndex(0)
				c.container.Objects[1] = choiceWidget

				if filter.HasDefaultValues() {
					value.Enable()
					c.value = filter.GetDefaultValues()[0]
					valueWidget := widget.NewSelect(filter.GetDefaultValues(), func(s string) {
						c.value = s
					})
					valueWidget.SetSelectedIndex(0)
					c.container.Objects[2] = valueWidget
					break
				}

				switch filter.Type {
				case jsonfilter.Int:
					c.container.Objects[2] = customs.NewNumericalEntry()
					c.container.Objects[2].(*customs.NumericalEntry).OnChanged = func(s string) {
						c.value = s
					}
				default:
					c.container.Objects[2] = widget.NewEntry()
					c.container.Objects[2].(*widget.Entry).OnChanged = func(s string) {
						c.value = s
					}
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
