package controller

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"screenplay/conf"
	"screenplay/cons"
	"screenplay/entity"
	"screenplay/global"
	"screenplay/utils"
	"sync"
)

var (
	vbox_time    = container.NewVBox()
	vbox_aside   = container.NewVBox()
	vbox_scene   = container.NewVBox()
	vbox_lens    = container.NewVBox()
	vbox_shot    = container.NewVBox()
	vbox_motion  = container.NewVBox()
	vbox_cmra    = container.NewVBox()
	vbox_action  = container.NewVBox()
	vbox_trans   = container.NewVBox()
	vbox_voice   = container.NewVBox()
	vbox_ext     = container.NewVBox()
	label_time   = widget.NewLabel("时间")
	label_aside  = widget.NewLabel("旁白")
	label_scene  = widget.NewLabel("环境")
	label_lens   = widget.NewLabel("焦段")
	label_shot   = widget.NewLabel("景别")
	label_motion = widget.NewLabel("运镜")
	label_cmera  = widget.NewLabel("镜头")
	label_action = widget.NewLabel("表演")
	label_trans  = widget.NewLabel("转场")
	label_voice  = widget.NewLabel("声音")
	label_ext    = widget.NewLabel("备注")
	col_time     = container.NewVBox(label_time, vbox_time)
	col_aside    = container.NewVBox(label_aside, vbox_aside)
	col_scene    = container.NewVBox(label_scene, vbox_scene)
	col_lens     = container.NewVBox(label_lens, vbox_lens)
	col_shot     = container.NewVBox(label_shot, vbox_shot)
	col_motion   = container.NewVBox(label_motion, vbox_motion)
	col_cmera    = container.NewVBox(label_cmera, vbox_cmra)
	col_action   = container.NewVBox(label_action, vbox_action)
	col_trans    = container.NewVBox(label_trans, vbox_trans)
	col_voice    = container.NewVBox(label_voice, vbox_voice)
	col_ext      = container.NewVBox(label_ext, vbox_ext)
	h_box        = container.NewHBox(
		col_aside,
		col_scene,
		col_lens,
		col_shot,
		col_motion,
		col_cmera,
		col_action,
		col_trans,
		col_voice,
		col_ext,
	)
	hscroll = container.NewHScroll(h_box)
	root    = container.NewHBox(
		col_time,
		hscroll,
	)
)

var (
	current_idx         = -1
	current_colume_type cons.ColumeType
	curr_total_sec      int
)

var (
	ctrl_input_sec *InputSecController
	lock           sync.Mutex
)

func CreateTable(ctrlInputSec *InputSecController) fyne.CanvasObject {

	ctrl_input_sec = ctrlInputSec
	ctrl_input_sec.OnChange = func(sec int) {
		if current_colume_type == cons.COLUME_TYPE_ASIDE {
			t := getDocTextByColumeType(current_colume_type, current_idx)
			if t.Txt != "" {
				return
			}
		}
		onInputSec(sec)
	}
	//if !global.App.Driver().Device().IsMobile() {
	//	hscroll.SetMinSize(fyne.NewSize(300, 0))
	//}
	return root
}

func ResizeTable() {

	//w := global.Win.Content().Size()
	////h := label_time.Size().Height + theme.Padding() + getHeightByCec(curr_total_sec)
	//h := label_time.Size().Height + theme.Padding() + vbox_time.Size().Height
	////h := time_col.Size().Height
	//hscroll.Resize(fyne.NewSize(w.Width-col_time.Size().Width, h))
	////for i := range time_col.Objects {
	////	fmt.Printf("-----i:%+v h:%+v\n", i, time_col.Objects[i].Size())
	////	for j := range vbox_time.Objects {
	////		fmt.Printf("-----j:%+v h:%+v\n", j, vbox_time.Objects[j].Size())
	////	}
	////}
	////fmt.Printf("w %+v %+v %+v \n", curr_total_sec, hscroll.Size(), h)

	col_time.ResizeByChild()
	h_box.ResizeByChild()
	w := global.Win.Content().Size()
	hscroll.Resize(fyne.NewSize(w.Width-col_time.Size().Width, col_time.Size().Height))
}

