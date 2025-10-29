package mediatools

import (
	"context"
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
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
	mt.filterBar = components.NewFilterBar(mt.window, mt.onFilterApply, mt.onFilterClear)
	mt.cleanButton = widget.NewButtonWithIcon("Clean", theme.DeleteIcon(), mt.onCleanButtonClicked)
	mt.selectAllBtn = widget.NewButtonWithIcon("Select All", theme.CheckButtonCheckedIcon(), mt.onSelectAllClicked)
	mt.unselectAllBtn = widget.NewButtonWithIcon("Unselect All", theme.CheckButtonIcon(), mt.onUnselectAllClicked)

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
	)

	// Section de la liste principale des fichiers scannés
	mediaListHeader := widget.NewLabelWithStyle("Scanned Files", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

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

	resultsLabel := widget.NewLabel("No filter applied")
	resultsLabel.TextStyle = fyne.TextStyle{Bold: true}

	applyButton := widget.NewButtonWithIcon("Apply Filter", theme.SearchIcon(), func() {
		filterStr := mt.filterBar.GetFilterText()
		if filterStr == "" {
			resultsLabel.SetText("No filter applied - Showing all files")
			mt.filteredMediaItems = mt.allMediaItems
		} else {
			mt.onFilterApply(filterStr)
			resultsLabel.SetText(fmt.Sprintf("Filter: %s - %d results", filterStr, len(mt.filteredMediaItems)))
		}
		mt.filterResultsList.Refresh()
	})
	applyButton.Importance = widget.HighImportance

	clearButton := widget.NewButtonWithIcon("Clear Filter", theme.ContentClearIcon(), func() {
		mt.onFilterClear()
		resultsLabel.SetText("Filter cleared - Showing all files")
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

	return container.NewTabItem("Filter", content)
}

// createMergeTab crée l'onglet pour fusionner des vidéos
func (mt *MediaTools) createMergeTab() *container.TabItem {
	placeholder := widget.NewLabel("Select at least 2 files above, then click 'Start Merge' to begin.")

	startButton := widget.NewButtonWithIcon("Start Merge", theme.MediaPlayIcon(), func() {
		selected := mt.listView.GetSelectedItems()
		if len(selected) < 2 {
			placeholder.SetText("Please select at least 2 files above.")
			return
		}
		mt.mergeComponent = components.NewMergeVideosComponent(mt.window, selected, mt.ffmpegService)
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

	return container.NewTabItem("Merge Videos", content)
}

// createRemoveStreamsTab crée l'onglet pour supprimer des pistes
func (mt *MediaTools) createRemoveStreamsTab() *container.TabItem {
	placeholder := container.NewCenter(
		widget.NewLabel("Select at least 1 file above, then click 'Start Processing' to begin."),
	)

	startButton := widget.NewButtonWithIcon("Start Processing", theme.ContentCutIcon(), func() {
		selected := mt.listView.GetSelectedItems()
		if len(selected) == 0 {
			placeholder := container.NewCenter(
				widget.NewLabel("Please select at least 1 file above."),
			)
			mt.removeStreamsTab.Content = placeholder
			mt.operationTabs.Refresh()
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
		placeholder,
	)

	return container.NewTabItem("Remove/Keep Streams", content)
}

// createCheckVideosTab crée l'onglet pour vérifier l'intégrité des vidéos
func (mt *MediaTools) createCheckVideosTab() *container.TabItem {
	placeholder := container.NewCenter(
		widget.NewLabel("Select at least 1 file above, then click 'Start Checking' to verify video integrity."),
	)

	startButton := widget.NewButtonWithIcon("Start Checking", theme.MediaPlayIcon(), func() {
		selected := mt.listView.GetSelectedItems()
		if len(selected) == 0 {
			placeholder := container.NewCenter(
				widget.NewLabel("Please select at least 1 file above."),
			)
			mt.checkVideosTab.Content = placeholder
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
		placeholder,
	)

	return container.NewTabItem("Check Videos", content)
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
	// Traiter les fichiers au fur et à mesure du scan
	if progress.CurrentFile != "" && !progress.IsComplete {
		mt.processMediaFile(progress.CurrentFile)
	}
}

func (mt *MediaTools) onFilterApply(filterStr string) {
	logger.Infof("Applying filter: %s", filterStr)
	mt.currentFilter = filterStr

	// Appliquer le filtre
	filtered, err := mt.filterService.FilterMediaList(mt.allMediaItems, filterStr)
	if err != nil {
		logger.Errorf("Filter error: %v", err)
		// TODO: Afficher une erreur à l'utilisateur
		return
	}

	mt.filteredMediaItems = filtered

	// Mettre à jour la liste affichée
	mt.updateDisplayedItems()
	logger.Infof("Filter applied: %d/%d items match", len(filtered), len(mt.allMediaItems))
}

func (mt *MediaTools) onFilterClear() {
	logger.Info("Clearing filter")
	mt.currentFilter = ""
	mt.filteredMediaItems = mt.allMediaItems

	// Mettre à jour la liste affichée
	mt.updateDisplayedItems()
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
		_, err := mt.mediaService.ScanFolder(ctx, folderPath, progressChan)

		if err != nil && err != context.Canceled {
			logger.Errorf("Folder scan error: %v", err)
		}
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

	// Appliquer le filtre si présent
	if mt.currentFilter != "" {
		// Re-filtrer toute la liste
		filtered, err := mt.filterService.FilterMediaList(mt.allMediaItems, mt.currentFilter)
		if err == nil {
			mt.filteredMediaItems = filtered
		}
	} else {
		mt.filteredMediaItems = mt.allMediaItems
	}

	// Mettre à jour l'affichage
	mt.updateDisplayedItems()
}

// updateDisplayedItems met à jour la liste affichée avec les items filtrés
func (mt *MediaTools) updateDisplayedItems() {
	mt.listView.Clear()

	itemsToDisplay := mt.filteredMediaItems
	if mt.currentFilter == "" {
		itemsToDisplay = mt.allMediaItems
	}

	for _, item := range itemsToDisplay {
		mt.listView.AddItem(item)
	}

	mt.listView.Refresh()
}

func (mt *MediaTools) onSelectAllClicked() {
	mt.listView.SelectAll()
}

func (mt *MediaTools) onUnselectAllClicked() {
	mt.listView.UnselectAll()
}

// Run démarre l'application
func (mt *MediaTools) Run() {
	mt.window.ShowAndRun()
}
