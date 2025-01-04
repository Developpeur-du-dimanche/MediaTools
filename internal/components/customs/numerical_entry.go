package customs

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/widget"
)

type NumericalEntry struct {
	widget.Entry
}

func NewNumericalEntry() *NumericalEntry {
	entry := &NumericalEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *NumericalEntry) TypedRune(r rune) {
	if (r >= '0' && r <= '9') || r == '.' || r == ',' {
		e.Entry.TypedRune(r)
	}
}

func (e *NumericalEntry) TypedShortcut(shortcut fyne.Shortcut) {
	paste, ok := shortcut.(*fyne.ShortcutPaste)
	if !ok {
		e.Entry.TypedShortcut(shortcut)
		return
	}

	content := paste.Clipboard.Content()
	if _, err := strconv.ParseFloat(content, 64); err == nil {
		e.Entry.TypedShortcut(shortcut)
	}
}

func (e *NumericalEntry) Keyboard() mobile.KeyboardType {
	return mobile.NumberKeyboard
}

func (e *NumericalEntry) Enable() {
	e.Entry.Enable()
}

func (e *NumericalEntry) Disable() {
	e.Entry.Disable()
}

func (e *NumericalEntry) Hide() {
	e.Entry.Hide()
}

func (e *NumericalEntry) MinSize() fyne.Size {
	return e.Entry.MinSize()
}

func (e *NumericalEntry) Move(pos fyne.Position) {
	e.Entry.Move(pos)
}

func (e *NumericalEntry) Position() fyne.Position {
	return e.Entry.Position()
}

func (e *NumericalEntry) Refresh() {
	e.Entry.Refresh()
}

func (e *NumericalEntry) Resize(size fyne.Size) {
	e.Entry.Resize(size)
}

func (e *NumericalEntry) Show() {
	e.Entry.Show()
}

func (e *NumericalEntry) Size() fyne.Size {
	return e.Entry.Size()
}

func (e *NumericalEntry) Visible() bool {
	return e.Entry.Visible()
}
