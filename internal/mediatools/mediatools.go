package mediatools

import (
	"context"
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/components"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/services"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/utils"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/logger"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/medias"
	"github.com/kbinani/screenshot"
)

// MediaTools représente l'application principale
type MediaTools struct {
	app      fyne.App
	window   fyne.Window
	listView *components.ListView

	// Services
	mediaService   *services.MediaService
	historyService *services.HistoryService
	filterService  *services.FilterService
	ffmpegService  *services.FFmpegService

	// UI Components
	openFolder     *components.OpenFolder
	openFile       *components.OpenFile
	history        *components.LastScanSelector
	filterBar      *components.FilterBar
	cleanButton    *widget.Button
	selectAllBtn   *widget.Button
	unselectAllBtn *widget.Button
	settingsButton *widget.Button
	settingsDialog *components.SettingsDialog

	// Tabs for operations (below media list)
	operationTabs    *container.AppTabs
	filterTab        *container.TabItem
	mergeTab         *container.TabItem
	removeStreamsTab *container.TabItem
	checkVideosTab   *container.TabItem

	// Components for tabs
	filterResultsList      *widget.List
	mergeComponent         *components.MergeVideosComponent
	removeStreamsComponent *components.RemoveStreamsComponent
	checkVideosComponent   *components.CheckVideosComponent

	// Data
	allMediaItems      []*medias.FfprobeResult
	filteredMediaItems []*medias.FfprobeResult
	currentFilter      string
}

