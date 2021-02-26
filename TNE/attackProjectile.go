package TNE

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/mortim-portim/GraphEng/GE"
)

//Creates an Attack that consist of one Projectile

type ProjectileAttParam struct {
	Name                           string
	Id                             int
	Damage                         float32
	Speed, Range                   float64
	obj                            *GE.WObj
	HitboxW, HitboxH               float64
	spriteWidth, squareSize, layer int
}

func (param *ProjectileAttParam) Init(img *ebiten.Image) {
	daynight := GE.GetDayNightAnim(0, 0, 10, 10, param.spriteWidth, 0, img)
	param.obj = GE.GetWObj(daynight, param.HitboxW, param.HitboxH, 0, 0, param.squareSize, int8(param.layer), param.Name)
}

func (param *ProjectileAttParam) Createattack(e *Entity, x, y float64, data interface{}) []Attack {
	px, py, _ := e.GetMiddle()
	dir := (&GE.Vector{x - px, y - py, 0}).Normalize().Mul(param.Speed)
	return []Attack{param.createProjectileAtt(dir, px, py)}
}

func (param *ProjectileAttParam) FromBytes(bs []byte) Attack {
	dir := GE.XYVectorFromBytes(bs[:16])
	pos := GE.XYVectorFromBytes(bs[16:])
	return param.createProjectileAtt(dir, pos.X, pos.Y)
}
func (param *ProjectileAttParam) createProjectileAtt(dir *GE.Vector, px, py float64) Attack {
	nWobj := param.obj.Copy()
	nWobj.SetMiddle(px, py)
	nWobj.GetAnim().SetRotation(dir.GetRotationZ())
	return &ProjectileAttack{WObj: nWobj, ProjectileAttParam: param, direction: dir, finished: false}
}

func (param *ProjectileAttParam) GetName() string {
	return param.Name
}

type ProjectileAttack struct {
	*GE.WObj
	*ProjectileAttParam
	direction *GE.Vector
	finished  bool
	frame     float64
}

func (attack *ProjectileAttack) Start(e *Entity, w *SmallWorld) {
	idx := EOBJ_ATTACKING_RIGHT
	if e.currentAnim%2 == 0 {
		idx = EOBJ_ATTACKING_LEFT
	}
	e.SetAnimManual(idx)
}

func (attack *ProjectileAttack) Update(e *Entity, w *SmallWorld) {
	attack.WObj.MoveBy(attack.direction.X, attack.direction.Y)
	attack.frame++

	OnRectWithWorldStructObjCollision(attack.Hitbox, w.Struct, func(so *GE.StructureObj, ent *Entity, ply *Player) {
		if ply != nil && e.ID == ply.ID {
			return
		}
		if ply != nil {
			ply.DealDamage(attack.Damage)
		}
		attack.finished = true
	})

	if attack.frame >= attack.Range/attack.Speed {
		attack.finished = true
	}
}

func (attack *ProjectileAttack) IsFinished() bool {
	return attack.finished
}

func (attack *ProjectileAttack) ToBytes() []byte {
	bytarray := make([]byte, 0)
	bytarray = append(bytarray, byte(attack.Id))
	bytarray = append(bytarray, attack.direction.XYToBytes()...)
	x, y, _ := attack.GetMiddle()
	bytarray = append(bytarray, (&GE.Vector{x, y, 0}).XYToBytes()...)
	return bytarray
}
