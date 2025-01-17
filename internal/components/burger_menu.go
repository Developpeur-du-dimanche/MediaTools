package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// BurgerMenu est un widget personnalisé qui implémente un menu burger
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
}

// NewBurgerMenu crée une nouvelle instance du menu burger
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

// CreateRenderer implémente l'interface Widget
func (b *BurgerMenu) CreateRenderer() fyne.WidgetRenderer {
	// Créer les trois lignes du burger
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

	// Container pour le menu déroulant
	b.menuPanel = container.NewBorder(
		b.top,
		b.bottom,
		b.left,
		b.right,
		b.content,
	)
	b.menuPanel.Hide()

	return &burgerMenuRenderer{
		menu:      b,
		lines:     []*canvas.Line{line1, line2, line3},
		menuPanel: b.menuPanel,
		window:    b.window,
	}
}

// Tapped gère les clics sur le menu
func (b *BurgerMenu) Tapped(_ *fyne.PointEvent) {
	b.expanded = !b.expanded
	if b.expanded {
		b.menuPanel.Show()
	} else {
		b.menuPanel.Hide()
	}
	b.Refresh()
}

// MouseIn gère l'entrée de la souris
func (b *BurgerMenu) MouseIn(_ *desktop.MouseEvent) {
	b.Refresh()
}

// MouseOut gère la sortie de la souris
func (b *BurgerMenu) MouseOut() {
	b.Refresh()
}

// burgerMenuRenderer implémente le renderer du widget
type burgerMenuRenderer struct {
	menu      *BurgerMenu
	lines     []*canvas.Line
	menuPanel *fyne.Container
	window    fyne.Window
}

// Layout gère le positionnement des éléments
func (r *burgerMenuRenderer) Layout(size fyne.Size) {
	// Dimensions des lignes du burger

	lineWidth := float32(25)
	lineHeight := float32(2)
	spacing := float32(6)

	// Positionner les lignes
	y := size.Height/2 - (lineHeight*3+spacing*2)/2
	for _, line := range r.lines {
		line.Position1 = fyne.NewPos(0, y)
		line.Position2 = fyne.NewPos(lineWidth, y)
		y += spacing + lineHeight
	}

	// Positionner le panel du menu
	if r.menu.expanded {
		r.menuPanel.Resize(fyne.NewSize(r.menu.window.Canvas().Size().Width-200, 300))
		r.menuPanel.Move(fyne.NewPos(0, size.Height))
	}
}

// MinSize retourne la taille minimale du widget
func (r *burgerMenuRenderer) MinSize() fyne.Size {
	return fyne.NewSize(30, 30)
}

// Refresh rafraîchit le rendu du widget
func (r *burgerMenuRenderer) Refresh() {
	for _, line := range r.lines {
		line.StrokeColor = theme.Color(theme.ColorNameForeground)
		line.Refresh()
	}
}

// Objects retourne tous les objets rendus
func (r *burgerMenuRenderer) Objects() []fyne.CanvasObject {
	objects := make([]fyne.CanvasObject, 0)
	for _, line := range r.lines {
		objects = append(objects, line)
	}
	objects = append(objects, r.menuPanel)
	return objects
}

// Destroy nettoie les ressources
func (r *burgerMenuRenderer) Destroy() {}
