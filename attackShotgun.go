package TNE

import (
	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/mortim-portim/GraphEng/GE"
)

type ShotgunAttackParam struct {
	ProjectileAttParam
	count    int
	rotation float64
}

func (param *ShotgunAttackParam) Init(img *ebiten.Image) {
	daynight := GE.GetDayNightAnim(0, 0, 10, 10, 10, 0, img)
	param.obj = GE.GetWObj(daynight, 0.42, 0.42, 0, 0, 24, 10, param.Name)
}

func (param *ShotgunAttackParam) Createattack(e *Entity, x, y float64, data interface{}) []Attack {
	px, py, _ := e.GetMiddle()
	maindir := (&GE.Vector{x - px, y - py, 0}).Normalize().Mul(param.Speed)
	attacklist := make([]Attack, param.count)

	for i := range attacklist {
		dir := maindir.Copy()
		dir.RotateZ(param.rotation * (float64(i)/float64(param.count-1) - 0.5))
		attacklist[i] = param.createProjectileAtt(dir, px, py)
	}
	return attacklist
}
