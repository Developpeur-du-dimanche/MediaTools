package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/helper"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/fileinfo"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/list"
)

type FileListComponent struct {
	widget.BaseWidget
	files       *list.List[*helper.FileMetadata]
	list        *widget.List
	c           chan string
	OnFileClick func(file *helper.FileMetadata)
}

func NewFileListComponent(parent *fyne.Window) *FileListComponent {
	files := list.NewList[*helper.FileMetadata]()
	list := new(widget.List)
	c := &FileListComponent{
		files:       files,
		list:        list,
		c:           make(chan string),
		OnFileClick: func(file *helper.FileMetadata) {},
	}

	c.list = widget.NewList(
		func() int {
			return files.GetLength()
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewButtonWithIcon("", theme.DeleteIcon(), nil),
				widget.NewButtonWithIcon("", theme.InfoIcon(), nil),
				widget.NewButtonWithIcon("", theme.DocumentIcon(), nil),
				widget.NewLabel(""),
			)
		},
		func(i widget.ListItemID, item fyne.CanvasObject) {

			item.(*fyne.Container).Objects[0].(*widget.Button).OnTapped = func() {
				files.RemoveItem(files.GetItem(i))
				c.list.Refresh()
			}
			item.(*fyne.Container).Objects[1].(*widget.Button).OnTapped = func() {
				d := dialog.NewCustom("Info", lang.L("close"), NewFileInfoComponent(parent, c.files.GetItem(i)), *parent)
				d.Resize(fyne.NewSize(400, 400))
				d.Show()
			}

			item.(*fyne.Container).Objects[2].(*widget.Button).OnTapped = func() {
				d := dialog.NewCustom("Edit", lang.L("close"), NewFileEditComponent(parent, c.files.GetItem(i)), *parent)
				d.Resize(fyne.NewSize(400, 400))
				d.Show()
			}

			item.(*fyne.Container).Objects[3].(*widget.Label).SetText(files.GetItem(i).FileName)
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

			metadata, err := helper.NewFileMetadata(f.GetInfo())

			if err != nil {
				dialog.ShowError(err, *parent)
				continue
			}

			files.AddItem(metadata)
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

func (f *FileListComponent) RemoveFile(file *helper.FileMetadata) {
	f.files.RemoveItem(file)
}

func (f *FileListComponent) GetFiles() *list.List[*helper.FileMetadata] {
	return f.files
}

func (f *FileListComponent) Update() {
	f.Refresh()
}

func (f *FileListComponent) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(f.list)
}
