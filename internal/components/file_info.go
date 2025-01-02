package components

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/components/info"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/fileinfo"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type FileInfoComponent struct {
	widget.BaseWidget
	file  fileinfo.FileInfo
	tree  *widget.Tree
	title string
}

func NewFileInfoComponent(file fileinfo.FileInfo) *FileInfoComponent {

	c := &FileInfoComponent{
		file:  file,
		title: file.GetFilename(),
	}

	c.tree = widget.NewTree(c.childUIDs, c.isBranch, c.create, c.update)
	c.ExtendBaseWidget(c)
	return c
}

func (f *FileInfoComponent) CreateRenderer() fyne.WidgetRenderer {
	b := container.NewBorder(
		container.NewVBox(
			widget.NewLabel(lang.L("folder")+f.file.GetFolder()),
			widget.NewLabel("Filename: "+f.file.GetFilename()),
		),
		nil,
		nil,
		nil,
		f.tree,
	)
	return widget.NewSimpleRenderer(b)
}

func (f *FileInfoComponent) childUIDs(id widget.TreeNodeID) []widget.TreeNodeID {
	file := f.file
	switch id {
	case "":
		return []widget.TreeNodeID{"streams", "format"}
	case "streams":
		return []widget.TreeNodeID{"Video", "Audio", "Subtitle"}
	case "Video":
		return generateStreamNodes(file.GetVideoStreams(), "Video")
	case "Audio":
		return generateStreamNodes(file.GetAudioStreams(), "Audio")
	case "Subtitle":
		return generateStreamNodes(file.GetSubtitleStreams(), "Subtitle")
	case "format":
		return []widget.TreeNodeID{"format_name", "duration", "size"}
	}

	switch {
	case strings.HasPrefix(id, "Video"):
		return info.NewVideoInfo(id).GetNodes()
	case strings.HasPrefix(id, "Audio"):
		return info.NewAudioInfo(id).GetNodes()
	case strings.HasPrefix(id, "Subtitle"):
		return info.NewSubtitleInfo(id).GetNodes()
	}

	return []string{}
}

func (f *FileInfoComponent) isBranch(id widget.TreeNodeID) bool {
	file := f.file
	switch id {
	case "", "format", "streams", "Video", "Audio", "Subtitle":
		return true
	}

	for i := 0; i < len(file.GetVideoStreams()); i++ {
		if id == "Video "+fmt.Sprint(i) {
			return true
		}
	}

	for i := 0; i < len(file.GetAudioStreams()); i++ {
		if id == "Audio "+fmt.Sprint(i) {
			return true
		}
	}

	for i := 0; i < len(file.GetSubtitleStreams()); i++ {
		if id == "Subtitle "+fmt.Sprint(i) {
			return true
		}
	}

	return false
}

func (f *FileInfoComponent) create(b bool) fyne.CanvasObject {
	return widget.NewLabel("template")
}

func (f *FileInfoComponent) update(id widget.TreeNodeID, isBranch bool, co fyne.CanvasObject) {
	file := f.file
	if isBranch {
		co.(*widget.Label).SetText(id)
		return
	}

	switch id {
	case "format_name":
		co.(*widget.Label).SetText("Format name: " + file.GetInfo().Format.FormatName)
	case "duration":
		co.(*widget.Label).SetText("Duration: " + fmt.Sprint(file.GetInfo().Format.Duration()))
	case "size":
		co.(*widget.Label).SetText(lang.L("size") + file.GetInfo().Format.Size)
	default:
		stream := getStreamByID(file, id)
		switch id[0] {
		case 'v':
			info.NewVideoInfo(id).From(co.(*widget.Label), &stream, id)
		case 'a':
			info.NewAudioInfo(id).From(co.(*widget.Label), &stream, id)
		case 's':
			info.NewSubtitleInfo(id).From(co.(*widget.Label), &stream, id)
		}
	}
}

func generateStreamNodes(streams []ffprobe.Stream, prefix string) []widget.TreeNodeID {
	nodes := make([]widget.TreeNodeID, len(streams))
	for i := 0; i < len(streams); i++ {
		str := prefix + " " + fmt.Sprint(i)
		nodes[i] = str
	}
	return nodes
}

func getStreamByID(file fileinfo.FileInfo, id widget.TreeNodeID) ffprobe.Stream {
	i := int(id[1] - '0')
	switch id[0] {
	case 'v':
		return file.GetVideoStreams()[i]
	case 'a':
		return file.GetAudioStreams()[i]
	case 's':
		return file.GetSubtitleStreams()[i]
	default:
		return ffprobe.Stream{}
	}
}
