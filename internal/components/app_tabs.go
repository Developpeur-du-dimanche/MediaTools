package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/helper"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/list"
)

type CustomAppTabs struct {
	*container.AppTabs
	window   *fyne.Window
	fileList *list.List[*helper.FileMetadata]
}

type CustomTabItem struct {
	*container.TabItem
}

func NewCustomAppTabs(window *fyne.Window, fileList *list.List[*helper.FileMetadata]) *CustomAppTabs {
	return &CustomAppTabs{
		AppTabs:  container.NewAppTabs(),
		window:   window,
		fileList: fileList,
	}
}

func (c *CustomAppTabs) AddTabItem(name string, callback func(window *fyne.Window, fileList *list.List[*helper.FileMetadata]) Component) *CustomAppTabs {
	c.AppTabs.Append(container.NewTabItem(name, callback(c.window, c.fileList).Content()))
	return c
}

func NewCustomTabItem() *CustomTabItem {
	return &CustomTabItem{}
}
