package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/services"
)

type LastScanSelector struct {
	widget.BaseWidget

	selector       *widget.Select
	historyService *services.HistoryService
	onSelect       func(path string)
}

func NewLastScanSelector(historyService *services.HistoryService, onSelect func(path string)) *LastScanSelector {
	folderHistory := historyService.GetHistory()

	selector := widget.NewSelect(folderHistory, onSelect)

	if len(folderHistory) == 0 {
		selector.Disable()
	}

	lss := &LastScanSelector{
		selector:       selector,
		historyService: historyService,
		onSelect:       onSelect,
	}
	lss.ExtendBaseWidget(lss)
	return lss
}

func (lss *LastScanSelector) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(lss.selector)
}

// Refresh updates the selector with the latest history
func (lss *LastScanSelector) Refresh() {
	lss.selector.Options = lss.historyService.GetHistory()

	if len(lss.selector.Options) > 0 {
		lss.selector.Enable()
	} else {
		lss.selector.Disable()
	}

	lss.selector.Refresh()
}
