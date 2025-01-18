package mediatools

import (
	"context"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/components"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"
	"github.com/kbinani/screenshot"
)

type MediaTools struct {
	app        *fyne.App
	window     fyne.Window
	listView   *components.ListView
	burgerMenu *components.BurgerMenu

	isScanning chan bool
}

func NewMediaTools(app fyne.App) *MediaTools {

	mediaTools := &MediaTools{
		app:    &app,
		window: app.NewWindow("MediaTools"),

		isScanning: make(chan bool),
	}

	screen := screenshot.GetDisplayBounds(0)
	mediaTools.window.Resize(fyne.NewSize(
		float32(screen.Dx()/2), float32(screen.Dy()/2),
	))

	mediaTools.listView = components.NewListView(nil)

	history := components.NewLastScanSelector(func(path string) {
		fmt.Printf("Folder selected: %s\n", path)
	})

	openFolder := components.NewOpenFolder(&mediaTools.window, func(path string) {
		history.AddFolder(path)
	},
		mediaTools.onNewFileDetected,
	)

	openFile := components.NewOpenFile(&mediaTools.window, mediaTools.onNewFileDetected)

	openFolder.OnScanTerminated = func() {
		mediaTools.listView.Refresh()
	}

	openFile.OnScanTerminated = func() {
		mediaTools.listView.Refresh()
	}

	burgerMenu := components.NewBurgerMenu(
		container.NewHBox(
			openFolder,
			openFile,
			widget.NewButtonWithIcon("Clean", theme.DeleteIcon(), func() {
				mediaTools.listView.Clear()
			}),
			history,
		),
		nil, nil, nil, mediaTools.listView, mediaTools.window, func() {
			mediaTools.listView.Refresh()
		})

	mediaTools.window.SetContent(container.NewBorder(
		burgerMenu, nil, nil, nil,
		nil,
	))

	return mediaTools
}

func (mt *MediaTools) Run() {
	mt.window.ShowAndRun()
}

func (mt *MediaTools) onNewFileDetected(path string) {
	mediaInfo, err := getMediaInfo(path)
	if err != nil {
		fmt.Printf("Error while getting media info: %s\n", err)
		return
	}

	mt.listView.AddItem(mediaInfo)

}

func getMediaInfo(path string) (*medias.FfprobeResult, error) {
	ffprobe := medias.NewFfprobe(path,
		medias.FFPROBE_LOGLEVEL_FATAL,
		medias.PRINT_FORMAT_JSON,
		medias.SHOW_FORMAT,
		medias.SHOW_STREAMS,
		medias.EXPERIMENTAL,
	)

	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	return ffprobe.Probe(ctx)

}
