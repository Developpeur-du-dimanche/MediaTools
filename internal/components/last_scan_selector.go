package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type LastScanSelector struct {
	widget.BaseWidget

	selector *widget.Select

	onSelect func(path string)
}

func NewLastScanSelector(onSelect func(path string)) *LastScanSelector {
	app := fyne.CurrentApp()
	folderHistory := app.Preferences().StringListWithFallback("last_scan_selector", []string{})

	selector := widget.NewSelect(folderHistory, onSelect)

	if len(folderHistory) > 0 {
		selector.Disable()
	}

	lss := &LastScanSelector{
		selector: selector,
		onSelect: onSelect,
	}
	lss.ExtendBaseWidget(lss)
	return lss
}

func (lss *LastScanSelector) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(lss.selector)
}

func (lss *LastScanSelector) AddFolder(path string) {

	for _, p := range lss.selector.Options {
		if p == path {
			return
		}
	}

	lss.selector.Options = append(lss.selector.Options, path)
	fyne.CurrentApp().Preferences().SetStringList("last_scan_selector", lss.selector.Options)
	lss.selector.Refresh()
}
