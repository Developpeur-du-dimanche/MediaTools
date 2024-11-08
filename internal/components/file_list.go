package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/fileinfo"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/list"
)

type FileListComponent struct {
	widget.BaseWidget
	files       *list.List[*fileinfo.FileInfo]
	list        *widget.List
	c           chan string
	OnFileClick func(file *fileinfo.FileInfo)
}

func NewFileListComponent(parent *fyne.Window) *FileListComponent {
	files := list.NewList[*fileinfo.FileInfo]()
	list := new(widget.List)
	c := &FileListComponent{
		files:       files,
		list:        list,
		c:           make(chan string),
		OnFileClick: func(file *fileinfo.FileInfo) {},
	}

	c.list = widget.NewList(
		func() int {
			return files.GetLength()
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewButtonWithIcon("", theme.DeleteIcon(), nil), widget.NewButtonWithIcon("", theme.InfoIcon(), nil), widget.NewLabel(""))
		},
		func(i widget.ListItemID, item fyne.CanvasObject) {

			item.(*fyne.Container).Objects[0].(*widget.Button).OnTapped = func() {
				files.RemoveItem(files.GetItem(i))
				c.list.Refresh()
			}
			item.(*fyne.Container).Objects[1].(*widget.Button).OnTapped = func() {
				d := dialog.NewCustom("Info", "Close", NewFileInfoComponent(c.files.GetItem(i)), *parent)
				d.Resize(fyne.NewSize(400, 400))
				d.Show()
			}
			item.(*fyne.Container).Objects[2].(*widget.Label).SetText(files.GetItem(i).Filename)
		},
	)

	c.list.OnSelected = func(id widget.ListItemID) {
		c.OnFileClick(files.GetItem(id))
	}

	c.ExtendBaseWidget(c)

	go func() {
		for file := range c.c {
			f, err := fileinfo.NewFileInfo(file)
			if err != nil {
				dialog.ShowError(err, *parent)
				continue
			}

			files.AddItem(f)
			c.list.Refresh()
		}
	}()

	return c
}

func (f *FileListComponent) AddFile(file string) {
	f.c <- file
}

func (f *FileListComponent) Clear() {
	f.files.Clear()
	f.list.Refresh()
}

func (f *FileListComponent) RemoveFile(file *fileinfo.FileInfo) {
	f.files.RemoveItem(file)
}

func (f *FileListComponent) GetFiles() *list.List[*fileinfo.FileInfo] {
	return f.files
}

func (f *FileListComponent) Update() {
	f.Refresh()
}

func (f *FileListComponent) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(f.list)
}
