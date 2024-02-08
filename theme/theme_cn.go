package theme

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
	"screenplay/res"
)

type ThemeCn struct {
}

var _ fyne.Theme = (*ThemeCn)(nil)

func (*ThemeCn) Font(s fyne.TextStyle) fyne.Resource {

	return &fyne.StaticResource{
		StaticName:    "SmileySans-Oblique.ttf",
		StaticContent: res.ChineseTtf,
	}

}

func (*ThemeCn) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {

	if n == theme.ColorNameDisabled {
		if v == theme.VariantDark {
			return color.Gray16{0xaaaf}
		} else {
			return color.Gray16{0x666f}
		}
	}
	return theme.DefaultTheme().Color(n, v)

}

func (*ThemeCn) Icon(n fyne.ThemeIconName) fyne.Resource {

	return theme.DefaultTheme().Icon(n)

}

func (*ThemeCn) Size(n fyne.ThemeSizeName) float32 {

	//if n == theme.SizeNamePadding {
	//	return 0
	//}
	return theme.DefaultTheme().Size(n)

}