// NewMediaTools crée une nouvelle instance de l'application
func NewMediaTools(app fyne.App) *MediaTools {
	mt := &MediaTools{
		app:    app,
		window: app.NewWindow(lang.L("MediaTools")),
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
	// Initialiser les services
	mt.mediaService = services.NewMediaService(utils.GetValidExtensions(), 10*time.Second)
	mt.historyService = services.NewHistoryService(mt.app)
	mt.filterService = services.NewFilterService()
	mt.ffmpegService = services.NewFFmpegService()

	// Initialiser les données
	mt.allMediaItems = make([]*medias.FfprobeResult, 0)
	mt.filteredMediaItems = make([]*medias.FfprobeResult, 0)

	// Initialiser les composants UI
	mt.listView = components.NewListView(nil, mt.window, mt.ffmpegService)
	mt.history = components.NewLastScanSelector(mt.historyService, mt.onHistoryFolderSelected)
	mt.openFolder = components.NewOpenFolder(mt.window, mt.onFolderOpened, mt.onScanProgress)
	mt.openFile = components.NewOpenFile(mt.window, mt.onFileOpened)
	mt.filterBar = components.NewFilterBar(mt.window, nil, nil)
	mt.cleanButton = widget.NewButtonWithIcon(lang.L("Clean"), theme.DeleteIcon(), mt.onCleanButtonClicked)
	mt.selectAllBtn = widget.NewButtonWithIcon(lang.L("SelectAll"), theme.CheckButtonCheckedIcon(), mt.onSelectAllClicked)
	mt.unselectAllBtn = widget.NewButtonWithIcon(lang.L("UnselectAll"), theme.CheckButtonIcon(), mt.onUnselectAllClicked)
	mt.settingsButton = widget.NewButtonWithIcon(lang.L("Settings"), theme.SettingsIcon(), mt.onSettingsClicked)
	mt.settingsDialog = components.NewSettingsDialog(mt.app, mt.window)

	// Initialiser les composants pour les onglets (seront créés à la demande)
	mt.filterResultsList = nil
	mt.mergeComponent = nil
	mt.removeStreamsComponent = nil
	mt.checkVideosComponent = nil
}

// setupLayout configure la disposition des éléments dans la fenêtre
func (mt *MediaTools) setupLayout() {
	// Barre d'outils du haut
	toolBar := container.NewHBox(
		mt.openFile,
		mt.openFolder,
		mt.cleanButton,
		mt.history,
		widget.NewSeparator(),
		mt.selectAllBtn,
		mt.unselectAllBtn,
		widget.NewSeparator(),
		mt.settingsButton,
	)

	// Section de la liste principale des fichiers scannés
	mediaListHeader := widget.NewLabelWithStyle(lang.L("ScannedFiles"), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	mediaSection := container.NewBorder(
		container.NewVBox(
			toolBar,
			widget.NewSeparator(),
			mediaListHeader,
		),
		nil,
		nil,
		nil,
		mt.listView,
	)

	// Créer les onglets d'opérations
	mt.filterTab = mt.createFilterTab()
	mt.mergeTab = mt.createMergeTab()
	mt.removeStreamsTab = mt.createRemoveStreamsTab()
	mt.checkVideosTab = mt.createCheckVideosTab()

	// Onglets d'opérations en dessous
	mt.operationTabs = container.NewAppTabs(
		mt.filterTab,
		mt.mergeTab,
		mt.removeStreamsTab,
		mt.checkVideosTab,
	)

	backgroud := canvas.NewRectangle(color.RGBA{
		R: 34,
		G: 34,
		B: 34,
		A: 255,
	})

	// Layout principal : médias en haut (50%), opérations en bas (50%)
	mainContent := container.NewVSplit(
		mediaSection,
		container.NewStack(backgroud, mt.operationTabs),
	)
	mainContent.SetOffset(0.5) // 50/50 split

	mt.window.SetContent(mainContent)
}

// createFilterTab crée l'onglet pour filtrer et afficher les résultats
func (mt *MediaTools) createFilterTab() *container.TabItem {
	// Initialiser la liste des résultats filtrés
	mt.filterResultsList = widget.NewList(
		func() int {
			return len(mt.filteredMediaItems)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			if id < len(mt.filteredMediaItems) {
				item := mt.filteredMediaItems[id]
				label.SetText(fmt.Sprintf("%d. %s", id+1, item.Format.Filename))
			}
		},
	)

	resultsLabel := widget.NewLabel(lang.L("NoFilterApplied"))
	resultsLabel.TextStyle = fyne.TextStyle{Bold: true}

	applyButton := widget.NewButtonWithIcon(lang.L("ApplyFilter"), theme.SearchIcon(), func() {
		filterStr := mt.filterBar.GetFilterText()
		if filterStr == "" {
			resultsLabel.SetText(lang.L("NoFilterAppliedShowingAll"))
			mt.filteredMediaItems = mt.allMediaItems
		} else {
			// Apply filter without affecting the main list
			filtered, err := mt.filterService.FilterMediaList(mt.allMediaItems, filterStr)
			if err != nil {
				logger.Errorf("Filter error: %v", err)
				resultsLabel.SetText(lang.L("FilterError", map[string]any{"Error": err.Error()}))
				return
			}
			mt.filteredMediaItems = filtered
			resultsLabel.SetText(lang.L("FilterResults", map[string]any{
				"Filter": filterStr,
				"Count":  len(mt.filteredMediaItems),
			}))
			logger.Infof("Filter applied: %d/%d items match", len(filtered), len(mt.allMediaItems))
		}
		mt.filterResultsList.Refresh()
	})
	applyButton.Importance = widget.HighImportance

	clearButton := widget.NewButtonWithIcon(lang.L("ClearFilter"), theme.ContentClearIcon(), func() {
		mt.filteredMediaItems = mt.allMediaItems
		resultsLabel.SetText(lang.L("FilterCleared"))
		mt.filterResultsList.Refresh()
	})

	// Header avec les contrôles
	header := container.NewVBox(
		mt.filterBar,
		container.NewHBox(applyButton, clearButton),
		widget.NewSeparator(),
		resultsLabel,
	)

	// La liste prend tout l'espace disponible avec Border
	content := container.NewBorder(
		header,
		nil,
		nil,
		nil,
		mt.filterResultsList,
	)

	return container.NewTabItem(lang.L("Filter"), content)
}

// createMergeTab crée l'onglet pour fusionner des vidéos
func (mt *MediaTools) createMergeTab() *container.TabItem {
	placeholder := widget.NewLabel(lang.L("SelectAtLeast2Files"))

	startButton := widget.NewButtonWithIcon(lang.L("StartMerge"), theme.MediaPlayIcon(), func() {

		selected := mt.listView.GetSelectedItems()
		if len(selected) < 2 {
			placeholder.SetText(lang.L("PleaseSelectAtLeast2Files"))
			return
		}

		mt.mergeComponent = components.NewMergeVideosComponent(mt.window, mt.ffmpegService, func() []*medias.FfprobeResult {
			selected := mt.listView.GetSelectedItems()
			if len(selected) < 2 {
				placeholder.SetText(lang.L("PleaseSelectAtLeast2Files"))
				return []*medias.FfprobeResult{}
			}
			return selected
		})
		mt.mergeTab.Content = mt.mergeComponent
		mt.operationTabs.Refresh()
	})
	startButton.Importance = widget.HighImportance

	content := container.NewBorder(
		nil,
		container.NewCenter(
			container.NewHBox(startButton),
		),
		nil,
		nil,
		container.NewCenter(placeholder),
	)

	return container.NewTabItem(lang.L("MergeVideos"), content)
}

// createRemoveStreamsTab crée l'onglet pour supprimer des pistes
func (mt *MediaTools) createRemoveStreamsTab() *container.TabItem {
	placeholder := widget.NewLabel(lang.L("SelectAtLeast1File"))

	startButton := widget.NewButtonWithIcon(lang.L("StartProcessing"), theme.ContentCutIcon(), func() {
		selected := mt.listView.GetSelectedItems()
		if len(selected) == 0 {
			placeholder.SetText(lang.L("PleaseSelectAtLeast1File"))
			return
		}
		mt.removeStreamsComponent = components.NewRemoveStreamsComponent(mt.window, selected, mt.ffmpegService)
		mt.removeStreamsTab.Content = mt.removeStreamsComponent
		mt.operationTabs.Refresh()
	})
	startButton.Importance = widget.HighImportance

	content := container.NewBorder(
		nil,
		container.NewCenter(
			container.NewHBox(startButton),
		),
		nil,
		nil,
		container.NewCenter(placeholder),
	)

	return container.NewTabItem(lang.L("RemoveKeepStreams"), content)
}

// createCheckVideosTab crée l'onglet pour vérifier l'intégrité des vidéos
func (mt *MediaTools) createCheckVideosTab() *container.TabItem {

	placeholder := widget.NewLabel(lang.L("SelectAtLeast1FileCheck"))

	startButton := widget.NewButtonWithIcon(lang.L("StartChecking"), theme.MediaPlayIcon(), func() {
		selected := mt.listView.GetSelectedItems()
		if len(selected) == 0 {
			placeholder.SetText(lang.L("PleaseSelectAtLeast1File"))
			mt.operationTabs.Refresh()
			return
		}
		mt.checkVideosComponent = components.NewCheckVideosComponent(mt.window, selected, mt.ffmpegService)
		mt.checkVideosTab.Content = mt.checkVideosComponent
		mt.operationTabs.Refresh()
	})
	startButton.Importance = widget.HighImportance

	content := container.NewBorder(
		nil,
		container.NewCenter(
			container.NewHBox(startButton),
		),
		nil,
		nil,
		container.NewCenter(placeholder),
	)

	return container.NewTabItem(lang.L("CheckVideos"), content)
}

func (mt *MediaTools) onHistoryFolderSelected(path string) {
	logger.Infof("History folder selected: %s", path)
	mt.scanFolder(path)
}

func (mt *MediaTools) onFolderOpened(path string) {
	// Ajouter à l'historique via le service
	mt.historyService.AddFolder(path)
	mt.history.Refresh()

	// Lancer le scan
	mt.scanFolder(path)
}

func (mt *MediaTools) onFileOpened(path string) {
	// Analyser un fichier unique
	mt.processMediaFile(path)
}

func (mt *MediaTools) onCleanButtonClicked() {
	mt.listView.Clear()
	mt.allMediaItems = make([]*medias.FfprobeResult, 0)
	mt.filteredMediaItems = make([]*medias.FfprobeResult, 0)
}

func (mt *MediaTools) onScanProgress(progress services.ScanProgress) {
	// Ne rien faire pendant le scan, le traitement se fera après
	// Cette fonction est juste pour la mise à jour de l'UI via UpdateProgress
}

// scanFolder lance le scan d'un dossier
func (mt *MediaTools) scanFolder(folderPath string) {
	// Créer un contexte annulable
	ctx, cancel := context.WithCancel(context.Background())
	mt.openFolder.SetCancelFunc(cancel)

	// Créer un canal de progression
	progressChan := make(chan services.ScanProgress, 10)

	// Lancer le scan en arrière-plan
	go func() {
		defer close(progressChan)
		mediaFiles, err := mt.mediaService.ScanFolder(ctx, folderPath, progressChan)

		if err != nil && err != context.Canceled {
			logger.Errorf("Folder scan error: %v", err)
			return
		}

		// Traiter tous les fichiers trouvés après le scan
		logger.Infof("Processing %d media files...", len(mediaFiles))
		for _, filePath := range mediaFiles {
			select {
			case <-ctx.Done():
				logger.Info("File processing cancelled")
				return
			default:
				mt.processMediaFile(filePath)
			}
		}
		logger.Info("All files processed successfully")
	}()

	// Traiter les mises à jour de progression
	go func() {
		for progress := range progressChan {
			mt.openFolder.UpdateProgress(progress)
		}
	}()
}

// processMediaFile analyse un fichier média et l'ajoute à la liste
func (mt *MediaTools) processMediaFile(path string) {
	ctx := context.Background()
	mediaInfo, err := mt.mediaService.GetMediaInfo(ctx, path)
	if err != nil {
		logger.Debugf("Skipping file %s: %v", path, err)
		return
	}

	// Ajouter à la liste complète
	mt.allMediaItems = append(mt.allMediaItems, mediaInfo)

	// Ajouter à la liste principale
	mt.listView.AddItem(mediaInfo)
}

func (mt *MediaTools) onSelectAllClicked() {
	mt.listView.SelectAll()
}

func (mt *MediaTools) onUnselectAllClicked() {
	mt.listView.UnselectAll()
}

func (mt *MediaTools) onSettingsClicked() {
	mt.settingsDialog.Show()
}

// Run démarre l'application
func (mt *MediaTools) Run() {
	mt.window.ShowAndRun()
}
