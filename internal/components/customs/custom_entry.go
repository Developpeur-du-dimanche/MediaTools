package customs

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/mobile"
)

type CustomEntry interface {
	TypedRune(rune)
	TypedShortcut(fyne.Shortcut)
	Keyboard() mobile.KeyboardType
	Hide()
	MinSize() fyne.Size
	Move(fyne.Position)
	Position() fyne.Position
	Refresh()
	Resize(fyne.Size)
	Show()
	Size() fyne.Size
	Visible() bool
	Enable()
	Disable()
}
