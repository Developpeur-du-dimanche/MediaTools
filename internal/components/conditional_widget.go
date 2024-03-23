package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/filter"
)

type ConditionalWidget struct {
	widget.BaseWidget
	key       *widget.Select
	choice    filter.ConditionContract
	container *fyne.Container
}

var conditions = []filter.ConditionContract{
	filter.NewContainerFilter(),
	filter.NewAudioLanguageFilter(),
	filter.NewBitrateFilter(),
	filter.NewSubtitleForcedFilter(),
	filter.NewSubtitleLanguageFilter(),
	filter.NewSubtitleTitleFilter(),
	filter.NewSubtitleCodecFilter(),
	filter.NewVideoTitleFilter(),
}

var choices = make([]string, len(conditions))

func NewConditionalWidget() *ConditionalWidget {
	for i, condition := range conditions {
		choices[i] = condition.Name()
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
		for _, cond := range conditions {
			if cond.Name() == s {
				c.choice = cond.New()
				choiceWidget := widget.NewSelect(cond.GetPossibleConditions(), func(s string) {
					c.choice.SetCondition(s)
				})
				c.container.Objects[1] = choiceWidget
				c.container.Objects[2] = c.choice.GetEntry()
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
