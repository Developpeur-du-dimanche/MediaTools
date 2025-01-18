package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"
)

type FileInfoComponent struct {
	widget.BaseWidget
	file medias.FfprobeResult
}

func NewFileInfoComponent(file medias.FfprobeResult) *FileInfoComponent {
	fic := &FileInfoComponent{
		file: file,
	}
	fic.ExtendBaseWidget(fic)
	return fic
}

func (fic *FileInfoComponent) CreateRenderer() fyne.WidgetRenderer {
	size := fyne.CurrentApp().Driver().AllWindows()[0].Canvas().Size()

	b := container.NewBorder(
		container.NewVBox(
			widget.NewLabel("Folder: "+fic.file.Format.Filename),
			widget.NewLabel("Duration: "+fic.file.Format.DurationSeconds.String()),
		),
		nil,
		nil,
		nil,
	)
	b.Resize(fyne.NewSize(size.Width-150, size.Height-150))

	return widget.NewSimpleRenderer(b)
}
