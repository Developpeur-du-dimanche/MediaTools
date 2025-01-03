package components

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/helper"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/list"
	"gopkg.in/vansante/go-ffprobe.v2"
)

type TrackRemoverComponent struct {
	container         *fyne.Container
	fileList          *list.List[*helper.FileMetadata]
	window            *fyne.Window
	removeTrackButton *widget.Button
}

func NewTrackRemoverComponent(window *fyne.Window, fileList *list.List[*helper.FileMetadata]) Component {
	return &TrackRemoverComponent{
		container:         container.NewVBox(),
		fileList:          fileList,
		window:            window,
		removeTrackButton: widget.NewButton("Remove", nil),
	}
}

func (f *TrackRemoverComponent) Content() fyne.CanvasObject {
	/*var objects []fyne.CanvasObject = []fyne.CanvasObject{}
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
				// co.(*TrackFilter).NumberWidget()
				co.(*TrackFilter).SetFilter(track.NewVideoTrackRemover().Filter)
			case "v_title":
				co.(*TrackFilter).SetText("Title")
				co.(*TrackFilter).StreamType = Video
				co.(*TrackFilter).SetFilter(track.NewVideoTrackRemover().Filter)
			case "v_lang":
				co.(*TrackFilter).SetText("Language")
				co.(*TrackFilter).StreamType = Video
				co.(*TrackFilter).SetFilter(track.NewVideoTrackRemover().Filter)
			case "a_id":
				co.(*TrackFilter).SetText("ID")
				co.(*TrackFilter).StreamType = Audio
			case "a_title":
				co.(*TrackFilter).SetText("Title")
				co.(*TrackFilter).StreamType = Audio
			case "a_lang":
				co.(*TrackFilter).SetText("Language")
				co.(*TrackFilter).StreamType = Audio
			case "s_id":
				co.(*TrackFilter).SetText("ID")
				co.(*TrackFilter).StreamType = Subtitle
			case "s_title":
				co.(*TrackFilter).SetText("Title")
				co.(*TrackFilter).StreamType = Subtitle
			case "s_lang":
				co.(*TrackFilter).SetText("Language")
				co.(*TrackFilter).StreamType = Subtitle
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

	return container.NewBorder(nil, f.removeTrackButton, nil, nil, tree)*/
	return f.container
}

func (f *TrackRemoverComponent) RemoveTrack(trackRemover *[]fyne.CanvasObject) {
	_ = []ffprobe.Stream{}
	for _, file := range f.fileList.GetItems() {
		for _, track := range *trackRemover {
			_ = track.(*TrackFilter)

			/*switch trackFilter.StreamType {
			case Video:
				for _, stream := range file.GetVideoStreams() {
					if trackFilter.Filter(stream) {
						streamsToRemove = append(streamsToRemove, stream)
					}
				}
			case Audio:
				for _, stream := range file.GetAudioStreams() {
					if trackFilter.Filter(stream) {
						streamsToRemove = append(streamsToRemove, stream)
					}
				}
			case Subtitle:
				for _, stream := range file.GetSubtitleStreams() {
					if trackFilter.(stream) {
						streamsToRemove = append(streamsToRemove, stream)
					}
				}
			}*/
			fmt.Println(file.FileName)
		}
	}
}
