package components

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/components/widgets"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/fileinfo"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/list"
)

type MergeFilesContract interface {
	MergeFiles(files []string, output string) error
}

type inputFiles struct {
	file     string
	position int
	enabled  bool
}

type MergeFiles struct {
	MergeFilesContract
	listFiles   *list.List[*fileinfo.FileInfo]
	inputFiles  []inputFiles
	outputFile  string
	window      fyne.Window
	procesPopup *widget.PopUp
	processText string
}

func NewMergeFilesComponent(window *fyne.Window, fileList *list.List[*fileinfo.FileInfo]) *MergeFiles {
	return &MergeFiles{
		listFiles: fileList,
		window:    *window,
	}
}

func (f *MergeFiles) Content() fyne.CanvasObject {

	mergeButton := widget.NewButton("Merge", func() {
		f.Merge()
	})
	outputFile := widget.NewEntry()
	outputFile.OnChanged = func(s string) {
		f.outputFile = s
	}

	listFiles := widget.NewList(
		func() int {
			return len(f.inputFiles)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewCheck("", nil), widgets.NewNumericalEntry(), widget.NewLabel(""))
		},
		func(i widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[0].(*widget.Check).Checked = f.inputFiles[i].enabled
			item.(*fyne.Container).Objects[1].(*widgets.NumericalEntry).SetText(fmt.Sprintf("%d", f.inputFiles[i].position))
			item.(*fyne.Container).Objects[2].(*widget.Label).SetText(f.inputFiles[i].file)

		},
	)

	refreshButton := widget.NewButton("Refresh", func() {
		f.inputFiles = []inputFiles{}
		for i, file := range f.listFiles.GetItems() {
			f.inputFiles = append(f.inputFiles, inputFiles{file.Path, i, true})
		}
		listFiles.Refresh()
	})

	content := container.NewBorder(
		container.NewHBox(
			refreshButton,
		),
		container.NewVBox(
			container.NewHBox(
				widget.NewLabel("Output file:"),
				outputFile,
			),
			mergeButton,
		),
		nil,
		nil,
		listFiles,
	)

	f.processText = ""

	processEntry := widget.NewEntry()
	processEntry.SetText(f.processText)

	f.procesPopup = widget.NewModalPopUp(processEntry, f.window.Canvas())

	return content
}

func (f *MergeFiles) Merge() {
	f.processText = "Merging files..."
	f.procesPopup.Show()

	defer f.procesPopup.Hide()

	var finalInputFiles []string
	// sort input files by position
	sort.SliceStable(f.inputFiles, func(i, j int) bool {
		return f.inputFiles[i].position < f.inputFiles[j].position
	})
	for _, file := range f.inputFiles {
		if file.enabled {
			finalInputFiles = append(finalInputFiles, file.file)
		}
	}

	if len(finalInputFiles) == 0 {
		return
	}

	err := f.MergeFiles(finalInputFiles, f.outputFile)
	if err != nil {
		fmt.Printf("Error: %s\n", err)

		widget.NewLabel(err.Error())
	}

}

func (f *MergeFiles) MergeFiles(files []string, output string) error {

	// create temp txt file
	txtFile, err := os.Create(output + ".txt")

	if err != nil {
		fmt.Printf("Error: %s\n", err)

		return err
	}

	for _, file := range files {
		_, err := txtFile.WriteString(fmt.Sprintf("file %s\n", file))
		if err != nil {
			fmt.Printf("Error: %s\n", err)

			return err
		}
	}

	// run ffmpeg command
	cmd := fmt.Sprintf("ffmpeg -f concat -safe 0 -i %s -c copy %s", txtFile.Name(), output)
	err = runCommand(cmd)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}

	txtFile.Close()

	// remove temp txt file
	err = os.Remove(txtFile.Name())
	if err != nil {
		fmt.Printf("Error: %s\n", err)

		return err
	}

	return nil
}

func runCommand(cmd string) error {
	path, err := exec.LookPath("FFMPEG")
	if err != nil {
		return err
	}

	if path == "" {
		return errors.New("ffmpeg not found")
	}

	command := exec.Command(path, strings.Split(cmd, " ")...)
	command.Stdout = os.Stdout
	err = command.Run()
	if err != nil {
		return err
	}

	return nil
}