func DelCell() {
	if current_idx < 0 || current_colume_type == 0 {
		return
	}
	vb := getVboxByColumeType(current_colume_type)
	if len(vb.Objects) < 2 {
		return
	}
	if current_colume_type == cons.COLUME_TYPE_CMERA || current_colume_type == cons.COLUME_TYPE_LENS ||
		current_colume_type == cons.COLUME_TYPE_SHOT || current_colume_type == cons.COLUME_TYPE_MOTION {
		getVboxByColumeType(cons.COLUME_TYPE_LENS).RemoveAt(current_idx)
		getVboxByColumeType(cons.COLUME_TYPE_SHOT).RemoveAt(current_idx)
		getVboxByColumeType(cons.COLUME_TYPE_MOTION).RemoveAt(current_idx)
		getVboxByColumeType(cons.COLUME_TYPE_CMERA).RemoveAt(current_idx)

		removed := make([]entity.Camera, len(global.Doc.Camera)-1)
		copy(removed, global.Doc.Camera[:current_idx])
		copy(removed[current_idx:], global.Doc.Camera[current_idx+1:])
		global.Doc.Camera = removed

	} else {
		vb.RemoveAt(current_idx)

		l, _ := getDocTextListByColumeType(current_colume_type)

		ls := *l
		removed := make([]entity.Text, len(ls)-1)
		copy(removed, ls[:current_idx])
		copy(removed[current_idx:], ls[current_idx+1:])
		*l = removed

	}
	current_idx = -1
	recheckTotalSec()
	global.DocChanged = true

}
func InsertCell(up bool) {
	if global.Doc.TotalSec() == 0 {
		global.Doc.Aside = []entity.Text{{Sec: 1}}
		global.Doc.Scene = []entity.Text{{Sec: 1}}
		global.Doc.Camera = []entity.Camera{{Text: entity.Text{Sec: 1}}}
		global.Doc.Action = []entity.Text{{Sec: 1}}
		global.Doc.Trans = []entity.Text{{Sec: 1}}
		global.Doc.Voice = []entity.Text{{Sec: 1}}
		global.Doc.Ext = []entity.Text{{Sec: 1}}
		RefreshTableFromDoc()
		//global.DocChanged = true
		return
	}
	if current_idx < 0 || current_colume_type == 0 {
		dialog.ShowInformation("提示", "请选中一个单元格，再操作插入", global.Win)
		return
	}

	idx := current_idx
	if !up {
		idx += 1
	} else {
		current_idx++
	}

	l, l2 := getDocTextListByColumeType(current_colume_type)

	if l != nil {
		ls := *l
		added := make([]entity.Text, len(ls)+1)
		copy(added, ls[:idx])
		added[idx] = entity.Text{Sec: 1}
		copy(added[idx+1:], ls[idx:])
		*l = added
	} else if l2 != nil {
		ls := *l2
		added := make([]entity.Camera, len(ls)+1)
		copy(added, ls[:idx])
		added[idx] = entity.Camera{Text: entity.Text{Sec: 1}}
		copy(added[idx+1:], ls[idx:])
		*l2 = added
	}

	if current_colume_type == cons.COLUME_TYPE_CMERA || current_colume_type == cons.COLUME_TYPE_LENS ||
		current_colume_type == cons.COLUME_TYPE_SHOT || current_colume_type == cons.COLUME_TYPE_MOTION {
		createCellAt(cons.COLUME_TYPE_CMERA, nil, &((*l2)[idx]), idx)
	} else {
		createCellAt(current_colume_type, &((*l)[idx]), nil, idx)
	}
	recheckTotalSec()
	global.DocChanged = true
}

func clean() {
	vbox_time.RemoveAll()
	vbox_aside.RemoveAll()
	vbox_scene.RemoveAll()
	vbox_lens.RemoveAll()
	vbox_shot.RemoveAll()
	vbox_motion.RemoveAll()
	vbox_cmra.RemoveAll()
	vbox_action.RemoveAll()
	vbox_trans.RemoveAll()
	vbox_voice.RemoveAll()
	vbox_ext.RemoveAll()
	current_idx = -1
	current_colume_type = 0
}

