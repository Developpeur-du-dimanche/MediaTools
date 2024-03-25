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
				return generateStreamNodes(file.VideoStreams, "Video")
			case "Audio":
				return generateStreamNodes(file.AudioStreams, "Audio")
			case "Subtitle":
				return generateStreamNodes(file.SubtitleStreams, "Subtitle")
			case "format":
				return []widget.TreeNodeID{"format_name", "duration", "size"}
			}

			switch {
			case strings.HasPrefix(id, "Video"):
				return generateStreamDetailNodes(id, "v")
			case strings.HasPrefix(id, "Audio"):
				return generateStreamDetailNodes(id, "a")
			case strings.HasPrefix(id, "Subtitle"):
				return generateStreamDetailNodes(id, "s")
			}

			return []string{}
		},
		func(id widget.TreeNodeID) bool {
			switch id {
			case "", "format", "streams", "Video", "Audio", "Subtitle":
				return true
			}

			for i := 0; i < len(*file.VideoStreams); i++ {
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
			case "duration":
				co.(*widget.Label).SetText("Duration: " + fmt.Sprint(file.Info.Format.Duration()))
			case "size":
				co.(*widget.Label).SetText("Size: " + file.Info.Format.Size)
			default:
				stream := getStreamByID(file, id)
				if stream != nil {
					setStreamDetailText(co.(*widget.Label), stream, id)
				}
			}
		},
	)

	return tree
}

func generateStreamNodes(streams *[]ffprobe.Stream, prefix string) []widget.TreeNodeID {
	nodes := []widget.TreeNodeID{}
	for i := 0; i < len(*streams); i++ {
		str := prefix + " " + fmt.Sprint(i)
		nodes = append(nodes, str)
	}
	return nodes
}

func generateStreamDetailNodes(id string, prefix string) []widget.TreeNodeID {
	i := strings.Split(id, " ")[1]
	return []widget.TreeNodeID{
		prefix + fmt.Sprint(i) + "_codec_name",
		prefix + fmt.Sprint(i) + "_width",
		prefix + fmt.Sprint(i) + "_height",
		prefix + fmt.Sprint(i) + "_bit_rate",
	}
}

func getStreamByID(file *fileinfo.FileInfo, id widget.TreeNodeID) *ffprobe.Stream {
	i := int(id[1] - '0')
	switch id[0] {
	case 'v':
		return &(*file.VideoStreams)[i]
	case 'a':
		return &(*file.AudioStreams)[i]
	case 's':
		return &(*file.SubtitleStreams)[i]
	}
	return nil
}

func setStreamDetailText(label *widget.Label, stream *ffprobe.Stream, id widget.TreeNodeID) {
	switch true {
	case strings.HasSuffix(id, "codec_name"):
		label.SetText("Codec name: " + stream.CodecName)
	case strings.HasSuffix(id, "width"):
		label.SetText("Width: " + fmt.Sprint(stream.Width))
	case strings.HasSuffix(id, "height"):
		label.SetText("Height: " + fmt.Sprint(stream.Height))
	case strings.HasSuffix(id, "bit_rate"):
		label.SetText("Bit rate: " + fmt.Sprint(stream.BitRate))
	case strings.HasSuffix(id, "channels"):
		label.SetText("Channels: " + fmt.Sprint(stream.Channels))
	case strings.HasSuffix(id, "title"):
		label.SetText("Title: " + stream.Tags.Title)
	}
}
