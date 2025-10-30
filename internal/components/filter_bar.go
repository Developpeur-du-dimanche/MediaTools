package components

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Developpeur-du-dimanche/MediaTools/internal/filters"
)

// FilterConditionRow represents a single filter condition with dropdowns
type FilterConditionRow struct {
	container      *fyne.Container
	fieldSelect    *widget.Select
	operatorSelect *widget.Select
	valueEntry     *widget.Entry
	valueSelect    *widget.Select
	logicalOp      *widget.Select
	removeButton   *widget.Button
}

// FilterBar represents the visual filter builder component
type FilterBar struct {
	widget.BaseWidget

	app           fyne.App
	window        fyne.Window
	popupWindow   fyne.Window
	filterDialog  dialog.Dialog
	conditions    []*FilterConditionRow
	mainButton    *widget.Button
	badge         *widget.Label
	activeFilters int

	onFilterApply func(filterStr string)
	onFilterClear func()
}

// getFilterFieldConfigs returns all available filter configurations
func getFilterFieldConfigs() []filters.Filter {
	return filters.GetAllFilters()
}

// NewFilterBar creates a new visual filter bar component
func NewFilterBar(window fyne.Window, onApply func(string), onClear func()) *FilterBar {
	fb := &FilterBar{
		window:        window,
		onFilterApply: onApply,
		onFilterClear: onClear,
		conditions:    make([]*FilterConditionRow, 0),
		activeFilters: 0,
	}

	// Create badge for active filters count
	fb.badge = widget.NewLabel("")
	fb.badge.Hide()

	// Main button to open filter dialog
	fb.mainButton = widget.NewButtonWithIcon("Filters", theme.SearchIcon(), fb.showFilterDialog)

	fb.ExtendBaseWidget(fb)
	return fb
}

// CreateRenderer creates the renderer for the filter bar
func (fb *FilterBar) CreateRenderer() fyne.WidgetRenderer {
	badgeContainer := container.NewHBox(
		fb.mainButton,
		fb.badge,
	)
	return widget.NewSimpleRenderer(badgeContainer)
}

// showFilterDialog displays the filter configuration dialog
func (fb *FilterBar) showFilterDialog() {
	// Create dialog content
	content := fb.createDialogContent()

	// Create custom dialog with larger size to accommodate dropdowns
	fb.filterDialog = dialog.NewCustom("Configure Filters", "Close", content, fb.window)
	fb.filterDialog.Resize(fyne.NewSize(900, 600))
	fb.filterDialog.Show()
}

