package main

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"

	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func init() {
	application.RegisterEvent[string]("time")
}

type FileService struct {
}

func (f *FileService) ShowFileServiceDialog() (string, error) {
	app := application.Get()

	if app == nil {
		return "", errors.New("cannot get current app")
	}

	folder, err := app.Dialog.OpenFile().
		CanChooseFiles(false).
		CanChooseDirectories(true).
		SetTitle("Open Directory").PromptForSingleSelection()

	if err != nil {
		return "", err
	}

	fmt.Println("choosing " + folder)

	return folder, nil
}

func (f *FileService) GetFiles(directory string) []string {
	var files = []string{}
	filepath.Walk(directory, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return errors.New(path + " is a directory")
		}

		files = append(files, path)
		return nil
	})

	return files

}

func main() {
	app := application.New(application.Options{
		Name:        "MediaTools",
		Description: "A demo of using raw HTML & CSS",
		Services: []application.Service{
			application.NewService(&GreetService{}),
			application.NewService(&FileService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Window 1",
		// Window sized to the golden ratio (1000 / 618 ≈ 1.618).
		Width:  1000,
		Height: 618,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(6, 7, 15),
		URL:              "/",
	})

	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			app.Event.Emit("time", now)
			time.Sleep(time.Second)
		}
	}()

	err := app.Run()

	if err != nil {
		log.Fatal(err)
	}
}
