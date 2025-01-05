package components

import (
	"os"
	"os/exec"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/helper"
	jsonfilter "github.com/Developpeur-du-dimanche/MediaTools/pkg/filter"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/list"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type TrackRemoverComponent struct {
	container         *fyne.Container
	fileList          *list.List[*helper.FileMetadata]
	window            *fyne.Window
	removeTrackButton *widget.Button
}

type TrackRemoverCondition string

const (
	TrackRemoverEquals      TrackRemoverCondition = "equals"
	TrackRemoverContains    TrackRemoverCondition = "contains"
	TrackRemoverNotEquals   TrackRemoverCondition = "not equals"
	TrackRemoverNotContains TrackRemoverCondition = "not contains"
)

type trackFilter func(*ffprobe.Stream, StreamType, jsonfilter.FilterType, string) bool

func NewTrackRemoverComponent(window *fyne.Window, fileList *list.List[*helper.FileMetadata]) Component {
	return &TrackRemoverComponent{
		container:         container.NewVBox(),
		fileList:          fileList,
		window:            window,
		removeTrackButton: widget.NewButton("Remove", nil),
	}
}

func (f *TrackRemoverComponent) Content() fyne.CanvasObject {
	var objects []fyne.CanvasObject = []fyne.CanvasObject{}
	tree := widget.NewTree(
		func(id widget.TreeNodeID) []widget.TreeNodeID {
			switch id {
			case "":
				return []widget.TreeNodeID{"Video", "Audio", "Subtitles"}
			case "Video":
				return []widget.TreeNodeID{"v_id", "v_title", "v_lang"}
			case "Audio":
				return []widget.TreeNodeID{"a_id", "a_title", "a_lang"}
			case "Subtitles":
				return []widget.TreeNodeID{"s_id", "s_title", "s_lang"}
			default:
				return []string{}
			}
		},
		func(id widget.TreeNodeID) bool {
			switch id {
			case "":
				return true
			case "Video":
				return true
			case "Audio":
				return true
			case "Subtitles":
				return true
			default:
				return false
			}
		}, func(b bool) fyne.CanvasObject {
			if b {
				return widget.NewLabel("")
			}
			return NewTrackFilter("template")
		}, func(id widget.TreeNodeID, b bool, co fyne.CanvasObject) {
			if b {
				co.(*widget.Label).SetText(id)
				return
			}

			switch id {
			case "v_id":
				co.(*TrackFilter).SetText("ID")
				co.(*TrackFilter).StreamType = Video
				co.(*TrackFilter).filter = f.checkId
			case "v_title":
				co.(*TrackFilter).SetText("Title")
				co.(*TrackFilter).StreamType = Video
				co.(*TrackFilter).filter = f.checkTitle
			case "v_lang":
				co.(*TrackFilter).SetText("Language")
				co.(*TrackFilter).StreamType = Video
				co.(*TrackFilter).filter = f.checkLanguage
			case "a_id":
				co.(*TrackFilter).SetText("ID")
				co.(*TrackFilter).StreamType = Audio
				co.(*TrackFilter).filter = f.checkId
			case "a_title":
				co.(*TrackFilter).SetText("Title")
				co.(*TrackFilter).StreamType = Audio
				co.(*TrackFilter).filter = f.checkTitle
			case "a_lang":
				co.(*TrackFilter).SetText("Language")
				co.(*TrackFilter).StreamType = Audio
				co.(*TrackFilter).filter = f.checkLanguage
			case "s_id":
				co.(*TrackFilter).SetText("ID")
				co.(*TrackFilter).StreamType = Subtitle
				co.(*TrackFilter).filter = f.checkId
			case "s_title":
				co.(*TrackFilter).SetText("Title")
				co.(*TrackFilter).StreamType = Subtitle
				co.(*TrackFilter).filter = f.checkTitle
			case "s_lang":
				co.(*TrackFilter).SetText("Language")
				co.(*TrackFilter).StreamType = Subtitle
				co.(*TrackFilter).filter = f.checkLanguage
			}

			alreadyAppend := false
			for _, o := range objects {
				if co.(*TrackFilter).Equals(o.(*TrackFilter)) {
					alreadyAppend = true
					break
				}
			}

			if !alreadyAppend {
				objects = append(objects, co)
			}

			co.(*TrackFilter).conditionWidget.SetSelectedIndex(0)

		},
	)

	f.removeTrackButton.OnTapped = func() {
		f.RemoveTrack(&objects)
	}

	return container.NewBorder(nil, f.removeTrackButton, nil, nil, tree)
}

func (f *TrackRemoverComponent) checkId(stream *ffprobe.Stream, condition TrackRemoverCondition, id string) bool {
	// id to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return false
	}
	switch condition {
	case TrackRemoverEquals:
		return stream.Index == idInt
	case TrackRemoverNotEquals:
		return stream.Index != idInt
	}
	return false
}