// createDialogContent creates the content for the filter dialog
func (fb *FilterBar) createDialogContent() fyne.CanvasObject {
	// Header with title and description
	header := container.NewVBox(
		widget.NewLabelWithStyle("Filter Configuration", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
	)

	// Conditions container (without scroll to avoid z-index issues with dropdowns)
	conditionsContainer := container.NewVBox()
	for _, row := range fb.conditions {
		conditionsContainer.Add(row.container)
	}

	// If no conditions, add a helpful message
	if len(fb.conditions) == 0 {
		conditionsContainer.Add(widget.NewLabelWithStyle(
			"No filters configured. Click 'Add Condition' to start.",
			fyne.TextAlignCenter,
			fyne.TextStyle{Italic: true},
		))
	}

	// Add condition button
	addButton := widget.NewButtonWithIcon("Add Condition", theme.ContentAddIcon(), func() {
		fb.addCondition()
		conditionsContainer.Objects = nil
		for _, row := range fb.conditions {
			conditionsContainer.Add(row.container)
		}
		conditionsContainer.Refresh()
	})

	// Action buttons
	applyButton := widget.NewButtonWithIcon("Apply Filters", theme.ConfirmIcon(), func() {
		fb.applyFilters()
		fb.filterDialog.Hide()
	})
	applyButton.Importance = widget.HighImportance

	clearButton := widget.NewButtonWithIcon("Clear All", theme.DeleteIcon(), func() {
		fb.clearFilters()
		conditionsContainer.Objects = nil
		conditionsContainer.Add(widget.NewLabelWithStyle(
			"No filters configured. Click 'Add Condition' to start.",
			fyne.TextAlignCenter,
			fyne.TextStyle{Italic: true},
		))
		conditionsContainer.Refresh()
	})
	clearButton.Importance = widget.DangerImportance

	cancelButton := widget.NewButton("Cancel", func() {
		fb.filterDialog.Hide()
	})

	buttonsRow := container.NewHBox(
		addButton,
		widget.NewLabel(""), // Spacer
		clearButton,
		cancelButton,
		applyButton,
	)

	// Wrap conditions in a padded container for better spacing
	conditionsWithPadding := container.NewPadded(conditionsContainer)

	// Main layout without nested scrolls
	mainContent := container.NewBorder(
		container.NewVBox(header, widget.NewLabel("")),       // Header with spacing
		container.NewVBox(widget.NewSeparator(), buttonsRow), // Footer with separator
		nil,
		nil,
		conditionsWithPadding, // Center content
	)

	return mainContent
}

// addCondition adds a new filter condition row
func (fb *FilterBar) addCondition() {
	row := &FilterConditionRow{}

	// Get all available filters
	allFilters := getFilterFieldConfigs()

	// Create display options for the select widget
	displayOptions := make([]string, len(allFilters))
	for i, config := range allFilters {
		displayOptions[i] = config.GetFieldConfig().DisplayName
	}

	// Field selector with display names
	row.fieldSelect = widget.NewSelect(displayOptions, func(selected string) {
		// Find actual field config from display name
		var fieldConfig filters.Filter
		for i := range allFilters {
			if allFilters[i].GetFieldConfig().DisplayName == selected {
				fieldConfig = allFilters[i]
				break
			}
		}
		if fieldConfig != nil {
			fb.updateOperatorsForField(row, fieldConfig)
			fb.updateValueInputForField(row, fieldConfig)
		}
	})
	row.fieldSelect.PlaceHolder = "Select field..."

	// Operator selector
	row.operatorSelect = widget.NewSelect([]string{}, nil)
	row.operatorSelect.PlaceHolder = "Operator..."

	// Value input (Entry by default)
	row.valueEntry = widget.NewEntry()
	row.valueEntry.PlaceHolder = "Value..."

	// Value selector (for predefined values)
	row.valueSelect = widget.NewSelect([]string{}, nil)
	row.valueSelect.PlaceHolder = "Select value..."
	row.valueSelect.Hide()

	// Logical operator (AND/OR) - only shown if not the first condition
	row.logicalOp = widget.NewSelect([]string{"AND", "OR"}, nil)
	row.logicalOp.Selected = "AND"

	// Remove button
	row.removeButton = widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		fb.removeCondition(row)
	})
	row.removeButton.Importance = widget.DangerImportance

	// Build the row container with better styling
	var rowContent *fyne.Container

	if len(fb.conditions) > 0 {
		rowContent = container.NewVBox(
			widget.NewSeparator(),
			container.NewBorder(
				nil, nil,
				container.NewHBox(
					widget.NewLabel("  "),
					row.logicalOp,
				),
				nil,
				widget.NewLabel(""),
			),
			container.NewBorder(
				nil, nil,
				widget.NewLabel("  Where"),
				row.removeButton,
				container.NewHBox(
					container.NewGridWithColumns(3,
						row.fieldSelect,
						row.operatorSelect,
						container.NewStack(row.valueEntry, row.valueSelect),
					),
				),
			),
		)
	} else {
		// First condition - no logical operator
		row.logicalOp.Hide()
		rowContent = container.NewBorder(
			nil, nil,
			widget.NewLabel("  Where"),
			row.removeButton,
			container.NewHBox(
				container.NewGridWithColumns(3,
					row.fieldSelect,
					row.operatorSelect,
					container.NewStack(row.valueEntry, row.valueSelect),
				),
			),
		)
	}

	row.container = rowContent
	fb.conditions = append(fb.conditions, row)
}

// removeCondition removes a filter condition row
func (fb *FilterBar) removeCondition(row *FilterConditionRow) {
	newConditions := make([]*FilterConditionRow, 0)
	for _, c := range fb.conditions {
		if c != row {
			newConditions = append(newConditions, c)
		}
	}
	fb.conditions = newConditions

	// Update the first row to remove logical operator
	if len(fb.conditions) > 0 {
		fb.conditions[0].logicalOp.Hide()
	}

	// Refresh dialog if it's open
	if fb.filterDialog != nil {
		// This will be handled by the refresh in the dialog
	}
}

