package components

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/helper"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type FileEditComponent struct {
	widget.BaseWidget
	file   *helper.FileMetadata
	window *fyne.Window
}

func NewFileEditComponent(window *fyne.Window, file *helper.FileMetadata) *FileEditComponent {

	c := &FileEditComponent{
		file:   file,
		window: window,
	}

	c.ExtendBaseWidget(c)
	return c
}

func (f *FileEditComponent) CreateRenderer() fyne.WidgetRenderer {
	size := (*f.window).Canvas().Size()

	// menu avec onglet pour chaques type de stream
	menu := container.NewAppTabs(
		container.NewTabItem("Video", f.ListStreams(f.file.GetVideoStreams())),
		container.NewTabItem("Audio", f.ListStreams(f.file.GetAudioStreams())),
		container.NewTabItem("Subtitle", f.ListStreams(f.file.GetSubtitleStreams())),
	)

	b := container.NewBorder(
		container.NewVBox(
			widget.NewLabel("File: "+f.file.FileName),
		),
		container.NewHBox(
			widget.NewButton("Save", func() {

			}),
		),
		nil,
		nil,
		menu,
	)

	b.Resize(fyne.NewSize(size.Width-150, size.Height-150))

	return widget.NewSimpleRenderer(b)
}

func (f *FileEditComponent) ListStreams(streams []*ffprobe.Stream) *widget.List {
	return widget.NewList(
		func() int {
			return len(streams)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewCheck("", nil),
				widget.NewLabel(""),
			)
		},
		func(i widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[0].(*widget.Check).Checked = true

			var title string
			var err error
			if len(streams[i].TagList) > 0 {
				title, err = streams[i].TagList.GetString("title")
				if err != nil {
					title = "Unknown"
				}
			} else {
				title = "Unknown"
			}

			var codecName string
			if streams[i].CodecName == "" {
				codecName = "Unknown"
			} else {
				codecName = streams[i].CodecName
			}

			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(fmt.Sprintf("%d. %s - %s", i, codecName, title))

			item.Refresh()

		},
	)

}
