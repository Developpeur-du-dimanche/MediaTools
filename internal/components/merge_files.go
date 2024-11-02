package components

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
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
	processText *widget.Label
}

func NewMergeFilesComponent(window *fyne.Window, fileList *list.List[*fileinfo.FileInfo]) *MergeFiles {
	return &MergeFiles{
		listFiles:   fileList,
		window:      *window,
		processText: widget.NewLabel(""),
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

	f.procesPopup = widget.NewModalPopUp(f.processText, f.window.Canvas())

	return content
}

func (f *MergeFiles) Merge() {
	f.processText.SetText("Merging files... 0%")
	f.procesPopup.Show()

	defer f.procesPopup.Hide()

	var finalInputFiles []string
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
		return err
	}
	defer txtFile.Close()

	for _, file := range files {
		_, err := txtFile.WriteString(fmt.Sprintf("file '%s'\n", strings.ReplaceAll(file, "\\", "/")))
		if err != nil {
			return err
		}
	}

	// replace all '\' by '/'
	output = strings.ReplaceAll(output, "\\", "/")

	// run ffmpeg command with progress
	cmd := fmt.Sprintf("-progress pipe:1 -f concat -safe 0 -i %s -c copy %s",
		strings.ReplaceAll(txtFile.Name(), "\\", "/"),
		output)
	err = f.runCommandWithProgress(cmd)
	if err != nil {
		return err
	}

	// remove temp txt file
	err = os.Remove(txtFile.Name())
	if err != nil {
		return err
	}

	return nil
}

func (f *MergeFiles) runCommandWithProgress(cmd string) error {
	path, err := exec.LookPath("FFMPEG")
	if err != nil {
		return err
	}

	if path == "" {
		return errors.New("ffmpeg not found")
	}

	command := exec.Command(path, strings.Split(cmd, " ")...)

	// Créer un pipe pour la sortie standard
	stdout, err := command.StdoutPipe()
	if err != nil {
		return err
	}

	// Créer un pipe pour la sortie d'erreur
	stderr, err := command.StderrPipe()
	if err != nil {
		return err
	}

	// Démarrer la commande
	if err := command.Start(); err != nil {
		return err
	}

	// Scanner pour lire la progression
	scanner := bufio.NewScanner(stdout)

	// Variables pour suivre la progression
	var duration float64
	var currentTime float64

	// Lire la sortie d'erreur dans une goroutine séparée
	go func() {
		errScanner := bufio.NewScanner(stderr)
		for errScanner.Scan() {
			line := errScanner.Text()
			// Extraire la durée totale
			if strings.Contains(line, "Duration:") {
				durationStr := strings.Split(strings.Split(line, "Duration: ")[1], ",")[0]
				duration = parseFFmpegTime(durationStr)
			}
		}
	}()

	// Lire la progression
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "out_time_us=") {
			timeStr := strings.Split(line, "=")[1]
			currentTime = float64(parseInt64(timeStr)) / 1000000 // convertir microsecondes en secondes

			if duration > 0 {
				progress := (currentTime / duration) * 100
				// Mettre à jour l'interface utilisateur avec le pourcentage
				f.processText.SetText(fmt.Sprintf("Merging files... %.0f%%", progress))
			}
		}
	}

	// Attendre que la commande se termine
	return command.Wait()
}

// Fonction utilitaire pour parser le temps FFmpeg (format HH:MM:SS.ms)
func parseFFmpegTime(timeStr string) float64 {
	parts := strings.Split(timeStr, ":")
	if len(parts) != 3 {
		return 0
	}

	hours, _ := strconv.ParseFloat(parts[0], 64)
	minutes, _ := strconv.ParseFloat(parts[1], 64)
	seconds, _ := strconv.ParseFloat(parts[2], 64)

	return hours*3600 + minutes*60 + seconds
}

// Fonction utilitaire pour parser les entiers 64 bits
func parseInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}
