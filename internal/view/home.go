package view

import (
	"errors"
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
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/fileinfo"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/list"
	"github.com/kbinani/screenshot"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type HomeView struct {
	c        chan string
	listFile *list.List[*fileinfo.FileInfo]
	window   fyne.Window
	list     *widget.List
}

var acceptedExtensions = "mp4,avi,mkv,mov,flv,wmv,webm,mpg,mpeg,mp3,wav,flac,ogg"

func NewHomeView(app fyne.App) View {

	home := &HomeView{
		listFile: list.NewList[*fileinfo.FileInfo](),
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
				tree := FileInfoToTree(home.listFile.GetItem(i))
				if tree == nil {
					dialog.ShowError(errors.New("error while getting file information"), *home.GetWindow())
					return
				}
				d := dialog.NewCustom("Info", "Close", tree, *home.GetWindow())
				d.Resize(fyne.NewSize(400, 400))
				d.Show()
			}
			item.(*fyne.Container).Objects[2].(*widget.Label).SetText(home.listFile.GetItem(i).Path)
		},
	)

	home.list.Resize(fyne.NewSize(float32(screen.Dx()/4), float32(screen.Dy()/2)))

	home.window.SetContent(home.Content())
	home.window.SetMainMenu(home.GetMainMenu())

	fmt.Printf("listFile address: %p\n", &home.listFile)

	go func() {
		for file := range home.c {
			f, err := fileinfo.NewFileInfo(file)
			if err != nil {
				dialog.ShowError(err, *home.GetWindow())
				continue
			}
			home.listFile.AddItem(f)
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

func FileInfoToTree(file *fileinfo.FileInfo) *widget.Tree {
	tree := widget.NewTree(
		func(id widget.TreeNodeID) []widget.TreeNodeID {
			switch id {
			case "":
				return []widget.TreeNodeID{"streams", "format"}
			case "streams":
				return []widget.TreeNodeID{"Video", "Audio", "Subtitle"}
			case "Video":
				// one node per video stream
				nodes := []widget.TreeNodeID{}
				for i := 0; i < len(*file.VideoStreams); i++ {
					nodes = append(nodes, "Video "+fmt.Sprint(i))
				}
				return nodes
			case "Audio":
				// one node per audio stream
				nodes := []widget.TreeNodeID{}
				for i := 0; i < len(*file.AudioStreams); i++ {
					nodes = append(nodes, "Audio "+fmt.Sprint(i))
				}
				return nodes
			case "Subtitle":
				// one node per subtitle stream
				nodes := []widget.TreeNodeID{}
				for i := 0; i < len(*file.SubtitleStreams); i++ {
					nodes = append(nodes, "Subtitle "+fmt.Sprint(i))
				}
				return nodes
			case "format":
				return []widget.TreeNodeID{"format_name", "duration", "size"}
			}

			if strings.HasPrefix(id, "Video") {
				nodes := []widget.TreeNodeID{}
				for i := 0; i < len(*file.VideoStreams); i++ {
					nodes = append(nodes, "v"+fmt.Sprint(i)+"_codec_name", "v"+fmt.Sprint(i)+"_width", "v"+fmt.Sprint(i)+"_height", "v"+fmt.Sprint(i)+"_bit_rate")
				}
				return nodes
			}

			if strings.HasPrefix(id, "Audio") {
				nodes := []widget.TreeNodeID{}
				for i := 0; i < len(*file.AudioStreams); i++ {
					nodes = append(nodes, "a"+fmt.Sprint(i)+"_codec_name", "a"+fmt.Sprint(i)+"_channels", "a"+fmt.Sprint(i)+"_bit_rate", "a"+fmt.Sprint(i)+"_title")
				}
				return nodes
			}

			if strings.HasPrefix(id, "Subtitle") {
				nodes := []widget.TreeNodeID{}
				for i := 0; i < len(*file.SubtitleStreams); i++ {
					nodes = append(nodes, "s"+fmt.Sprint(i)+"_codec_name", "s"+fmt.Sprint(i)+"_title")
				}
				return nodes
			}

			return []string{}
		},
		func(id widget.TreeNodeID) bool {
			if id == "" || id == "format" || id == "streams" || id == "Video" || id == "Audio" || id == "Subtitle" {
				return true
			}

			for i := 0; i < len(*file.SubtitleStreams); i++ {
				if id == "Video "+fmt.Sprint(i) {
					return true
				}
			}

			for i := 0; i < len(*file.AudioStreams); i++ {
				if id == "Audio "+fmt.Sprint(i) {
					return true
				}
			}

			for i := 0; i < len(*file.SubtitleStreams); i++ {
				if id == "Subtitle "+fmt.Sprint(i) {
					return true
				}
			}

			return false

		},
		func(b bool) fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(id widget.TreeNodeID, isBranch bool, co fyne.CanvasObject) {
			if isBranch {
				co.(*widget.Label).SetText(id)
				return
			}

			switch id {
			case "format_name":
				co.(*widget.Label).SetText("Format name: " + file.Info.Format.FormatName)
				return
			case "duration":
				co.(*widget.Label).SetText("Duration: " + fmt.Sprint(file.Info.Format.Duration()))
				return
			case "size":
				co.(*widget.Label).SetText("Size: " + file.Info.Format.Size)
				return
			}

			var stream *ffprobe.Stream
			i := int(id[1] - '0')

			switch id[0] {
			case 'v':
				stream = &(*file.VideoStreams)[i]
			case 'a':
				stream = &(*file.AudioStreams)[i]
			case 's':
				stream = &(*file.SubtitleStreams)[i]
			}

			request := strings.TrimPrefix(id, string(id[0])+fmt.Sprint(i)+"_")

			switch request {
			case "codec_name":
				co.(*widget.Label).SetText("Codec name: " + stream.CodecName)
				return
			case "width":
				co.(*widget.Label).SetText("Width: " + fmt.Sprint(stream.Width))
				return
			case "height":
				co.(*widget.Label).SetText("Height: " + fmt.Sprint(stream.Height))
				return
			case "bit_rate":
				co.(*widget.Label).SetText("Bit rate: " + fmt.Sprint(stream.BitRate))
				return
			case "channels":
				co.(*widget.Label).SetText("Channels: " + fmt.Sprint(stream.Channels))
				return
			case "title":
				co.(*widget.Label).SetText("Title: " + stream.Tags.Title)
				return
			}

		},
	)

	return tree
}