func refreshTimeColum(totalSec int) {
	for i := 1; i <= totalSec; i++ {
		lb := widget.NewLabel(utils.SecondsToClockFormat(i))
		lb.Alignment = fyne.TextAlignTrailing
		ct := container.NewCenter(lb)
		ct.Resize(fyne.NewSize(conf.Colume_time_width, conf.Row_heigh))
		vbox_time.Add(ct)
	}
}
func getHeightByCec(sec int) float32 {
	return conf.Row_heigh*float32(sec) + theme.Padding()*float32(sec-1)
}
func onInputSec(sec int) {

	if sec < 1 {
		return
	}

	if current_idx >= 0 && current_colume_type > 0 {
		tx := getDocTextByColumeType(current_colume_type, current_idx)
		if tx != nil {
			tx.Sec = sec
		}

		if current_colume_type == cons.COLUME_TYPE_CMERA || current_colume_type == cons.COLUME_TYPE_LENS ||
			current_colume_type == cons.COLUME_TYPE_SHOT || current_colume_type == cons.COLUME_TYPE_MOTION {
			oldEntry, _ := getWidgetByTypeAndIndex(cons.COLUME_TYPE_CMERA, current_idx)
			if oldEntry != nil {
				oldEntry.Resize(fyne.NewSize(oldEntry.Size().Width, getHeightByCec(sec)))
			}
			_, oldSelet := getWidgetByTypeAndIndex(cons.COLUME_TYPE_LENS, current_idx)
			if oldSelet != nil {
				oldSelet.Resize(fyne.NewSize(oldSelet.Size().Width, getHeightByCec(sec)))
			}
			_, oldSelet = getWidgetByTypeAndIndex(cons.COLUME_TYPE_SHOT, current_idx)
			if oldSelet != nil {
				oldSelet.Resize(fyne.NewSize(oldSelet.Size().Width, getHeightByCec(sec)))
			}
			_, oldSelet = getWidgetByTypeAndIndex(cons.COLUME_TYPE_MOTION, current_idx)
			if oldSelet != nil {
				oldSelet.Resize(fyne.NewSize(oldSelet.Size().Width, getHeightByCec(sec)))
			}
		} else {
			oldEntry, _ := getWidgetByTypeAndIndex(current_colume_type, current_idx)
			if oldEntry != nil {

				oldEntry.Resize(fyne.NewSize(oldEntry.Size().Width, getHeightByCec(sec)))
			}
		}

		global.DocChanged = true
	}

	recheckTotalSec()
}
func recheckTotalSec() {
	lock.Lock()
	defer lock.Unlock()
	totalSec := global.Doc.TotalSec()
	if totalSec != curr_total_sec {
		vbox_time.RemoveAll()
		//fmt.Printf("rm %+v %+v\n", vbox_time.Objects, vbox_time.Size())
		refreshTimeColum(totalSec)
		curr_total_sec = totalSec
		ResizeTable()
	}
}

