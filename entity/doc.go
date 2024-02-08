package entity

import (
	"regexp"
	"screenplay/conf"
	"screenplay/cons"
	"unicode/utf8"
)

type Doc struct {
	Aside  []Text   `json:"aside"`  //旁白
	Scene  []Text   `json:"scene"`  //环境
	Camera []Camera `json:"camera"` //镜头，景别/运镜/焦段
	Action []Text   `json:"action"` //表演
	Trans  []Text   `json:"trans"`  //转场
	Voice  []Text   `json:"voice"`  //音乐/音效
	Ext    []Text   `json:"ext"`    //备注
}

type Text struct {
	Sec int    `json:"sec"` //时长，秒
	Txt string `json:"text"`
}

type Camera struct {
	Text
	Lens   cons.CameraLens   `json:"lens"`   //焦段
	Shot   cons.CameraShot   `json:"shot"`   //景别
	Motion cons.CameraMotion `json:"motion"` //运镜
}

func (d *Text) ComputeAsideSec() int {
	if d.Txt == "" {
		return d.Sec
	}
	re := regexp.MustCompile("[\u4e00-\u9fa5]+")
	ls := re.FindAllString(d.Txt, -1)
	l := 0
	for i := range ls {
		l += utf8.RuneCountInString(ls[i])
	}
	sh := l / conf.Words_per_sec
	md := l % conf.Words_per_sec
	if md > 0 {
		sh += 1
	}
	if sh == 0 {
		sh = 1
	}
	//fmt.Printf("len:%d \n res: %+v", sh, ls)
	d.Sec = sh
	return sh
}

func (d *Doc) TotalSec() int {
	res := 0
	tmp := 0
	for i := range d.Aside {
		tmp += d.Aside[i].ComputeAsideSec()
	}
	if tmp > res {
		res = tmp
	}

	tmp = 0
	for i := range d.Scene {
		tmp += d.Scene[i].Sec
	}
	if tmp > res {
		res = tmp
	}

	tmp = 0
	for i := range d.Camera {
		tmp += d.Camera[i].Sec
	}
	if tmp > res {
		res = tmp
	}

	tmp = 0
	for i := range d.Action {
		tmp += d.Action[i].Sec
	}
	if tmp > res {
		res = tmp
	}

	tmp = 0
	for i := range d.Trans {
		tmp += d.Trans[i].Sec
	}
	if tmp > res {
		res = tmp
	}
	tmp = 0
	for i := range d.Voice {
		tmp += d.Voice[i].Sec
	}
	if tmp > res {
		res = tmp
	}
	tmp = 0
	for i := range d.Ext {
		tmp += d.Ext[i].Sec
	}
	if tmp > res {
		res = tmp
	}

	return res
}
