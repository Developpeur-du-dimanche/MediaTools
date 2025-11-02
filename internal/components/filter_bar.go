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
	addChildButton *widget.Button        // "➕" to add a child filter
	level          int                   // indentation level (0 = root, 1 = child)
	parentRow      *FilterConditionRow   // reference to parent row (nil if root)
	children       []*FilterConditionRow // list of child rows
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

	// Add child button - initially visible, hidden when row becomes a parent
	row.addChildButton = widget.NewButtonWithIcon("", theme.ContentAddIcon(), func() {
		fb.addChildCondition(row)
	})

	// Initialize hierarchy fields
	row.level = 0
	row.parentRow = nil
	row.children = make([]*FilterConditionRow, 0)

	// Build the row container with better styling
	var rowContent *fyne.Container

	// Calculate indentation based on level
	indentPixels := row.level * 50

	indentLabel := widget.NewLabel(strings.Repeat("    ", row.level))

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
				container.NewHBox(indentLabel, widget.NewLabel("Where")),
				container.NewHBox(row.addChildButton, row.removeButton),
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
			container.NewHBox(indentLabel, widget.NewLabel("Where")),
			container.NewHBox(row.addChildButton, row.removeButton),
			container.NewHBox(
				container.NewGridWithColumns(3,
					row.fieldSelect,
					row.operatorSelect,
					container.NewStack(row.valueEntry, row.valueSelect),
				),
			),
		)
	}

	// Prevent unused variable warning
	_ = indentPixels

	row.container = rowContent
	fb.conditions = append(fb.conditions, row)
}

// removeCondition removes a filter condition row and its children
func (fb *FilterBar) removeCondition(row *FilterConditionRow) {
	// Remove all children first
	for _, child := range row.children {
		fb.removeCondition(child)
	}

	// Remove from parent's children list if this row has a parent
	if row.parentRow != nil {
		newChildren := make([]*FilterConditionRow, 0)
		for _, child := range row.parentRow.children {
			if child != row {
				newChildren = append(newChildren, child)
			}
		}
		row.parentRow.children = newChildren

		// If parent has no more children, show the add child button again
		if len(row.parentRow.children) == 0 {
			row.parentRow.addChildButton.Show()
		}
	}

	// Remove from main conditions list
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
}

// addChildCondition adds a child filter to the given parent row
func (fb *FilterBar) addChildCondition(parentRow *FilterConditionRow) {
	// Only allow one level of children (no infinite nesting)
	if parentRow.level > 0 {
		return
	}

	// Create a new condition row
	childRow := &FilterConditionRow{}

	// Get all available filters
	allFilters := getFilterFieldConfigs()

	// Create display options for the select widget
	displayOptions := make([]string, len(allFilters))
	for i, config := range allFilters {
		displayOptions[i] = config.GetFieldConfig().DisplayName
	}

	// Field selector with display names
	childRow.fieldSelect = widget.NewSelect(displayOptions, func(selected string) {
		// Find actual field config from display name
		var fieldConfig filters.Filter
		for i := range allFilters {
			if allFilters[i].GetFieldConfig().DisplayName == selected {
				fieldConfig = allFilters[i]
				break
			}
		}
		if fieldConfig != nil {
			fb.updateOperatorsForField(childRow, fieldConfig)
			fb.updateValueInputForField(childRow, fieldConfig)
		}
	})
	childRow.fieldSelect.PlaceHolder = "Select field..."

	// Operator selector
	childRow.operatorSelect = widget.NewSelect([]string{}, nil)
	childRow.operatorSelect.PlaceHolder = "Operator..."

	// Value input (Entry by default)
	childRow.valueEntry = widget.NewEntry()
	childRow.valueEntry.PlaceHolder = "Value..."

	// Value selector (for predefined values)
	childRow.valueSelect = widget.NewSelect([]string{}, nil)
	childRow.valueSelect.PlaceHolder = "Select value..."
	childRow.valueSelect.Hide()

	// Logical operator (AND/OR)
	childRow.logicalOp = widget.NewSelect([]string{"AND", "OR"}, nil)
	childRow.logicalOp.Selected = "AND"

	// Remove button
	childRow.removeButton = widget.NewButtonWithIcon("", theme.DeleteIcon(), func() {
		fb.removeCondition(childRow)
	})
	childRow.removeButton.Importance = widget.DangerImportance

	// Add child button - child rows can also have children (siblings)
	childRow.addChildButton = widget.NewButton("➕", func() {
		fb.addSiblingCondition(childRow)
	})

	// Set hierarchy: child is level 1, parent is level 0
	childRow.level = parentRow.level + 1
	childRow.parentRow = parentRow
	childRow.children = make([]*FilterConditionRow, 0)

	// Calculate indentation
	indentLabel := widget.NewLabel(strings.Repeat("    ", childRow.level))

	// Build child row container with indentation
	childRowContent := container.NewVBox(
		widget.NewSeparator(),
		container.NewBorder(
			nil, nil,
			container.NewHBox(
				widget.NewLabel("  "),
				childRow.logicalOp,
			),
			nil,
			widget.NewLabel(""),
		),
		container.NewBorder(
			nil, nil,
			container.NewHBox(indentLabel, widget.NewLabel("Where")),
			container.NewHBox(childRow.addChildButton, childRow.removeButton),
			container.NewHBox(
				container.NewGridWithColumns(3,
					childRow.fieldSelect,
					childRow.operatorSelect,
					container.NewStack(childRow.valueEntry, childRow.valueSelect),
				),
			),
		),
	)

	childRow.container = childRowContent

	// Find parent index and insert child after parent
	parentIndex := -1
	for i, c := range fb.conditions {
		if c == parentRow {
			parentIndex = i
			break
		}
	}

	if parentIndex != -1 {
		// Insert child after parent (or after last existing child)
		insertIndex := parentIndex + 1 + len(parentRow.children)

		// Insert into conditions slice
		newConditions := make([]*FilterConditionRow, 0, len(fb.conditions)+1)
		newConditions = append(newConditions, fb.conditions[:insertIndex]...)
		newConditions = append(newConditions, childRow)
		newConditions = append(newConditions, fb.conditions[insertIndex:]...)
		fb.conditions = newConditions

		// Add to parent's children list
		parentRow.children = append(parentRow.children, childRow)

		// Hide parent's add child button once it has children
		parentRow.addChildButton.Hide()

		// Refresh dialog
		fb.refreshDialog()
	}
}