// updateOperatorsForField updates available operators based on field type
func (fb *FilterBar) updateOperatorsForField(row *FilterConditionRow, fieldConfig filters.Filter) {
	operators := filters.OperatorsByType[fieldConfig.GetFieldConfig().Type]

	row.operatorSelect.Options = operators
	if len(operators) > 0 {
		row.operatorSelect.Selected = operators[0]
	}
	row.operatorSelect.Refresh()
}

// updateValueInputForField updates the value input based on field type
func (fb *FilterBar) updateValueInputForField(row *FilterConditionRow, fieldConfig filters.Filter) {
	// Reset visibility
	row.valueEntry.Hide()
	row.valueSelect.Hide()

	// If field has predefined values, show dropdown
	if len(fieldConfig.GetFieldConfig().PredefinedValues) > 0 {
		row.valueSelect.Options = fieldConfig.GetFieldConfig().PredefinedValues
		row.valueSelect.Show()
	} else {
		// Otherwise show text entry with placeholder
		row.valueEntry.PlaceHolder = fieldConfig.GetFieldConfig().Placeholder
		if row.valueEntry.PlaceHolder == "" {
			row.valueEntry.PlaceHolder = "Value..."
		}
		row.valueEntry.Show()
	}

	row.container.Refresh()
}

// buildFilterString builds the filter expression string from conditions
func (fb *FilterBar) buildFilterString() string {
	if len(fb.conditions) == 0 {
		return ""
	}

	parts := make([]string, 0)
	allFilters := getFilterFieldConfigs()

	for i, row := range fb.conditions {
		// Get field key from display name
		var fieldKey string
		for _, config := range allFilters {
			if config.GetFieldConfig().DisplayName == row.fieldSelect.Selected {
				fieldKey = config.GetFieldConfig().Key
				break
			}
		}

		operator := row.operatorSelect.Selected
		var value string

		if row.valueSelect.Visible() {
			value = row.valueSelect.Selected
		} else {
			value = row.valueEntry.Text
		}

		// Skip incomplete conditions
		if fieldKey == "" || operator == "" || value == "" {
			continue
		}

		condition := fmt.Sprintf("%s %s %s", fieldKey, operator, value)

		if i > 0 && len(parts) > 0 {
			logicalOp := row.logicalOp.Selected
			if logicalOp == "" {
				logicalOp = "AND"
			}
			parts = append(parts, logicalOp)
		}

		parts = append(parts, condition)
	}

	return strings.Join(parts, " ")
}

// applyFilters applies the current filter configuration
func (fb *FilterBar) applyFilters() {
	filterStr := fb.buildFilterString()

	// Count valid conditions
	validConditions := 0
	allFilters := getFilterFieldConfigs()
	for _, row := range fb.conditions {
		var fieldKey string
		for _, config := range allFilters {
			if config.GetFieldConfig().DisplayName == row.fieldSelect.Selected {
				fieldKey = config.GetFieldConfig().Key
				break
			}
		}

		operator := row.operatorSelect.Selected
		var value string
		if row.valueSelect.Visible() {
			value = row.valueSelect.Selected
		} else {
			value = row.valueEntry.Text
		}

		if fieldKey != "" && operator != "" && value != "" {
			validConditions++
		}
	}

	fb.activeFilters = validConditions
	fb.updateBadge()

	if fb.onFilterApply != nil {
		fb.onFilterApply(filterStr)
	}
}

// clearFilters removes all conditions
func (fb *FilterBar) clearFilters() {
	fb.conditions = make([]*FilterConditionRow, 0)
	fb.activeFilters = 0
	fb.updateBadge()

	if fb.onFilterClear != nil {
		fb.onFilterClear()
	}
}

// updateBadge updates the filter count badge
func (fb *FilterBar) updateBadge() {
	if fb.activeFilters > 0 {
		fb.badge.SetText(fmt.Sprintf("(%d active)", fb.activeFilters))
		fb.badge.Show()
	} else {
		fb.badge.Hide()
	}
	fb.badge.Refresh()
}

// GetFilterText returns the current filter expression as text
func (fb *FilterBar) GetFilterText() string {
	return fb.buildFilterString()
}
