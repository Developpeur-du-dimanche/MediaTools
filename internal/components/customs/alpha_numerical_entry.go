package customs

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/widget"
)

type AlphaNumericalEntry struct {
	widget.Entry
}

func NewAlphaNumericalEntry() *AlphaNumericalEntry {
	entry := &AlphaNumericalEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *AlphaNumericalEntry) TypedRune(r rune) {
	e.Entry.TypedRune(r)
}

func (e *AlphaNumericalEntry) TypedShortcut(shortcut fyne.Shortcut) {
	_, ok := shortcut.(*fyne.ShortcutPaste)
	if !ok {
		e.Entry.TypedShortcut(shortcut)
		return
	}

	e.Entry.TypedShortcut(shortcut)
}

func (e *AlphaNumericalEntry) Keyboard() mobile.KeyboardType {
	return mobile.DefaultKeyboard
}

func (e *AlphaNumericalEntry) Enable() {
	e.Entry.Enable()
}

func (e *AlphaNumericalEntry) Disable() {
	e.Entry.Disable()
}

func (e *AlphaNumericalEntry) Hide() {
	e.Entry.Hide()
}

func (e *AlphaNumericalEntry) MinSize() fyne.Size {
	return e.Entry.MinSize()
}

func (e *AlphaNumericalEntry) Move(pos fyne.Position) {
	e.Entry.Move(pos)
}

func (e *AlphaNumericalEntry) Position() fyne.Position {
	return e.Entry.Position()
}

func (e *AlphaNumericalEntry) Refresh() {
	e.Entry.Refresh()
}

func (e *AlphaNumericalEntry) Resize(size fyne.Size) {
	e.Entry.Resize(size)
}

func (e *AlphaNumericalEntry) Show() {
	e.Entry.Show()
}

func (e *AlphaNumericalEntry) Size() fyne.Size {
	return e.Entry.Size()
}

func (e *AlphaNumericalEntry) Visible() bool {
	return e.Entry.Visible()
}