// addSiblingCondition adds a sibling to the given child row (another child of the same parent)
func (fb *FilterBar) addSiblingCondition(siblingRow *FilterConditionRow) {
	if siblingRow.parentRow != nil {
		fb.addChildCondition(siblingRow.parentRow)
	}
}

// refreshDialog refreshes the filter dialog content
func (fb *FilterBar) refreshDialog() {
	if fb.filterDialog != nil {
		// Rebuild the dialog content
		content := fb.createDialogContent()
		fb.filterDialog.Hide()
		fb.filterDialog = dialog.NewCustom("Configure Filters", "Close", content, fb.window)
		fb.filterDialog.Resize(fyne.NewSize(900, 600))
		fb.filterDialog.Show()
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

	// Process only root-level conditions (children will be processed within their parents)
	for i, row := range fb.conditions {
		// Skip child rows - they'll be processed with their parent
		if row.parentRow != nil {
			continue
		}

		// Add logical operator before this condition (if not first root element)
		if len(parts) > 0 {
			logicalOp := row.logicalOp.Selected
			if logicalOp == "" {
				logicalOp = "AND"
			}
			parts = append(parts, logicalOp)
		}

		// Build this condition
		conditionStr := fb.buildConditionString(row, allFilters)
		if conditionStr == "" {
			continue // Skip incomplete conditions
		}

		// If this row has children, wrap it and its children in parentheses
		if len(row.children) > 0 {
			groupParts := make([]string, 0)
			groupParts = append(groupParts, conditionStr)

			// Add all children
			for _, child := range row.children {
				childConditionStr := fb.buildConditionString(child, allFilters)
				if childConditionStr != "" {
					// Add logical operator before child
					logicalOp := child.logicalOp.Selected
					if logicalOp == "" {
						logicalOp = "AND"
					}
					groupParts = append(groupParts, logicalOp)
					groupParts = append(groupParts, childConditionStr)
				}
			}

			// Wrap in parentheses
			parts = append(parts, "("+strings.Join(groupParts, " ")+")")
		} else {
			// No children, just add the condition
			parts = append(parts, conditionStr)
		}

		// Avoid unused variable warning
		_ = i
	}

	return strings.Join(parts, " ")
}

// buildConditionString builds a single condition string from a row
func (fb *FilterBar) buildConditionString(row *FilterConditionRow, allFilters []filters.Filter) string {
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

	// Return empty string if incomplete
	if fieldKey == "" || operator == "" || value == "" {
		return ""
	}

	return fmt.Sprintf("%s %s %s", fieldKey, operator, value)
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