func createCellAt(tp cons.ColumeType, text *entity.Text, camera *entity.Camera, addIdx int) {

	switch tp {
	case cons.COLUME_TYPE_ASIDE:
		et := widget.NewEntry()
		if global.Disable {
			et.Disable()
		}
		et.Text = text.Txt
		et.MultiLine = true
		et.Wrapping = fyne.TextWrapWord
		et.SetMinRowsVisible(1)
		et.OnChanged = func(s string) {
			idx := vbox_aside.IndexOf(et)
			if idx >= 0 {
				global.Doc.Aside[idx].Txt = et.Text
				sec := global.Doc.Aside[idx].ComputeAsideSec()
				//et.Resize(fyne.NewSize(w, getHeightByCec(sec)))
				//fmt.Printf("aside sec:%+v h:%+v\n", sec, conf.Row_heigh*float32(sec)+theme.Padding()*float32(sec-1))
				//fmt.Printf(" :%+v h:%+v\n", vbox_time.Size().Height, vbox_lens.Size().Height)
				ctrl_input_sec.SetValue(sec)
				onInputSec(sec)
				global.DocChanged = true
			}
		}
		et.OnFocusChanged = func(f bool) {
			if f {
				idx := vbox_aside.IndexOf(et)
				onFoucsChanged(idx, cons.COLUME_TYPE_ASIDE)
				ctrl_input_sec.SetValue(global.Doc.Aside[idx].Sec)
			}
		}
		sec := text.ComputeAsideSec()
		et.Resize(fyne.NewSize(100, getHeightByCec(sec)))

		if addIdx < 0 {
			vbox_aside.Add(et)
		} else {
			vbox_aside.AddAt(addIdx, et)
		}
	case cons.COLUME_TYPE_SCENE:
		et := widget.NewEntry()
		if global.Disable {
			et.Disable()
		}
		et.Text = text.Txt
		et.MultiLine = true
		et.Wrapping = fyne.TextWrapWord
		et.SetMinRowsVisible(1)
		et.OnChanged = func(s string) {
			idx := vbox_scene.IndexOf(et)
			if idx >= 0 {
				global.Doc.Scene[idx].Txt = et.Text
				global.DocChanged = true
			}
		}
		et.OnFocusChanged = func(f bool) {
			if f {
				idx := vbox_scene.IndexOf(et)
				onFoucsChanged(idx, cons.COLUME_TYPE_SCENE)
				ctrl_input_sec.SetValue(global.Doc.Scene[idx].Sec)
			}
		}
		sec := text.Sec
		et.Resize(fyne.NewSize(100, getHeightByCec(sec)))

		if addIdx < 0 {
			vbox_scene.Add(et)
		} else {
			vbox_scene.AddAt(addIdx, et)
		}
	case cons.COLUME_TYPE_LENS, cons.COLUME_TYPE_SHOT, cons.COLUME_TYPE_MOTION, cons.COLUME_TYPE_CMERA:
		var (
			sl_lens, sl_shot, sl_motion *widget.Select
		)
		sec := camera.Sec

		// 焦段
		sl_lens = widget.NewSelect(global.Options_lens.GetLabels(), func(s string) {
			idx := vbox_lens.IndexOf(sl_lens)
			if idx >= 0 {
				global.Doc.Camera[idx].Lens = global.Options_lens.GetValueByLabel(s)
				global.DocChanged = true
			}
		})
		if global.Disable {
			sl_lens.Disable()
		}
		sl_lens.SetSelected(global.Options_lens[camera.Lens])
		sl_lens.OnTapped = func(event *fyne.PointEvent) {
			idx := vbox_lens.IndexOf(sl_lens)
			onFoucsChanged(idx, cons.COLUME_TYPE_LENS)
			ctrl_input_sec.SetValue(global.Doc.Camera[idx].Sec)
		}
		sl_lens.Resize(fyne.NewSize(70, getHeightByCec(sec)))
		if addIdx < 0 {
			vbox_lens.Add(sl_lens)
		} else {
			vbox_lens.AddAt(addIdx, sl_lens)
		}

		// 景别
		sl_shot = widget.NewSelect(global.Options_shot.GetLabels(), func(s string) {
			idx := vbox_shot.IndexOf(sl_shot)
			if idx >= 0 {
				global.Doc.Camera[idx].Shot = global.Options_shot.GetValueByLabel(s)
				global.DocChanged = true
			}
		})
		if global.Disable {
			sl_shot.Disable()
		}
		sl_shot.SetSelected(global.Options_shot[camera.Shot])
		sl_shot.OnTapped = func(event *fyne.PointEvent) {
			idx := vbox_shot.IndexOf(sl_shot)
			onFoucsChanged(idx, cons.COLUME_TYPE_SHOT)
			ctrl_input_sec.SetValue(global.Doc.Camera[idx].Sec)
		}
		sl_shot.Resize(fyne.NewSize(70, getHeightByCec(sec)))
		if addIdx < 0 {
			vbox_shot.Add(sl_shot)
		} else {
			vbox_shot.AddAt(addIdx, sl_shot)
		}

		// 运镜
		sl_motion = widget.NewSelect(global.Options_cmotion.GetLabels(), func(s string) {
			idx := vbox_motion.IndexOf(sl_motion)
			if idx >= 0 {
				global.Doc.Camera[idx].Motion = global.Options_cmotion.GetValueByLabel(s)
				global.DocChanged = true
			}
		})
		if global.Disable {
			sl_motion.Disable()
		}
		sl_motion.SetSelected(global.Options_cmotion[camera.Motion])
		sl_motion.OnTapped = func(event *fyne.PointEvent) {
			idx := vbox_motion.IndexOf(sl_motion)
			onFoucsChanged(idx, cons.COLUME_TYPE_MOTION)
			ctrl_input_sec.SetValue(global.Doc.Camera[idx].Sec)
		}
		sl_motion.Resize(fyne.NewSize(70, getHeightByCec(sec)))
		if addIdx < 0 {
			vbox_motion.Add(sl_motion)
		} else {
			vbox_motion.AddAt(addIdx, sl_motion)
		}

		// 镜头描述
		et := widget.NewEntry()
		if global.Disable {
			et.Disable()
		}
		et.Text = camera.Txt
		et.MultiLine = true
		et.Wrapping = fyne.TextWrapWord
		et.SetMinRowsVisible(1)
		et.OnChanged = func(s string) {
			idx := vbox_cmra.IndexOf(et)
			if idx >= 0 {
				global.Doc.Camera[idx].Txt = et.Text
				global.DocChanged = true
			}
		}
		et.OnFocusChanged = func(f bool) {
			if f {
				idx := vbox_cmra.IndexOf(et)
				onFoucsChanged(idx, cons.COLUME_TYPE_CMERA)
				ctrl_input_sec.SetValue(global.Doc.Camera[idx].Sec)
			}
		}
		et.Resize(fyne.NewSize(100, getHeightByCec(sec)))
		if addIdx < 0 {
			vbox_cmra.Add(et)
		} else {
			vbox_cmra.AddAt(addIdx, et)
		}
	case cons.COLUME_TYPE_ACTION:
		et := widget.NewEntry()
		if global.Disable {
			et.Disable()
		}
		et.Text = text.Txt
		et.MultiLine = true
		et.Wrapping = fyne.TextWrapWord
		et.SetMinRowsVisible(1)
		et.OnChanged = func(s string) {
			idx := vbox_action.IndexOf(et)
			if idx >= 0 {
				global.Doc.Action[idx].Txt = et.Text
				global.DocChanged = true
			}
		}
		et.OnFocusChanged = func(f bool) {
			if f {
				idx := vbox_action.IndexOf(et)
				onFoucsChanged(idx, cons.COLUME_TYPE_ACTION)
				ctrl_input_sec.SetValue(global.Doc.Action[idx].Sec)
			}
		}
		sec := text.Sec
		et.Resize(fyne.NewSize(300, getHeightByCec(sec)))
		if addIdx < 0 {
			vbox_action.Add(et)
		} else {
			vbox_action.AddAt(addIdx, et)
		}
	case cons.COLUME_TYPE_TRANS:
		et := widget.NewEntry()
		if global.Disable {
			et.Disable()
		}
		et.Text = text.Txt
		et.MultiLine = true
		et.Wrapping = fyne.TextWrapWord
		et.SetMinRowsVisible(1)
		et.OnChanged = func(s string) {
			idx := vbox_trans.IndexOf(et)
			if idx >= 0 {
				global.Doc.Trans[idx].Txt = et.Text
				global.DocChanged = true
			}
		}
		et.OnFocusChanged = func(f bool) {
			if f {
				idx := vbox_trans.IndexOf(et)
				onFoucsChanged(idx, cons.COLUME_TYPE_TRANS)
				ctrl_input_sec.SetValue(global.Doc.Trans[idx].Sec)
			}
		}
		sec := text.Sec
		et.Resize(fyne.NewSize(100, getHeightByCec(sec)))
		if addIdx < 0 {
			vbox_trans.Add(et)
		} else {
			vbox_trans.AddAt(addIdx, et)
		}
	case cons.COLUME_TYPE_VOICE:
		et := widget.NewEntry()
		if global.Disable {
			et.Disable()
		}
		et.Text = text.Txt
		et.MultiLine = true
		et.Wrapping = fyne.TextWrapWord
		et.SetMinRowsVisible(1)
		et.OnChanged = func(s string) {
			idx := vbox_voice.IndexOf(et)
			if idx >= 0 {
				global.Doc.Voice[idx].Txt = et.Text
				global.DocChanged = true
			}
		}
		et.OnFocusChanged = func(f bool) {
			if f {
				idx := vbox_voice.IndexOf(et)
				onFoucsChanged(idx, cons.COLUME_TYPE_VOICE)
				ctrl_input_sec.SetValue(global.Doc.Voice[idx].Sec)
			}
		}
		sec := text.Sec
		et.Resize(fyne.NewSize(100, getHeightByCec(sec)))
		if addIdx < 0 {
			vbox_voice.Add(et)
		} else {
			vbox_voice.AddAt(addIdx, et)
		}
	case cons.COLUME_TYPE_EXT:
		et := widget.NewEntry()
		if global.Disable {
			et.Disable()
		}
		et.Text = text.Txt
		et.MultiLine = true
		et.Wrapping = fyne.TextWrapWord
		et.SetMinRowsVisible(1)
		et.OnChanged = func(s string) {
			idx := vbox_ext.IndexOf(et)
			if idx >= 0 {
				global.Doc.Ext[idx].Txt = et.Text
				global.DocChanged = true
			}
		}
		et.OnFocusChanged = func(f bool) {
			if f {
				idx := vbox_ext.IndexOf(et)
				onFoucsChanged(idx, cons.COLUME_TYPE_EXT)
				ctrl_input_sec.SetValue(global.Doc.Ext[idx].Sec)
			}
		}
		sec := text.Sec
		et.Resize(fyne.NewSize(100, getHeightByCec(sec)))
		if addIdx < 0 {
			vbox_ext.Add(et)
		} else {
			vbox_ext.AddAt(addIdx, et)
		}
	}
}
func RefreshTableFromDoc() {
	clean()
	totalSec := global.Doc.TotalSec()
	curr_total_sec = totalSec
	if totalSec == 0 {
		return
	}
	refreshTimeColum(totalSec)

	//旁白
	for i := range global.Doc.Aside {
		createCellAt(cons.COLUME_TYPE_ASIDE, &global.Doc.Aside[i], nil, -1)
	}
	//环境
	for i := range global.Doc.Scene {
		createCellAt(cons.COLUME_TYPE_SCENE, &global.Doc.Scene[i], nil, -1)
	}
	//镜头
	for i := range global.Doc.Camera {
		createCellAt(cons.COLUME_TYPE_CMERA, nil, &global.Doc.Camera[i], -1)
	}

	//表演
	for i := range global.Doc.Action {
		createCellAt(cons.COLUME_TYPE_ACTION, &global.Doc.Action[i], nil, -1)
	}
	//转场
	for i := range global.Doc.Trans {
		createCellAt(cons.COLUME_TYPE_TRANS, &global.Doc.Trans[i], nil, -1)
	}
	//音效
	for i := range global.Doc.Voice {
		createCellAt(cons.COLUME_TYPE_VOICE, &global.Doc.Voice[i], nil, -1)
	}
	//备注
	for i := range global.Doc.Ext {
		createCellAt(cons.COLUME_TYPE_EXT, &global.Doc.Ext[i], nil, -1)
	}
}

