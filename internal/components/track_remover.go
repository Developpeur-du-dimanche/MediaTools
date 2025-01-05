package components

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
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
	processText       *widget.Label
	procesPopup       *widget.PopUp
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
		processText:       widget.NewLabel(""),
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

	f.procesPopup = widget.NewPopUp(f.processText, (*f.window).Canvas())

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
	default:
		return false
	}
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
		var streams []*ffprobe.Stream

		for _, track := range *trackRemover {
			tr := track.(*TrackFilter)

			var streamsToCheck []*ffprobe.Stream
			switch tr.StreamType {
			case Video:
				streamsToCheck = file.GetVideoStreams()
			case Audio:
				streamsToCheck = file.GetAudioStreams()
			case Subtitle:
				streamsToCheck = file.GetSubtitleStreams()
			}

			for _, stream := range streamsToCheck {
				switch tr.condition {
				case "ignore":
					continue
				case "equals":
					if tr.filter(stream, TrackRemoverEquals, tr.value) {
						streams = append(streams, stream)
					}
				case "not equals":
					if tr.filter(stream, TrackRemoverNotEquals, tr.value) {
						streams = append(streams, stream)
					}
				case "contains":
					if tr.filter(stream, TrackRemoverContains, tr.value) {
						streams = append(streams, stream)
					}
				case "not contains":
					if tr.filter(stream, TrackRemoverNotContains, tr.value) {
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

	if len(streamToRemove) == 0 {
		dialog.ShowInformation("No track found", "No track found", *f.window)
		return
	}

	var streamToDelete []string
	str := "(%s) Stream %d from %s, title: %s, codec: %s"

	for _, s := range streamToRemove {
		for _, stream := range s.stream {
			streamToDelete = append(streamToDelete, fmt.Sprintf(str, stream.CodecType, stream.Index, s.file.FileName, stream.Tags.Title, stream.CodecName))
		}
	}

	// show dialog to confirm deletion
	dialog := dialog.NewConfirm("Confirm deletion", "Are you sure you want to delete the selected tracks?\n"+strings.Join(streamToDelete, "\n"), func(b bool) {
		if b {
			for _, s := range streamToRemove {
				f.RemoveStream(s)
			}
		}
	}, *f.window)

	dialog.Show()
}

func (f *TrackRemoverComponent) RemoveStream(s StreamFileRemove) {

	var streamsIndex []string
	for _, stream := range s.stream {
		var streamType string
		switch stream.CodecType {
		case "video":
			streamType = "v"
		case "audio":
			streamType = "a"
		case "subtitle":
			streamType = "s"
		}
		streamsIndex = append(streamsIndex, fmt.Sprintf("-0:%s:%d", streamType, stream.Index))
	}

	args := []string{
		"-progress",
		"pipe:1",
		"-i", s.file.Directory + "/" + s.file.FileName,
		"-map", "0",
	}

	for _, index := range streamsIndex {
		args = append(args, "-map", index)
	}

	args = append(args, "-c", "copy", "-y", s.file.Directory+"/"+strings.TrimSuffix(s.file.FileName, filepath.Ext(s.file.FileName))+"-new."+s.file.Extension)

	progress := helper.NewCmd()

	totalSize, err := helper.CalculateTotalSize([]string{s.file.Directory + "/" + s.file.FileName})

	if err != nil {
		dialog.ShowError(err, *f.window)
		return
	}

	command, err := progress.RunCommandWithProgress(args, totalSize)
	if err != nil {
		dialog.ShowError(err, *f.window)
		return
	}

	f.procesPopup.Show()
	defer f.procesPopup.Hide()

	for p := range progress.CProgress {
		if p.PercentComplete == 100 {
			f.processText.SetText(lang.L("merging_files") + " 100%")
			break
		}
		f.processText.SetText(fmt.Sprintf("%s %.1f%% (%s: %.1fx, Bitrate: %.1f kbits/s)",
			fmt.Sprintf("Supression de piste %s", s.file.FileName),
			p.PercentComplete,
			lang.L("speed"),
			p.Speed,
			p.Bitrate))
	}

	err = command.Wait()

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
