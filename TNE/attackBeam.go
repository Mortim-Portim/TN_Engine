package TNE

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/mortim-portim/GraphEng/GE"
)

type BeamAttParam struct {
	*ProjectileAttParam
	Count int
}

func (param *BeamAttParam) Init(img *ebiten.Image) {
	daynight := GE.GetDayNightAnim(0, 0, 10, 10, param.spriteWidth, 0, img)
	param.obj = GE.GetWObj(daynight, param.HitboxW, param.HitboxH, 0, 0, param.squareSize, int8(param.layer), param.Name)
}

func (param *BeamAttParam) Createattack(e *Entity, x, y float64, data interface{}) []Attack {
	px, py, _ := e.GetMiddle()
	dir := (&GE.Vector{x - px, y - py, 0}).Normalize().Mul(param.Speed)

	attacks := make([]Attack, param.Count)
	for i := 0; i < param.Count; i++ {
		attacks[i] = param.createProjectileAtt(dir, px, py, float64(i*1))
	}
	return attacks
}

func (param *BeamAttParam) FromBytes(bs []byte) Attack {
	dir := GE.XYVectorFromBytes(bs[:16])
	pos := GE.VectorFromBytes(bs[16:])
	return param.createProjectileAtt(dir, pos.X, pos.Y, pos.Z)
}
func (param *BeamAttParam) createProjectileAtt(dir *GE.Vector, px, py, delay float64) Attack {
	nWobj := param.obj.Copy()
	nWobj.SetMiddle(px, py)
	nWobj.GetAnim().SetRotation(dir.GetRotationZ())
	return &BeamAttack{WObj: nWobj, BeamAttParam: param, direction: dir, finished: false, delay: delay}
}

func (param *BeamAttParam) GetName() string {
	return param.Name
}

type BeamAttack struct {
	*GE.WObj
	*BeamAttParam
	direction    *GE.Vector
	finished     bool
	frame, delay float64
}

func (attack *BeamAttack) Start(e *Entity, w *SmallWorld) {
	idx := EOBJ_ATTACKING_RIGHT
	if e.currentAnim%2 == 0 {
		idx = EOBJ_ATTACKING_LEFT
	}
	e.SetAnimManual(idx)
}

func (attack *BeamAttack) Update(e *Entity, w *SmallWorld) {
	attack.frame++

	if attack.frame < attack.delay {
		return
	}

	attack.WObj.MoveBy(attack.direction.X, attack.direction.Y)

	OnRectWithWorldStructObjCollision(attack.Hitbox, w.Struct, func(so *GE.StructureObj, ent *Entity, ply *Player) {
		if ply != nil && e.ID == ply.ID {
			return
		}
		if ply != nil {
			ply.DealDamage(attack.Damage, w.IsOnServer())
		}
		attack.finished = true
	})

	if attack.frame+attack.delay >= attack.Range/attack.Speed {
		attack.finished = true
	}
}

func (attack *BeamAttack) IsFinished() bool {
	return attack.finished
}

func (attack *BeamAttack) ToBytes() []byte {
	bytarray := make([]byte, 0)
	bytarray = append(bytarray, byte(attack.Id))
	bytarray = append(bytarray, attack.direction.XYToBytes()...)
	x, y, _ := attack.GetMiddle()
	bytarray = append(bytarray, (&GE.Vector{x, y, attack.delay}).ToBytes()...)
	return bytarray
}
