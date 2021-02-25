package TNE

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/mortim-portim/GraphEng/GE"
)

/**
Add the index of every Attack as a constant
**/
const (
	ATTACK_FIREBALL = iota
)

type Attackparams interface {
	Init(img *ebiten.Image)
	Createattack(e *Entity, x, y float64, data interface{}) Attack
	FromBytes(bs []byte) Attack
	GetName() string
}

type Projectileattparam struct {
	Name   string
	Id     int
	Damage int
	Speed  float64
	obj    *GE.WObj
}

func (param Projectileattparam) Init(img *ebiten.Image) {
	daynight := GE.GetDayNightAnim(0, 0, 10, 10, 1, 1, img)
	param.obj = GE.GetWObj(daynight, 5, 5, 0, 0, 16, 0, param.Name)
}

func (param Projectileattparam) Createattack(e *Entity, x, y float64, data interface{}) Attack {
	px, py, _ := e.GetMiddle()
	vector := (&GE.Vector{x - px, y - py, 0}).Normalize().Mul(param.Speed)
	return &ProjectileAttack{WObj: param.obj.Copy(), Projectileattparam: param, direction: vector, finished: false}
}

func (param Projectileattparam) FromBytes(bs []byte) Attack {
	vector := GE.VectorFromBytes(bs[:24])
	return &ProjectileAttack{WObj: param.obj.Copy(), Projectileattparam: param, direction: vector, finished: false}
}

func (param Projectileattparam) GetName() string {
	return param.Name
}

/**
Add every Attack to this list according to its index
**/
var Attacks = []Attackparams{
	Projectileattparam{"Fireball", ATTACK_FIREBALL, 5, 5, nil},
}

type Attack interface {
	GE.Drawable

	/**
	-> Starts and initializes the attack
	-> e is the player who started the attack
	-> w is the world, on the client w = nil
	-> if w != nil the attack should modifiy other entities (health) that it hits
	-> if w == nil the attack should just be displayed on the client
	**/
	Start(e *Entity, w *World)
	/**
	-> updates the attack
	-> is called every frame
	-> e is the player who started the attack
	-> if w != nil the attack should modifiy other entities (health) that it hits
	-> if w == nil the attack should just be displayed on the client
	**/
	Update(e *Entity, w *World)
	/**
	-> returns if the attack is finished and can be deleted
	**/
	IsFinished() bool
	/**
	-> encodes the attack as a byte slice
	-> will be transfered to the server and other clients to reconstruct the attack
	**/
	ToBytes() []byte
}

/**
-> should load the attack from bytes
-> will be followed by a call to Start(pl *Player, w *World, x, y float64)
**/
func GetAttackFromBytes(bs []byte) (a Attack, err error) {
	id := int(bs[0])
	return Attacks[id].FromBytes(bs[1:]), nil
}

type ProjectileAttack struct {
	*GE.WObj
	Projectileattparam
	direction *GE.Vector
	finished  bool
}

func (attack *ProjectileAttack) Start(e *Entity, w *World) {

}

func (attack *ProjectileAttack) Update(e *Entity, w *World) {
	attack.WObj.MoveBy(attack.direction.X, attack.direction.Y)
}

func (attack *ProjectileAttack) IsFinished() bool {
	return false
}

func (attack *ProjectileAttack) ToBytes() []byte {
	bytarray := make([]byte, 0)
	bytarray = append(bytarray, byte(attack.Id))
	bytarray = append(bytarray, attack.direction.ToBytes()...)
	return bytarray
}
