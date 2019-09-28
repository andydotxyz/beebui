package beebui

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

type beebTheme struct{}

func (beebTheme) BackgroundColor() color.Color {
	return color.RGBA{0x10, 0x10, 0x10, 0xff}
}

func (beebTheme) ButtonColor() color.Color {
	return color.White
}

func (beebTheme) DisabledButtonColor() color.Color {
	return color.White
}

func (beebTheme) IconColor() color.Color {
	return color.White
}

func (beebTheme) DisabledIconColor() color.Color {
	return color.Black
}

func (beebTheme) HyperlinkColor() color.Color {
	return color.White
}

func (beebTheme) TextColor() color.Color {
	return color.White
}

func (beebTheme) DisabledTextColor() color.Color {
	return color.White
}

func (beebTheme) HoverColor() color.Color {
	return color.Black
}

func (beebTheme) PlaceHolderColor() color.Color {
	return color.White
}

func (beebTheme) PrimaryColor() color.Color {
	return color.White
}

func (beebTheme) FocusColor() color.Color {
	return color.White
}

func (beebTheme) ScrollBarColor() color.Color {
	return color.Black
}

func (beebTheme) ShadowColor() color.Color {
	return color.Black
}

func (beebTheme) TextSize() int {
	return 18
}

func (beebTheme) TextFont() fyne.Resource {
	return theme.DefaultTextFont()
}

func (beebTheme) TextBoldFont() fyne.Resource {
	return theme.DefaultTextBoldFont()
}

func (beebTheme) TextItalicFont() fyne.Resource {
	return theme.DefaultTextItalicFont()
}

func (beebTheme) TextBoldItalicFont() fyne.Resource {
	return theme.DefaultTextBoldItalicFont()
}

func (beebTheme) TextMonospaceFont() fyne.Resource {
	return font
}

func (beebTheme) Padding() int {
	return 0
}

func (beebTheme) IconInlineSize() int {
	return 10
}

func (beebTheme) ScrollBarSize() int {
	return 10
}

func (beebTheme) ScrollBarSmallSize() int {
	return 10
}
