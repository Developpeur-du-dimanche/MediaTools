package view

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/components"
	"github.com/kbinani/screenshot"
)

type HomeView struct {
	window fyne.Window
	list   *components.FileListComponent
}

var acceptedExtensions = []string{".mp4", ".avi", ".mkv", ".mov", ".flv", ".wmv", ".webm", ".mpg", ".mpeg", ".wav", ".flac", ".ogg"}

func NewHomeView(app fyne.App) View {

	home := &HomeView{
		window: app.NewWindow("MediaTools"),
	}

	screen := screenshot.GetDisplayBounds(0)
	home.window.Resize(fyne.NewSize(float32(screen.Dx()/2), float32(screen.Dy()/2)))

	list := components.NewFileListComponent(&home.window)

	home.list = list

	home.window.SetContent(home.Content())
	home.window.SetMainMenu(home.GetMainMenu())

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
				h.list.Clear()
			}),
		),
		nil,
		nil,
		nil,
		h.list,
	)

	c.Resize((h.window.Canvas().Size()))
	layout := container.NewAdaptiveGrid(2, c, container.NewAppTabs(
		container.NewTabItem("Filter", components.NewFilterComponent(&h.window, h.list.GetFiles()).Content()),
		container.NewTabItem("Track Remover", components.NewTrackRemoverComponent(&h.window, h.list.GetFiles()).Content()),
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
		h.list.AddFile(uc.URI().Path())
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

					if isValidExtension(filepath.Ext(path)) {
						h.list.AddFile(path)
					}
					return nil
				})
			} else {
				scanFolder.SetText("Scanning folder: " + file.Path())
				fmt.Printf("ext: %s\n", file.Extension())
				if isValidExtension(strings.ToLower(file.Extension())) {
					h.list.AddFile(file.Path())
				}
			}

		}
		dialog.Hide()
	}, *h.GetWindow())
	size := (*h.GetWindow()).Canvas().Size()
	dialog.Resize(fyne.NewSize(size.Width-150, size.Height-150))
	return dialog
}

func isValidExtension(extension string) bool {
	for _, ext := range acceptedExtensions {
		if ext == extension {
			return true
		}
	}
	return false
}
