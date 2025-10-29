package components

import (
	"fmt"
	"path/filepath"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"
)

type FileInfoComponent struct {
	widget.BaseWidget
	file    *medias.FfprobeResult
	appTabs []container.AppTabs
	window  fyne.Window
}

func NewFileInfoComponent(file *medias.FfprobeResult, window fyne.Window) *FileInfoComponent {
	fic := &FileInfoComponent{
		file:    file,
		appTabs: []container.AppTabs{},
		window:  window,
	}
	fic.ExtendBaseWidget(fic)
	return fic
}

func (fic *FileInfoComponent) CreateRenderer() fyne.WidgetRenderer {
	// File header info with better formatting
	filename := filepath.Base(fic.file.Format.Filename)
	directory := filepath.Dir(fic.file.Format.Filename)

	duration := fic.file.Format.DurationSeconds.String()
	size := formatSizeString(fic.file.Format.Size)
	bitrate := formatBitrateString(fic.file.Format.BitRate)

	// File info section with grid layout
	fileInfoGrid := container.NewVBox(
		widget.NewLabelWithStyle("File Information", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		fic.createInfoRow("Name:", filename),
		fic.createInfoRow("Location:", directory),
		fic.createInfoRow("Duration:", duration),
		fic.createInfoRow("Size:", size),
		fic.createInfoRow("Bitrate:", bitrate),
	)

	// Stream tabs
	tabs := container.NewAppTabs()

	if len(fic.file.Videos) != 0 {
		tabs.Append(container.NewTabItem("Video Streams", fic.createVideoTabs()))
	}

	if len(fic.file.Audios) != 0 {
		tabs.Append(container.NewTabItem("Audio Streams", fic.createAudioTabs()))
	}

	if len(fic.file.Subtitles) != 0 {
		tabs.Append(container.NewTabItem("Subtitle Streams", fic.createSubtitleTabs()))
	}

	// Main layout
	content := container.NewBorder(
		fileInfoGrid,
		nil,
		nil,
		nil,
		tabs,
	)

	content.Resize(fyne.NewSize(fic.window.Canvas().Size().Width-150, fic.window.Canvas().Size().Height-150))

	return widget.NewSimpleRenderer(content)
}

func (fic *FileInfoComponent) createInfoRow(label, value string) fyne.CanvasObject {
	return container.NewHBox(
		widget.NewLabelWithStyle(label, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(value),
	)
}

func formatSizeString(sizeStr string) string {
	if sizeStr == "" {
		return "N/A"
	}
	// Parse string to int64
	bytes, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		return sizeStr // Return original if parsing fails
	}

	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func formatBitrateString(bitrateStr string) string {
	if bitrateStr == "" {
		return "N/A"
	}
	// Parse string to int64
	bitrate, err := strconv.ParseInt(bitrateStr, 10, 64)
	if err != nil {
		return bitrateStr // Return original if parsing fails
	}

	if bitrate == 0 {
		return "N/A"
	}
	return fmt.Sprintf("%.2f Mbps", float64(bitrate)/1000000)
}

func (fic *FileInfoComponent) createVideoTabs() *container.AppTabs {
	videoAppTabs := container.NewAppTabs()

	for _, stream := range fic.file.Videos {
		streamInfo := container.NewVBox(
			widget.NewLabelWithStyle(fmt.Sprintf("Video Stream #%d", stream.StreamIndex), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewSeparator(),
			fic.createInfoRow("Codec:", stream.CodecName),
			fic.createInfoRow("Resolution:", fmt.Sprintf("%dx%d", stream.Width, stream.Height)),
		)

		videoAppTab := container.NewTabItem(fmt.Sprintf("Stream #%d", stream.StreamIndex), streamInfo)
		videoAppTabs.Append(videoAppTab)
	}

	return videoAppTabs
}

func (fic *FileInfoComponent) createAudioTabs() *container.AppTabs {
	audioAppTabs := container.NewAppTabs()

	for _, stream := range fic.file.Audios {
		language := stream.Language
		if language == "" {
			language = "Unknown"
		}

		streamInfo := container.NewVBox(
			widget.NewLabelWithStyle(fmt.Sprintf("Audio Stream #%d", stream.StreamIndex), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewSeparator(),
			fic.createInfoRow("Codec:", stream.CodecName),
			fic.createInfoRow("Channels:", fmt.Sprintf("%d", stream.Channels)),
			fic.createInfoRow("Language:", language),
		)

		audioAppTab := container.NewTabItem(fmt.Sprintf("Stream #%d", stream.StreamIndex), streamInfo)
		audioAppTabs.Append(audioAppTab)
	}

	return audioAppTabs
}

func (fic *FileInfoComponent) createSubtitleTabs() *container.AppTabs {
	subtitleAppTabs := container.NewAppTabs()

	for _, stream := range fic.file.Subtitles {
		language := stream.Language
		if language == "" {
			language = "Unknown"
		}

		streamInfo := container.NewVBox(
			widget.NewLabelWithStyle(fmt.Sprintf("Subtitle Stream #%d", stream.StreamIndex), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewSeparator(),
			fic.createInfoRow("Codec:", stream.CodecName),
			fic.createInfoRow("Language:", language),
		)

		subtitleAppTab := container.NewTabItem(fmt.Sprintf("Stream #%d", stream.StreamIndex), streamInfo)
		subtitleAppTabs.Append(subtitleAppTab)
	}

	return subtitleAppTabs
}
