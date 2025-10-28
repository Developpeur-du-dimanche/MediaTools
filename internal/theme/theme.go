//go:generate fyne bundle -o bundled.go -package theme ../../assets
package theme

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type mediaToolsTheme struct {
	fyne.Theme
}

func NewMediaToolsTheme() fyne.Theme {
	return &mediaToolsTheme{Theme: theme.DefaultTheme()}
}

func (m *mediaToolsTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameBackground {
		return color.RGBA{R: 0x22, G: 0x22, B: 0x22, A: 0xff}
	}
	return m.Theme.Color(name, variant)
}

func (t *mediaToolsTheme) Font(s fyne.TextStyle) fyne.Resource {
	if s.Symbol || s.Monospace {
		return t.Theme.Font(s)
	}

	if s.Bold {
		if s.Italic {
			return resourcePoppinsBoldItalicTtf
		} else {
			return resourcePoppinsBoldTtf
		}
	}
	if s.Italic {
		return resourcePoppinsItalicTtf
	}
	return resourcePoppinsRegularTtf
}

func (t *mediaToolsTheme) Size(name fyne.ThemeSizeName) float32 {
	// Ajustement des tailles pour un design compact
	if name == theme.SizeNameText {
		return 12
	}

	if name == theme.SizeNamePadding {
		return 4
	}

	return t.Theme.Size(name)
}