func (f *TrackRemoverComponent) checkTitle(stream *ffprobe.Stream, condition TrackRemoverCondition, title string) bool {
	switch condition {
	case TrackRemoverEquals:
		return stream.Tags.Title == title
	case TrackRemoverContains:
		return strings.Contains(stream.Tags.Title, title)
	case TrackRemoverNotEquals:
		return stream.Tags.Title != title
	case TrackRemoverNotContains:
		return !strings.Contains(stream.Tags.Title, title)
	}
	return false
}

func (f *TrackRemoverComponent) checkLanguage(stream *ffprobe.Stream, condition TrackRemoverCondition, language string) bool {
	switch condition {
	case TrackRemoverEquals:
		return stream.Tags.Language == language
	case TrackRemoverContains:
		return strings.Contains(stream.Tags.Language, language)
	case TrackRemoverNotEquals:
		return stream.Tags.Language != language
	case TrackRemoverNotContains:
		return !strings.Contains(stream.Tags.Language, language)
	}
	return false
}

type StreamFileRemove struct {
	stream []*ffprobe.Stream
	file   *helper.FileMetadata
}

func (f *TrackRemoverComponent) RemoveTrack(trackRemover *[]fyne.CanvasObject) {
	streamToRemove := []StreamFileRemove{}
	for _, file := range f.fileList.GetItems() {
		streams := []*ffprobe.Stream{}
		for _, stream := range file.GetVideoStreams() {
			for _, track := range *trackRemover {
				trackRemover := track.(*TrackFilter)
				switch trackRemover.condition {
				case "ignore":
					continue
				case "equals":
					if trackRemover.filter(stream, TrackRemoverEquals, trackRemover.value) {
						streams = append(streams, stream)
					}
				case "not equals":
					if trackRemover.filter(stream, TrackRemoverNotEquals, trackRemover.value) {
						streams = append(streams, stream)
					}
				case "contains":
					if trackRemover.filter(stream, TrackRemoverContains, trackRemover.value) {
						streams = append(streams, stream)
					}
				case "not contains":
					if trackRemover.filter(stream, TrackRemoverNotContains, trackRemover.value) {
						streams = append(streams, stream)
					}

				}
			}
		}

		for _, stream := range file.GetAudioStreams() {
			for _, track := range *trackRemover {
				trackRemover := track.(*TrackFilter)
				switch trackRemover.condition {
				case "ignore":
					continue
				case "equals":
					if trackRemover.filter(stream, TrackRemoverEquals, trackRemover.value) {
						streams = append(streams, stream)
					}
				case "not equals":
					if trackRemover.filter(stream, TrackRemoverNotEquals, trackRemover.value) {
						streams = append(streams, stream)
					}
				case "contains":
					if trackRemover.filter(stream, TrackRemoverContains, trackRemover.value) {
						streams = append(streams, stream)
					}
				case "not contains":
					if trackRemover.filter(stream, TrackRemoverNotContains, trackRemover.value) {
						streams = append(streams, stream)
					}

				}
			}
		}

		for _, stream := range file.GetSubtitleStreams() {
			for _, track := range *trackRemover {
				trackRemover := track.(*TrackFilter)
				switch trackRemover.condition {
				case "ignore":
					continue
				case "equals":
					if trackRemover.filter(stream, TrackRemoverEquals, trackRemover.value) {
						streams = append(streams, stream)
					}
				case "not equals":
					if trackRemover.filter(stream, TrackRemoverNotEquals, trackRemover.value) {
						streams = append(streams, stream)
					}
				case "contains":
					if trackRemover.filter(stream, TrackRemoverContains, trackRemover.value) {
						streams = append(streams, stream)
					}
				case "not contains":
					if trackRemover.filter(stream, TrackRemoverNotContains, trackRemover.value) {
						streams = append(streams, stream)
					}

				}
			}
		}

		if len(streams) > 0 {
			streamToRemove = append(streamToRemove, StreamFileRemove{
				stream: streams,
				file:   file,
			})
		}

	}

	// show dialog to confirm deletion
	dialog := dialog.NewConfirm("Confirm deletion", "Are you sure you want to delete the selected tracks?", func(b bool) {
		if b {
			for _, s := range streamToRemove {
				f.RemoveStream(s)
			}
		}
	}, *f.window)

	dialog.Show()
}

func (f *TrackRemoverComponent) RemoveStream(s StreamFileRemove) {

	// for all files, create ffmpeg command to remove the selected streams
	for _, stream := range s.stream {

		args := []string{
			"-i", s.file.FileName,
			"-map", "0",
			"-map", "-1",
			"-c", "copy",
		}

		if stream.Index > 0 {
			args = append(args, "-map", stream.ID)
		}

		args = append(args, "-c", "copy", "-y", s.file.FileName+"-new"+s.file.Extension)

		// execute command
		exec := exec.Command("ffmpeg", args...)

		err := exec.Run()
		if err != nil {
			dialog.ShowError(err, *f.window)
		}

		// remove old file
		err = os.Remove(s.file.FileName)

		if err != nil {
			dialog.ShowError(err, *f.window)
		}

		// rename new file
		err = os.Rename(s.file.FileName+"-new"+s.file.Extension, s.file.FileName)

		if err != nil {
			dialog.ShowError(err, *f.window)
		}

	}

}
