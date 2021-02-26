package TNE

/*

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/mortim-portim/GraphEng/GE"
)

type ShotgunAttackParam struct {
	ProjectileAttParam
	count, rotation int
}

func (param *ShotgunAttackParam) Init(img *ebiten.Image) {
	daynight := GE.GetDayNightAnim(0, 0, 10, 10, 10, 0, img)
	param.obj = GE.GetWObj(daynight, 0.42, 0.42, 0, 0, 24, 10, param.Name)
}

func (param *ShotgunAttackParam) Createattack(e *Entity, x, y float64, data interface{}) []Attack {
	px, py, _ := e.GetMiddle()
	maindir := (&GE.Vector{x - px, y - py, 0}).Normalize().Mul(param.Speed)
	attacklist := make([]Attack)
	return []Attack{param.createProjectileAtt(dir, px, py)}
}

/*func (param *ShotgunAttackParam) FromBytes(bs []byte) Attack {
	dir := GE.XYVectorFromBytes(bs[:16])
	pos := GE.XYVectorFromBytes(bs[16:])
	return param.createProjectileAtt(dir, pos.X, pos.Y)
}
func (param *ShotgunAttackParam) createProjectileAtt(dir *GE.Vector, px, py float64) Attack {
	nWobj := param.obj.Copy()
	nWobj.SetMiddle(px, py)
	nWobj.GetAnim().SetRotation(dir.GetRotationZ())
	return &ProjectileAttack{WObj: nWobj, ProjectileAttParam: param, direction: dir, finished: false}
}
*/
