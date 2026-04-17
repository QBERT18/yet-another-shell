package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type terminalTheme struct{}

func (t *terminalTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 30, G: 30, B: 30, A: 255} // #1e1e1e
	case theme.ColorNameForeground:
		return color.NRGBA{R: 204, G: 204, B: 204, A: 255} // #cccccc
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 45, G: 45, B: 45, A: 255} // #2d2d2d
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{R: 102, G: 102, B: 102, A: 255} // #666666
	case theme.ColorNameScrollBar:
		return color.NRGBA{R: 85, G: 85, B: 85, A: 255} // #555555
	case theme.ColorNameInputBorder:
		return color.NRGBA{R: 68, G: 68, B: 68, A: 255} // #444444
	default:
		return theme.DefaultTheme().Color(name, theme.VariantDark)
	}
}

func (t *terminalTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTextMonospaceFont()
}

func (t *terminalTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *terminalTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameText {
		return 14
	}
	return theme.DefaultTheme().Size(name)
}
