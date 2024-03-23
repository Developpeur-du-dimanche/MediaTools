package view

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/components"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/list"
	"github.com/kbinani/screenshot"
)

type HomeView struct {
	c        chan string
	listFile *list.List
	window   fyne.Window
	list     *widget.List
}

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
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(home.listFile.GetItem(i))
		},
	)

	home.list.Resize(fyne.NewSize(float32(screen.Dx()/4), float32(screen.Dy()/2)))

	home.window.SetContent(home.Content())
	home.window.SetMainMenu(home.GetMainMenu())

	fmt.Printf("listFile address: %p\n", &home.listFile)

	go func() {
		for file := range home.c {
			fmt.Println(file)
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

		for _, file := range list {
			h.c <- file.Path()
		}
	}, *h.GetWindow())
	size := (*h.GetWindow()).Canvas().Size()
	dialog.Resize(fyne.NewSize(size.Width-150, size.Height-150))
	return dialog
}
