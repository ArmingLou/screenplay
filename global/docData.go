package global

import (
	"fyne.io/fyne/v2"
	"screenplay/entity"
)

var (
	Doc        entity.Doc
	DocFileUri fyne.URI
	DocChanged bool

	Disable bool

	Win     fyne.Window
	App     fyne.App
	WinSize fyne.Size

	FixSize func()
)
