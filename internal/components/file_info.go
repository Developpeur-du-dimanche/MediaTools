package components

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/fileinfo"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type FileInfoComponent struct {
	widget.BaseWidget
	file  *fileinfo.FileInfo
	tree  *widget.Tree
	title string
}

func NewFileInfoComponent(file *fileinfo.FileInfo) *FileInfoComponent {

	c := &FileInfoComponent{
		file:  file,
		title: file.Filename,
	}

	c.tree = widget.NewTree(c.childUIDs, c.isBranch, c.create, c.update)
	c.ExtendBaseWidget(c)
	return c
}

func (f *FileInfoComponent) CreateRenderer() fyne.WidgetRenderer {
	b := container.NewBorder(
		container.NewVBox(
			widget.NewLabel("Folder: "+f.file.Folder),
			widget.NewLabel("Filename: "+f.file.Filename),
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
}

func (f *FileInfoComponent) isBranch(id widget.TreeNodeID) bool {
	file := f.file
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