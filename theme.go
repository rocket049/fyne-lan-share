package main

import (
	_ "embed"
	"image/color"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

//go:embed ttf/stfangso.ttf
var fontWqy []byte

type MyTheme struct {
	Res fyne.Resource
}

var theTheme *MyTheme
var once sync.Once

func GetMyTheme() *MyTheme {
	once.Do(func() {
		theTheme = &MyTheme{Res: fyne.NewStaticResource("stfangso.ttf", fontWqy)}
	})
	return theTheme
}

func (m *MyTheme) Font(s fyne.TextStyle) fyne.Resource {
	return m.Res
}

func (m *MyTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(n, v)
}

func (m *MyTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (m *MyTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
}
