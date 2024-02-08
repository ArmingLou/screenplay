package global

import (
	"screenplay/cons"
	"screenplay/utils"
)

var (
	Options_lens    = make(utils.Options[cons.CameraLens])
	Options_shot    = make(utils.Options[cons.CameraShot])
	Options_cmotion = make(utils.Options[cons.CameraMotion])
)

func init() {
	Options_lens[cons.LENS_24] = "24"
	Options_lens[cons.CameraLens(35)] = "35"
	Options_lens[cons.CameraLens(50)] = "50"
	Options_lens[cons.CameraLens(85)] = "85"
	Options_lens[cons.CameraLens(100)] = "100"
	Options_lens[cons.CameraLens(200)] = "200"
	Options_lens[cons.CameraLens(300)] = "300"
	Options_lens[cons.CameraLens(400)] = "400"
	Options_lens[cons.CameraLens(600)] = "600"
	//======
	Options_shot[cons.SHOT_CLOSE_UP] = "特写"
	Options_shot[cons.SHOT_CLOSE] = "近景"
	Options_shot[cons.SHOT_MEDIUM] = "中景"
	Options_shot[cons.SHOT_PANORAMIC] = "全景"
	Options_shot[cons.SHOT_LONG] = "远景"
	//=====
	Options_cmotion[cons.MOTION_FIXED] = "固定"
	Options_cmotion[cons.MOTION_PUSH] = "推"
	Options_cmotion[cons.MOTION_PULL] = "拉"
	Options_cmotion[cons.MOTION_FOLLOW] = "跟随"
	Options_cmotion[cons.MOTION_PANNING] = "平移"
	Options_cmotion[cons.MOTION_TURNING] = "摇"
	Options_cmotion[cons.MOTION_SURROUND] = "环绕"
	Options_cmotion[cons.MOTION_SUBJECTIVE] = "主观"
	Options_cmotion[cons.MOTION_SPECIAL] = "特殊"
}
