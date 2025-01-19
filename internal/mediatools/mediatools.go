package mediatools

import (
	"context"
	"fmt"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/components"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"
	"github.com/kbinani/screenshot"
)

// MediaTools représente l'application principale
type MediaTools struct {
	app      fyne.App
	window   fyne.Window
	listView *components.ListView

	// UI Components
	burgerMenu  *components.BurgerMenu
	openFolder  *components.OpenFolder
	openFile    *components.OpenFile
	history     *components.LastScanSelector
	cleanButton *widget.Button

	// État de l'application
	scanState struct {
		sync.Mutex
		isScanning bool
		cancel     context.CancelFunc
	}
}

// NewMediaTools crée une nouvelle instance de l'application
func NewMediaTools(app fyne.App) *MediaTools {
	mt := &MediaTools{
		app:    app,
		window: app.NewWindow("MediaTools"),
	}

	mt.initWindow()
	mt.initComponents()
	mt.setupLayout()

	return mt
}

// initWindow initialise la fenêtre principale
func (mt *MediaTools) initWindow() {
	screen := screenshot.GetDisplayBounds(0)
	mt.window.Resize(fyne.NewSize(
		float32(screen.Dx()/2),
		float32(screen.Dy()/2),
	))
}

// initComponents initialise tous les composants de l'interface
func (mt *MediaTools) initComponents() {
	mt.listView = components.NewListView(nil, mt.window)

	mt.history = components.NewLastScanSelector(mt.onHistoryFolderSelected)
	mt.openFolder = components.NewOpenFolder(&mt.window, mt.onFolderOpened, mt.onNewFileDetected)
	mt.openFile = components.NewOpenFile(&mt.window, mt.onNewFileDetected)
	mt.cleanButton = widget.NewButtonWithIcon("Clean", theme.DeleteIcon(), mt.onCleanButtonClicked)

	// Configuration des callbacks
	mt.openFolder.OnScanTerminated = mt.onScanTerminated
	mt.openFile.OnScanTerminated = mt.onScanTerminated
}

// setupLayout configure la disposition des éléments dans la fenêtre
func (mt *MediaTools) setupLayout() {
	topBar := container.NewHBox(
		mt.openFolder,
		mt.openFile,
		mt.cleanButton,
		mt.history,
	)

	mt.burgerMenu = components.NewBurgerMenu(
		topBar, // top
		nil,    // bottom
		nil,    // left
		nil,    // right
		mt.listView,
		mt.window,
		mt.onBurgerMenuRefresh,
	)

	mt.window.SetContent(container.NewBorder(
		mt.burgerMenu,
		nil,
		nil,
		nil,
		nil,
	))
}

// Callbacks

func (mt *MediaTools) onHistoryFolderSelected(path string) {
	fmt.Printf("History folder selected: %s\n", path)
	// TODO: Implémenter le chargement du dossier historique
}

func (mt *MediaTools) onFolderOpened(path string) {
	mt.history.AddFolder(path)
	mt.listView.Refresh()
}

func (mt *MediaTools) onCleanButtonClicked() {
	mt.listView.Clear()
}

func (mt *MediaTools) onBurgerMenuRefresh() {
	mt.listView.Refresh()
}

func (mt *MediaTools) onScanTerminated() {
	mt.scanState.Lock()
	mt.scanState.isScanning = false
	if mt.scanState.cancel != nil {
		mt.scanState.cancel()
		mt.scanState.cancel = nil
	}
	mt.scanState.Unlock()

}

func (mt *MediaTools) onNewFileDetected(path string) {
	mediaInfo, err := mt.getMediaInfo(path)
	if err != nil {
		fmt.Printf("Error while getting media info: %s\n", err)
		return
	}

	mt.listView.AddItem(mediaInfo)
	mt.listView.Refresh()
}

// getMediaInfo récupère les informations médias d'un fichier
func (mt *MediaTools) getMediaInfo(path string) (*medias.FfprobeResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	mt.scanState.Lock()
	mt.scanState.isScanning = true
	mt.scanState.cancel = cancel
	mt.scanState.Unlock()

	defer func() {
		mt.scanState.Lock()
		if mt.scanState.cancel != nil {
			mt.scanState.isScanning = false
			mt.scanState.cancel = nil
		}
		mt.scanState.Unlock()
		cancel()
	}()

	ffprobe := medias.NewFfprobe(path,
		medias.FFPROBE_LOGLEVEL_FATAL,
		medias.PRINT_FORMAT_JSON,
		medias.SHOW_FORMAT,
		medias.SHOW_STREAMS,
		medias.EXPERIMENTAL,
	)

	return ffprobe.Probe(ctx)
}

// Run démarre l'application
func (mt *MediaTools) Run() {
	mt.window.ShowAndRun()
}

// Stop arrête proprement l'application
func (mt *MediaTools) Stop() {
	mt.scanState.Lock()
	if mt.scanState.cancel != nil {
		mt.scanState.cancel()
	}
	mt.scanState.Unlock()
}
