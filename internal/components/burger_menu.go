package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type BurgerMenu struct {
	widget.BaseWidget
	top       fyne.CanvasObject
	bottom    fyne.CanvasObject
	left      fyne.CanvasObject
	right     fyne.CanvasObject
	content   fyne.CanvasObject
	expanded  bool
	OnRefresh func()
	menuPanel *fyne.Container
	window    fyne.Window
	renderer  *burgerMenuRenderer
}

func NewBurgerMenu(top fyne.CanvasObject, bottom fyne.CanvasObject,
	left fyne.CanvasObject, right fyne.CanvasObject, content fyne.CanvasObject, w fyne.Window, onRefresh func()) *BurgerMenu {
	menu := &BurgerMenu{
		top:       top,
		bottom:    bottom,
		left:      left,
		right:     right,
		content:   content,
		expanded:  false,
		window:    w,
		OnRefresh: onRefresh,
	}
	menu.ExtendBaseWidget(menu)
	return menu
}

// Refresh surcharge la méthode Refresh de BaseWidget pour propager le rafraîchissement
func (b *BurgerMenu) Refresh() {
	b.BaseWidget.Refresh()
	if b.expanded {
		b.refreshContent()
	}
}

func (b *BurgerMenu) refreshContent() {
	if b.content != nil {
		if widget, ok := b.content.(fyne.Widget); ok {
			widget.Refresh()
		}
	}
	b.menuPanel.Refresh()
	if b.OnRefresh != nil {
		b.OnRefresh()
	}
}

func (b *BurgerMenu) CreateRenderer() fyne.WidgetRenderer {
	line1 := canvas.NewLine(theme.Color(theme.ColorNameForeground))
	line2 := canvas.NewLine(theme.Color(theme.ColorNameForeground))
	line3 := canvas.NewLine(theme.Color(theme.ColorNameForeground))

	contentSize := b.window.Canvas().Size().Height

	if b.top != nil {
		contentSize -= b.top.Size().Height
	}
	if b.bottom != nil {
		contentSize -= b.bottom.Size().Height
	}

	b.content.Resize(fyne.NewSize(b.window.Canvas().Size().Width-150, contentSize))

	b.menuPanel = container.NewBorder(
		b.top,
		b.bottom,
		b.left,
		b.right,
		b.content,
	)
	b.menuPanel.Hide()

	background := canvas.NewRectangle(color.Color(color.RGBA{R: 0, G: 0, B: 0, A: 0}))

	renderer := &burgerMenuRenderer{
		menu:      b,
		lines:     []*canvas.Line{line1, line2, line3},
		menuPanel: container.NewStack(background, b.menuPanel),
		window:    b.window,
	}

	b.renderer = renderer
	return renderer
}

func (b *BurgerMenu) Tapped(_ *fyne.PointEvent) {
	b.expanded = !b.expanded
	if b.expanded {
		b.menuPanel.Show()
		// Force un rafraîchissement complet à l'ouverture
		b.refreshContent()
	} else {
		b.menuPanel.Hide()
	}
	b.Refresh()
}

func (b *BurgerMenu) MouseIn(_ *desktop.MouseEvent) {
	if b.renderer != nil {
		b.renderer.Refresh()
	}
	b.Refresh()
}

func (b *BurgerMenu) MouseOut() {
	if b.renderer != nil {
		b.renderer.Refresh()
	}
	b.Refresh()
}

type burgerMenuRenderer struct {
	menu      *BurgerMenu
	lines     []*canvas.Line
	menuPanel *fyne.Container
	window    fyne.Window
}

func (r *burgerMenuRenderer) Layout(size fyne.Size) {
	lineWidth := float32(25)
	lineHeight := float32(2)
	spacing := float32(6)

	y := size.Height/2 - (lineHeight*3+spacing*2)/2
	for _, line := range r.lines {
		line.Position1 = fyne.NewPos(0, y)
		line.Position2 = fyne.NewPos(lineWidth, y)
		y += spacing + lineHeight
	}

	if r.menu.expanded {
		r.menuPanel.Resize(fyne.NewSize(r.menu.window.Canvas().Size().Width-200, 300))
		r.menuPanel.Move(fyne.NewPos(0, size.Height))

		// Force le rafraîchissement du contenu lors du redimensionnement
		if widget, ok := r.menu.content.(fyne.Widget); ok {
			widget.Refresh()
		}
	}
}

func (r *burgerMenuRenderer) MinSize() fyne.Size {
	return fyne.NewSize(30, 30)
}

func (r *burgerMenuRenderer) Refresh() {
	for _, line := range r.lines {
		line.StrokeColor = theme.Color(theme.ColorNameForeground)
		line.Refresh()
	}

	// Rafraîchir également le contenu si le menu est ouvert
	if r.menu.expanded {
		if widget, ok := r.menu.content.(fyne.Widget); ok {
			widget.Refresh()
		}
	}
}

func (r *burgerMenuRenderer) Objects() []fyne.CanvasObject {
	objects := make([]fyne.CanvasObject, 0)
	for _, line := range r.lines {
		objects = append(objects, line)
	}
	objects = append(objects, r.menuPanel)
	return objects
}

func (r *burgerMenuRenderer) Destroy() {}
