package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"log"
	"screenplay/conf"
	"screenplay/controller"
	"screenplay/global"
	theme2 "screenplay/theme"
	"screenplay/view"
)

func main() {
	global.App = app.New()
	global.App.Settings().SetTheme(&theme2.ThemeCn{}) //支持中文
	global.Win = global.App.NewWindow("Screenplay Editor")
	if global.App.Driver().Device().IsMobile() {
		//global.App.Settings().SetScale(0.5)
		conf.App_size_width = 400
		conf.App_size_heigh = 900
	} else {
		conf.App_size_width = 1220
		conf.App_size_heigh = 600
	}
	global.Win.Resize(fyne.NewSize(conf.App_size_width, conf.App_size_heigh))

	inputSecCtrl := controller.NewInputSecController()

	ipsec := view.CreatInputSec(inputSecCtrl)
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.FolderNewIcon(), func() {
			controller.OnNewDoc()
		}),
		widget.NewToolbarAction(theme.FolderOpenIcon(), func() {
			controller.OnOpenDoc()
		}),
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			controller.OnSaveDoc(nil)
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.MoveUpIcon(), func() {
			if global.Disable {
				return
			}
			log.Println("new row up")
			controller.InsertCell(true)
		}),
		widget.NewToolbarAction(theme.CancelIcon(), func() {
			if global.Disable {
				return
			}
			log.Println("delete row")
			controller.DelCell()
		}),
		widget.NewToolbarAction(theme.MoveDownIcon(), func() {
			if global.Disable {
				return
			}
			log.Println("new row down")
			controller.InsertCell(false)
		}),
		widget.NewToolbarAction(theme.VisibilityIcon(), func() {
			log.Println("lock")
			global.Disable = !global.Disable
			if global.Disable {
				ipsec.Disable()
			} else {
				ipsec.Enable()
			}
			controller.ToggleEitable()
			global.FixSize()
		}),
	)

	vs := container.NewVScroll(controller.CreateTable(inputSecCtrl))
	topbar := container.NewHScroll(container.NewHBox(toolbar, ipsec))
	content := container.NewVBox(
		topbar,
		vs,
	)
	global.Win.SetContent(content)

	global.FixSize = func() {
		//time.AfterFunc(time.Millisecond*500, func() {
		//pading := theme.Padding()
		//ss := fyne.NewSize(global.WinSize.Width-pading*2, global.WinSize.Height-pading*3)
		ss := global.Win.Content().Size()
		//ss := size
		topbar.Resize(fyne.NewSize(ss.Width, topbar.Size().Height))
		vs.Resize(fyne.NewSize(ss.Width, ss.Height-topbar.Size().Height))
		controller.ResizeTable()
		//})
	}
	global.Win.SetOnResize(func(size fyne.Size) {
		fmt.Printf("+++++++ resize callback%+v\n", size)
		fmt.Printf("+++++++ theme.Padding()%+v\n", theme.Padding())
		global.WinSize = size
		global.FixSize()
	})
	global.Win.SetCloseIntercept(func() {
		if global.DocChanged {
			//提示保存旧
			dialog := dialog.NewConfirm("提示", "还有未保存的文件，请先保存文件！", func(b bool) {
				if b {
					controller.OnSaveDoc(func(err error) {
						if err == nil {
							global.Win.Close()
						}
					})
				} else {
					global.Win.Close()
				}
			}, global.Win)
			dialog.SetConfirmText("保存")
			dialog.SetDismissText("不保存")
			dialog.Show()
		} else {
			global.Win.Close()
		}
	})

	//w.SetFullScreen(true)

	//global.FixSize()
	//
	global.Win.ShowAndRun()
}
