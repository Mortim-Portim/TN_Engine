package TNE

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/mortim-portim/GraphEng/GE"
)

type Projectileattparam struct {
	Name   string
	Id     int
	Damage int
	Speed  float64
	obj    *GE.WObj
}

func (param *Projectileattparam) Init(img *ebiten.Image) {
	daynight := GE.GetDayNightAnim(0, 0, 10, 10, 10, 0, img)
	param.obj = GE.GetWObj(daynight, 0.42, 0.42, 0, 0, 24, 10, param.Name)
}

func (param *Projectileattparam) Createattack(e *Entity, x, y float64, data interface{}) Attack {
	px, py, _ := e.GetMiddle()
	dir := (&GE.Vector{x - px, y - py, 0}).Normalize().Mul(param.Speed)
	return param.createProjectileAtt(dir, px, py)
}

func (param *Projectileattparam) FromBytes(bs []byte) Attack {
	dir := GE.XYVectorFromBytes(bs[:16])
	pos := GE.XYVectorFromBytes(bs[16:])
	return param.createProjectileAtt(dir, pos.X, pos.Y)
}
func (param *Projectileattparam) createProjectileAtt(dir *GE.Vector, px, py float64) Attack {
	nWobj := param.obj.Copy()
	nWobj.SetMiddle(px, py)
	nWobj.GetAnim().SetRotation(dir.GetRotationZ())
	return &ProjectileAttack{WObj: nWobj, Projectileattparam: param, direction: dir, finished: false}
}

func (param *Projectileattparam) GetName() string {
	return param.Name
}

type ProjectileAttack struct {
	*GE.WObj
	*Projectileattparam
	direction *GE.Vector
	finished  bool
}

func (attack *ProjectileAttack) Start(e *Entity, w *SmallWorld) {
	idx := EOBJ_ATTACKING_LEFT
	dir := ENTITY_ORIENTATION_L
	if attack.direction.GetRotationZ() < 180 {
		idx = EOBJ_ATTACKING_RIGHT
		dir = ENTITY_ORIENTATION_R
	}
	e.ChangeOrientation((&Direction{ID: dir}).FromID())
	e.SetAnimManual(idx)
}

func (attack *ProjectileAttack) Update(e *Entity, w *SmallWorld) {
	attack.WObj.MoveBy(attack.direction.X, attack.direction.Y)
}

func (attack *ProjectileAttack) IsFinished() bool {
	return false
}

func (attack *ProjectileAttack) ToBytes() []byte {
	bytarray := make([]byte, 0)
	bytarray = append(bytarray, byte(attack.Id))
	bytarray = append(bytarray, attack.direction.XYToBytes()...)
	x, y, _ := attack.GetMiddle()
	bytarray = append(bytarray, (&GE.Vector{x, y, 0}).XYToBytes()...)
	return bytarray
}
