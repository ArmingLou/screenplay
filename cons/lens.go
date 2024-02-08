package cons

type CameraLens int   //焦段
type CameraShot int   //景别
type CameraMotion int //运镜

const (
	SHOT_CLOSE_UP  CameraShot = iota + 1 //特写
	SHOT_CLOSE                           //近景
	SHOT_MEDIUM                          //中景
	SHOT_PANORAMIC                       //全景
	SHOT_LONG                            //远景
)

const (
	LENS_24 CameraLens = 24
)

const (
	MOTION_FIXED      CameraMotion = iota + 1 //固定镜头
	MOTION_PUSH                               //推镜头
	MOTION_PULL                               //拉镜头
	MOTION_FOLLOW                             //跟随镜头
	MOTION_PANNING                            //平移镜头
	MOTION_TURNING                            //摇镜头
	MOTION_SURROUND                           //环绕镜头
	MOTION_SUBJECTIVE                         //主观镜头
	MOTION_SPECIAL                            //特殊镜头运动
)

type ColumeType int

const (
	COLUME_TYPE_ASIDE ColumeType = iota + 1
	COLUME_TYPE_SCENE
	COLUME_TYPE_LENS
	COLUME_TYPE_SHOT
	COLUME_TYPE_MOTION
	COLUME_TYPE_CMERA
	COLUME_TYPE_ACTION
	COLUME_TYPE_TRANS
	COLUME_TYPE_VOICE
	COLUME_TYPE_EXT
)
