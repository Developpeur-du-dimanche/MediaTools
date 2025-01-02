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
	"github.com/Developpeur-du-dimanche/MediaTools/internal/configuration"
	"github.com/kbinani/screenshot"
)

type HomeView struct {
	window        fyne.Window
	list          *components.FileListComponent
	configuration configuration.Configuration
}

func NewHomeView(window fyne.Window, configuration configuration.Configuration) View {

	home := &HomeView{
		window:        window,
		configuration: configuration,
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
		container.NewTabItem("Merge files", components.NewMergeFilesComponent(&h.window, h.list.GetFiles()).Content()),
	))

	return layout
}

func (h *HomeView) GetWindow() *fyne.Window {
	return &h.window
}

func (h *HomeView) GetMainMenu() *fyne.MainMenu {
	extentionsEntry := widget.NewEntry()
	extentionsEntry.SetText(strings.Join(h.configuration.GetList("extension"), ","))
	return fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Settings", func() {
				// create settings view
				view := container.NewBorder(
					widget.NewLabel("Settings"),
					nil,
					nil,
					nil,
					container.NewHSplit(
						widget.NewLabel("Extension:"),
						extentionsEntry,
					),
				)
				// show popup with settings
				popup := dialog.NewCustom("Settings", "Save", view, h.window)

				size := (*h.GetWindow()).Canvas().Size()
				popup.Resize(fyne.NewSize(size.Width-150, size.Height-150))
				popup.Show()

				popup.SetOnClosed(func() {
					ext := strings.Split(strings.Trim(extentionsEntry.Text, " "), ",")
					h.configuration.SaveExtensions(ext)
				})
			}),
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

					if h.configuration.IsValidExtension(filepath.Ext(path)) {
						h.list.AddFile(path)
					}
					return nil
				})
			} else {
				scanFolder.SetText("Scanning folder: " + file.Path())
				fmt.Printf("ext: %s\n", file.Extension())
				if h.configuration.IsValidExtension(strings.ToLower(file.Extension())) {
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
