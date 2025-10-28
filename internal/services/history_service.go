package services

import (
	"slices"

	"fyne.io/fyne/v2"
	"github.com/Developpeur-du-dimanche/MediaTools/pkg/logger"
)

const (
	// PreferenceKeyHistory is the key used to store folder history in preferences
	PreferenceKeyHistory = "last_scan_selector"
	// MaxHistoryItems is the maximum number of items to keep in history
	MaxHistoryItems = 10
)

// HistoryService manages folder scan history
type HistoryService struct {
	app fyne.App
}

// NewHistoryService creates a new history service instance
func NewHistoryService(app fyne.App) *HistoryService {
	return &HistoryService{
		app: app,
	}
}

// GetHistory returns the list of previously scanned folders
func (hs *HistoryService) GetHistory() []string {
	history := hs.app.Preferences().StringListWithFallback(PreferenceKeyHistory, []string{})
	logger.Debugf("Retrieved history: %d items", len(history))
	return history
}

// AddFolder adds a folder to the history
func (hs *HistoryService) AddFolder(path string) {
	history := hs.GetHistory()

	// Don't add if already exists
	if slices.Contains(history, path) {
		logger.Debugf("Folder already in history: %s", path)
		return
	}

	// Add to the beginning of the list
	history = append([]string{path}, history...)

	// Limit history size
	if len(history) > MaxHistoryItems {
		history = history[:MaxHistoryItems]
	}

	hs.app.Preferences().SetStringList(PreferenceKeyHistory, history)
	logger.Infof("Added folder to history: %s", path)
}

// ClearHistory removes all items from the history
func (hs *HistoryService) ClearHistory() {
	hs.app.Preferences().SetStringList(PreferenceKeyHistory, []string{})
	logger.Info("History cleared")
}

// RemoveFolder removes a specific folder from the history
func (hs *HistoryService) RemoveFolder(path string) {
	history := hs.GetHistory()

	// Find and remove the path
	for i, folder := range history {
		if folder == path {
			history = append(history[:i], history[i+1:]...)
			break
		}
	}

	hs.app.Preferences().SetStringList(PreferenceKeyHistory, history)
	logger.Infof("Removed folder from history: %s", path)
}
