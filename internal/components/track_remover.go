package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/fileinfo"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/list"
)

type TrackRemoverComponent struct {
	choices           *[]*ConditionalWidget
	container         *fyne.Container
	fileList          *list.List[fileinfo.FileInfo]
	window            *fyne.Window
	removeTrackButton *widget.Button
}

func NewTrackRemoverComponent(window *fyne.Window, fileList *list.List[fileinfo.FileInfo]) *TrackRemoverComponent {
	return &TrackRemoverComponent{
		choices:           &[]*ConditionalWidget{},
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
				//co.(*widgets.TrackFilter).NumberWidget()
			case "v_title":
				co.(*TrackFilter).SetText("Title")
				//co.(*widgets.TrackFilter).SetFilter(globalfilter.NewVideoTitleFilter())
			case "v_lang":
				co.(*TrackFilter).SetText("Language")
			case "a_id":
				co.(*TrackFilter).SetText("ID")
			case "a_title":
				co.(*TrackFilter).SetText("Title")
			case "a_lang":
				co.(*TrackFilter).SetText("Language")
			case "s_id":
				co.(*TrackFilter).SetText("ID")
			case "s_title":
				co.(*TrackFilter).SetText("Title")
			case "s_lang":
				co.(*TrackFilter).SetText("Language")

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

			//co.(*fyne.Container).Objects[1].(*widget.Select).SetSelectedIndex(0)

		},
	)

	f.removeTrackButton.OnTapped = func() {
		f.RemoveTrack(&objects)
	}

	return container.NewBorder(nil, f.removeTrackButton, nil, nil, tree)
}

func (f *TrackRemoverComponent) RemoveTrack(trackRemover *[]fyne.CanvasObject) {

	/*for _, file := range f.fileList.GetItems() {
		for _, track := range *trackRemover {
			fmt.Println(track.(*widgets.TrackFilter).GetCondition())
		}
	}*/
}
