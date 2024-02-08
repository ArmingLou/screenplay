package controller

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"io"
	"log"
	"screenplay/entity"
	"screenplay/global"
	"time"
)

func OnNewDoc() {

	if global.DocChanged {
		//提示保存旧
		dialog := dialog.NewConfirm("提示", "还有未保存的文件，请先保存文件！", func(b bool) {
			if b {
				OnSaveDoc(func(err error) {
					if err == nil {
						newDoc()
					}
				})
			} else {
				newDoc()
			}
		}, global.Win)
		dialog.SetConfirmText("保存")
		dialog.SetDismissText("不保存")
		dialog.Show()
	} else {
		newDoc()
	}
}
func OnOpenDoc() {

	if global.DocChanged {
		//提示保存旧
		dialog := dialog.NewConfirm("提示", "还有未保存的文件，请先保存文件！", func(b bool) {
			if b {
				OnSaveDoc(func(err error) {
					if err == nil {
						openDoc()
					}
				})
			} else {
				openDoc()
			}
		}, global.Win)
		dialog.SetConfirmText("保存")
		dialog.SetDismissText("不保存")
		dialog.Show()
	} else {
		openDoc()
	}
}

func OnSaveDoc(callback func(err error)) {
	if global.DocChanged {
		if global.DocFileUri == nil {
			dialog := dialog.NewFileSave(func(closer fyne.URIWriteCloser, err error) {
				if err != nil {
					dialog.ShowError(err, global.Win)
					if callback != nil {
						callback(err)
					}
					return
				}
				if closer == nil {
					//没有选择任何文件
					if callback != nil {
						callback(fmt.Errorf("取消保存"))
					}
					return
				}

				global.DocFileUri = closer.URI()
				e := saveDoc(closer)
				if e != nil {
					dialog.ShowError(fmt.Errorf("%+v (%+v)", e, global.DocFileUri), global.Win)
					if callback != nil {
						callback(e)
					}
					return
				}
				if callback != nil {
					callback(nil)
				}
			}, global.Win)
			dialog.SetFilter(&storage.ExtensionFileFilter{
				Extensions: []string{".sp"},
			})
			dialog.SetFileName(fmt.Sprintf("%v.sp", time.Now().UnixMilli()))
			dialog.Show()
		} else {
			e := saveDoc(nil)
			if e != nil {
				dialog.ShowError(fmt.Errorf("%+v (%+v)", e, global.DocFileUri), global.Win)
				if callback != nil {
					callback(e)
				}
				return
			}
			if callback != nil {
				callback(nil)
			}

		}

	}
}

func saveDoc(closer fyne.URIWriteCloser) error {

	bys, e := json.Marshal(global.Doc)
	if e != nil {
		return e
	}

	if closer == nil {
		closer, e = storage.Writer(global.DocFileUri)
		if e != nil {
			return e
		}
	}

	n2, err := closer.Write(bys)
	if err != nil {
		fmt.Println(err)
		closer.Close()
		return err
	}
	fmt.Println(n2, "bytes written successfully")
	err = closer.Close()
	if err != nil {
		fmt.Println(err)
	}

	log.Println("Save document")
	global.DocChanged = false
	return nil
}

func newDoc() error {
	log.Println("New document")
	global.Doc = entity.Doc{}
	global.DocFileUri = nil
	global.DocChanged = false
	RefreshTableFromDoc()
	return nil
}

func openDoc() error {
	dialog := dialog.NewFileOpen(func(closer fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, global.Win)
			return
		}
		if closer == nil {
			//没有选择任何文件
			return
		}

		strAll := ""
		bys := make([]byte, 256)
		//reader := bufio.NewReader(closer)
		for {
			i, e := closer.Read(bys)
			if e == nil {
				strAll += string(bys[:i])
			} else if e == io.EOF {
				//strAll += string(bys[:])
				break
			} else {
				dialog.ShowError(e, global.Win)
				closer.Close()
				return
			}
		}
		bty := []byte(strAll)

		closer.Close()

		e := json.Unmarshal(bty, &global.Doc)
		if e != nil {
			dialog.ShowError(e, global.Win)
			return
		}
		log.Println("Open document")
		global.DocFileUri = closer.URI()
		global.DocChanged = false
		//fmt.Printf("---%+v, err:%+v, path:%+v \n all:\n%+v\n ,doc:\n%+v\n", closer, err, global.DocFilePath, strAll, global.Doc)
		RefreshTableFromDoc()
	}, global.Win)
	dialog.SetFilter(&storage.ExtensionFileFilter{
		Extensions: []string{".sp"},
	})
	dialog.Show()
	return nil
}
