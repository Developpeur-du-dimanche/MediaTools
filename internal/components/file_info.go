package components

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"
)

type FileInfoComponent struct {
	widget.BaseWidget
	file    medias.FfprobeResult
	appTabs []container.AppTabs
	window  fyne.Window
}

func NewFileInfoComponent(file medias.FfprobeResult, window fyne.Window) *FileInfoComponent {
	fic := &FileInfoComponent{
		file:    file,
		appTabs: []container.AppTabs{},
		window:  window,
	}
	fic.ExtendBaseWidget(fic)
	return fic
}

func (fic *FileInfoComponent) CreateRenderer() fyne.WidgetRenderer {

	tabs := container.NewAppTabs()

	if len(fic.file.Videos) != 0 {
		tabs.Append(container.NewTabItem("Video", fic.createVideoTabs()))
	}

	if len(fic.file.Audios) != 0 {
		tabs.Append(container.NewTabItem("Audio", fic.createAudioTabs()))
	}

	if len(fic.file.Subtitles) != 0 {
		tabs.Append(container.NewTabItem("Subtitle", fic.createSubtitleTabs()))
	}

	b := container.NewBorder(
		container.NewVBox(
			widget.NewLabel("Folder: "+fic.file.Format.Filename),
			widget.NewLabel("Duration: "+fic.file.Format.DurationSeconds.String()),
			tabs,
		),
		nil,
		nil,
		nil,
	)

	b.Refresh()
	b.Resize(fyne.NewSize(fic.window.Canvas().Size().Width-150, fic.window.Canvas().Size().Height-150))

	return widget.NewSimpleRenderer(b)
}

func (fic *FileInfoComponent) createVideoTabs() *container.AppTabs {
	videoAppTabs := container.NewAppTabs()

	for i, stream := range fic.file.Videos {
		videoAppTab := container.NewTabItem("Video "+strconv.Itoa(i), container.NewVBox(
			widget.NewLabel("Codec: "+stream.CodecName),
			widget.NewLabel("Resolution: "+strconv.Itoa(stream.Width)+"x"+strconv.Itoa(stream.Height)),
		))

		videoAppTabs.Append(videoAppTab)

	}

	return videoAppTabs

}

func (fic *FileInfoComponent) createAudioTabs() *container.AppTabs {
	audioAppTabs := container.NewAppTabs()

	for i, stream := range fic.file.Audios {
		audioAppTab := container.NewTabItem("Audio "+strconv.Itoa(i), container.NewVBox(
			widget.NewLabel("Codec: "+stream.CodecName),
			widget.NewLabel("Channels: "+strconv.Itoa(stream.Channels)),
		))

		audioAppTabs.Append(audioAppTab)

	}

	return audioAppTabs

}

func (fic *FileInfoComponent) createSubtitleTabs() *container.AppTabs {
	subtitleAppTabs := container.NewAppTabs()

	for i, stream := range fic.file.Subtitles {
		subtitleAppTab := container.NewTabItem("Subtitle "+strconv.Itoa(i), container.NewVBox(
			widget.NewLabel("Codec: "+stream.CodecName),
		))

		subtitleAppTabs.Append(subtitleAppTab)

	}

	return subtitleAppTabs

}
