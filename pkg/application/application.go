package application

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/kbinani/screenshot"
)

type ApplicationInterface interface {
	GetApp() fyne.App
	GetWindow() WindowInterface
	GetDisplaySize() (float32, float32)
	SetMainMenu(menu *fyne.MainMenu)
}

type Application struct {
	App    fyne.App
	Window WindowInterface
	height float32
	width  float32
	files  []string
}

func NewApp(metadata *fyne.AppMetadata) ApplicationInterface {

	app.SetMetadata(*metadata)
	a := app.New()

	w := a.NewWindow("Mediatools")

	screenSize := screenshot.GetDisplayBounds(0)

	x := float32(screenSize.Dx() / 2)
	y := float32(screenSize.Dy() / 2)

	openFileButton := widget.NewButton("Open File", nil)
	openFolderButton := widget.NewButton("Open Folder", nil)
	removeAllButton := widget.NewButton("Remove All", nil)

	appli := &Application{
		App:    a,
		height: y,
		width:  x,
		files:  []string{},
	}

	appli.files = append(appli.files, "test")

	fileList := widget.NewList(
		func() int {
			return len(appli.files)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, item fyne.CanvasObject) {
			item.(*widget.Label).SetText(appli.files[i])
		},
	)

	appli.Window = newWindow(&w, NewLayout(
		openFolderButton,
		openFileButton,
		fileList,
	), &fyne.Size{Width: x, Height: y})

	openFileButton.OnTapped = func() {
		appli.GetWindow().NewFileOpen(func(f fyne.URIReadCloser, e error) {
			// if err != nil {
			// 	dialog.ShowError(err, w)
			// 	return
			// }
			if e != nil {
				dialog.ShowError(e, w)
				return
			}

			if f == nil {
				return
			}

			appli.files = append(appli.files, f.URI().String())
			fileList.Refresh()

		}, nil).Show()
	}

	openFolderButton.OnTapped = func() {
		appli.GetWindow().NewFolderOpen(func(f fyne.ListableURI, e error) {
			fmt.Println("Open Folder")
			if e != nil {
				dialog.ShowError(e, w)
				return
			}

			if f == nil {
				dialog.ShowError(fmt.Errorf("No folder selected"), w)
				return
			}

			loadingLabel := widget.NewLabel("Loading...")

			// show loading dialog
			custom := dialog.NewCustomWithoutButtons("Loading", loadingLabel, w)

			lists, err := f.List()

			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			custom.Show()

			for _, file := range lists {
				loadingLabel.SetText("Loading... " + file.Name())
				appli.files = append(appli.files, file.Path())
			}

			custom.Hide()

		}, nil).Show()
	}

	removeAllButton.OnTapped = func() {
		appli.files = []string{}
		fileList.Refresh()
	}

	openFolderButton.OnTapped = func() {
		appli.GetWindow().NewFolderOpen(func(fyne.ListableURI, error) {
			// TODO
		}, nil).Show()
	}

	w.Resize(fyne.NewSize(x, y))

	return appli

}

func (a *Application) GetApp() fyne.App {
	return a.App
}

func (a *Application) GetWindow() WindowInterface {
	return a.Window
}

func (a *Application) GetDisplaySize() (float32, float32) {
	return a.width, a.height
}

func (a *Application) SetMainMenu(menu *fyne.MainMenu) {
	a.Window.SetMainMenu(menu)
}