func onFoucsChanged(idx int, tp cons.ColumeType) {

	if current_idx == idx && current_colume_type == tp {
		return
	}
	if current_idx >= 0 && current_colume_type > 0 {
		oldEntry, oldSelet := getWidgetByTypeAndIndex(current_colume_type, current_idx)
		if oldEntry != nil {
			oldEntry.HighLight = false
			oldEntry.Refresh()
		} else if oldSelet != nil {
			oldSelet.HighLight = false
			oldSelet.Refresh()
		}
	}

	if idx >= 0 && tp > 0 {
		newEntry, newSelet := getWidgetByTypeAndIndex(tp, idx)
		if newEntry != nil {
			newEntry.HighLight = true
		} else if newSelet != nil {
			newSelet.HighLight = true
		}
	}

	current_idx = idx
	current_colume_type = tp
}

func ToggleEitable() {
	toggleEitableVbox(vbox_aside)
	toggleEitableVbox(vbox_scene)
	toggleEitableVbox(vbox_lens)
	toggleEitableVbox(vbox_shot)
	toggleEitableVbox(vbox_motion)
	toggleEitableVbox(vbox_cmra)
	toggleEitableVbox(vbox_action)
	toggleEitableVbox(vbox_trans)
	toggleEitableVbox(vbox_voice)
	toggleEitableVbox(vbox_ext)
}
func toggleEitableVbox(c *fyne.Container) {
	for i := range c.Objects {
		if o, ok := c.Objects[i].(fyne.Disableable); ok {
			if global.Disable {
				o.Disable()
			} else {
				o.Enable()
			}

		} else {
			fmt.Println("diasd bbbb")
		}
	}
	c.Refresh()
}

