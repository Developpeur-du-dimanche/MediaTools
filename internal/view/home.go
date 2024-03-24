package view

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/components"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/list"
	"github.com/kbinani/screenshot"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type HomeView struct {
	c        chan string
	listFile *list.List
	window   fyne.Window
	list     *widget.List
}

var acceptedExtensions = "mp4,avi,mkv,mov,flv,wmv,webm,mpg,mpeg,mp3,wav,flac,ogg"

func NewHomeView(app fyne.App) View {

	home := &HomeView{
		listFile: list.NewList(),
		c:        make(chan string),
	}

	home.window = app.NewWindow("MediaTools")
	screen := screenshot.GetDisplayBounds(0)
	home.window.Resize(fyne.NewSize(float32(screen.Dx()/2), float32(screen.Dy()/2)))
	home.list = widget.NewList(
		func() int {
			return home.listFile.GetLength()
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewButtonWithIcon("", theme.DeleteIcon(), nil), widget.NewButtonWithIcon("", theme.InfoIcon(), nil), widget.NewLabel(""))
		},
		func(i widget.ListItemID, item fyne.CanvasObject) {

			item.(*fyne.Container).Objects[0].(*widget.Button).OnTapped = func() {
				home.listFile.RemoveItem(home.listFile.GetItem(i))
				home.list.Refresh()
			}
			item.(*fyne.Container).Objects[1].(*widget.Button).OnTapped = func() {
				dialog.ShowInformation("Info", home.getFileInformation(home.listFile.GetItem(i)), *home.GetWindow())
			}
			item.(*fyne.Container).Objects[2].(*widget.Label).SetText(home.listFile.GetItem(i))
		},
	)

	home.list.Resize(fyne.NewSize(float32(screen.Dx()/4), float32(screen.Dy()/2)))

	home.window.SetContent(home.Content())
	home.window.SetMainMenu(home.GetMainMenu())

	fmt.Printf("listFile address: %p\n", &home.listFile)

	go func() {
		for file := range home.c {
			home.listFile.AddItem(file)
			home.list.Refresh()
		}
	}()

	return home
}

func (h HomeView) Content() fyne.CanvasObject {

	c := container.NewBorder(
		container.NewHBox(
			widget.NewButtonWithIcon("", theme.FolderIcon(), func() {
				h.OpenFolderDialog().Show()
			}),
			widget.NewButtonWithIcon("", theme.FileIcon(), func() {
				h.OpenFileDialog().Show()
			}),
			widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
				h.listFile.Clear()
				h.list.Refresh()
			}),
		),
		nil,
		nil,
		nil,
		h.list,
	)

	c.Resize((h.window.Canvas().Size()))

	// print address of listFile
	fmt.Printf("listFile address: %p\n", &h.listFile)

	layout := container.NewAdaptiveGrid(2, c, container.NewAppTabs(
		container.NewTabItem("Filter", components.NewFilterComponent(&h.window, h.listFile).Content()),
		container.NewTabItem("Track Remover", components.NewTrackRemoverComponent().Content()),
	))

	return layout
}

func (h *HomeView) GetWindow() *fyne.Window {
	return &h.window
}

func (h *HomeView) GetMainMenu() *fyne.MainMenu {
	return fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Settings", nil),
		),
	)
}

func (h *HomeView) ShowAndRun() {
	h.window.ShowAndRun()
}

func (h *HomeView) OpenFileDialog() *dialog.FileDialog {
	dialog := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, *h.GetWindow())
			return
		}

		if uc == nil {
			return
		}
		h.c <- uc.URI().Path()
	}, *h.GetWindow())
	size := (*h.GetWindow()).Canvas().Size()
	dialog.Resize(fyne.NewSize(size.Width-150, size.Height-150))
	return dialog

}

func (h *HomeView) OpenFolderDialog() *dialog.FileDialog {
	dialog := dialog.NewFolderOpen(func(lu fyne.ListableURI, err error) {
		if err != nil {
			dialog.ShowError(err, *h.GetWindow())
			return
		}

		if lu == nil {
			return
		}

		list, err := lu.List()
		if err != nil {
			dialog.ShowError(err, *h.GetWindow())
			return
		}

		scanFolder := widget.NewLabel("Scanning folder: " + lu.Path())
		dialog := dialog.NewCustomWithoutButtons("Scanning folder", scanFolder, *h.GetWindow())
		dialog.Show()

		for _, file := range list {
			if i, err := os.Stat(file.Path()); err == nil && i.IsDir() {
				filepath.WalkDir(file.Path(), func(path string, d os.DirEntry, err error) error {

					scanFolder.SetText("Scanning folder: " + path)

					if err != nil {
						return err
					}

					if filepath.Ext(path) == "" {
						return nil
					}

					if strings.Contains(acceptedExtensions, file.Extension()) {
						h.c <- path
					}
					return nil
				})
			} else {
				scanFolder.SetText("Scanning folder: " + file.Path())
				if strings.Contains(acceptedExtensions, file.Extension()) {
					h.c <- file.Path()
				}
			}

		}
		dialog.Hide()
	}, *h.GetWindow())
	size := (*h.GetWindow()).Canvas().Size()
	dialog.Resize(fyne.NewSize(size.Width-150, size.Height-150))
	return dialog
}

func (h *HomeView) getFileInformation(file string) string {
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()

	data, err := ffprobe.ProbeURL(ctx, file)

	if err != nil {
		dialog.ShowError(err, *h.GetWindow())
		return "error while probing file: " + file + ", " + err.Error()
	}

	return fmt.Sprintf("file: %s\nformat: %s\nsize: %s\nbitrate: %s\n", file, data.Format.FormatName, data.Format.Size, data.Format.BitRate)

}