func getWidgetByTypeAndIndex(tp cons.ColumeType, idx int) (*widget.Entry, *widget.Select) {
	vbox := getVboxByColumeType(tp)
	obj := vbox.Objects[idx]
	switch tp {
	case cons.COLUME_TYPE_LENS, cons.COLUME_TYPE_SHOT, cons.COLUME_TYPE_MOTION:
		return nil, obj.(*widget.Select)
	default:
		return obj.(*widget.Entry), nil
	}
}

func getVboxByColumeType(tp cons.ColumeType) *fyne.Container {
	switch tp {
	case cons.COLUME_TYPE_ASIDE:
		return vbox_aside
	case cons.COLUME_TYPE_SCENE:
		return vbox_scene
	case cons.COLUME_TYPE_LENS:
		return vbox_lens
	case cons.COLUME_TYPE_SHOT:
		return vbox_shot
	case cons.COLUME_TYPE_MOTION:
		return vbox_motion
	case cons.COLUME_TYPE_CMERA:
		return vbox_cmra
	case cons.COLUME_TYPE_ACTION:
		return vbox_action
	case cons.COLUME_TYPE_TRANS:
		return vbox_trans
	case cons.COLUME_TYPE_VOICE:
		return vbox_voice
	case cons.COLUME_TYPE_EXT:
		return vbox_ext
	}
	return nil
}
func getDocTextByColumeType(tp cons.ColumeType, idx int) *entity.Text {
	switch tp {
	case cons.COLUME_TYPE_ASIDE:
		return &global.Doc.Aside[idx]
	case cons.COLUME_TYPE_SCENE:
		return &global.Doc.Scene[idx]
	case cons.COLUME_TYPE_LENS, cons.COLUME_TYPE_SHOT, cons.COLUME_TYPE_MOTION, cons.COLUME_TYPE_CMERA:
		return &global.Doc.Camera[idx].Text
	case cons.COLUME_TYPE_ACTION:
		return &global.Doc.Action[idx]
	case cons.COLUME_TYPE_TRANS:
		return &global.Doc.Trans[idx]
	case cons.COLUME_TYPE_VOICE:
		return &global.Doc.Voice[idx]
	case cons.COLUME_TYPE_EXT:
		return &global.Doc.Ext[idx]
	}
	return nil
}

func getDocTextListByColumeType(tp cons.ColumeType) (*[]entity.Text, *[]entity.Camera) {
	switch tp {
	case cons.COLUME_TYPE_ASIDE:
		return &global.Doc.Aside, nil
	case cons.COLUME_TYPE_SCENE:
		return &global.Doc.Scene, nil
	case cons.COLUME_TYPE_LENS, cons.COLUME_TYPE_SHOT, cons.COLUME_TYPE_MOTION, cons.COLUME_TYPE_CMERA:
		return nil, &global.Doc.Camera
	case cons.COLUME_TYPE_ACTION:
		return &global.Doc.Action, nil
	case cons.COLUME_TYPE_TRANS:
		return &global.Doc.Trans, nil
	case cons.COLUME_TYPE_VOICE:
		return &global.Doc.Voice, nil
	case cons.COLUME_TYPE_EXT:
		return &global.Doc.Ext, nil
	}
	return nil, nil
}
