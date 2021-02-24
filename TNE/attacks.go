package TNE

import "github.com/mortim-portim/GraphEng/GE"

/**
Add the index of every Attack as a constant
**/
const (
	ATTACK_FIREBALL = iota
)

type attackparams interface {
	createattack() Attack
}

type projectileattackparams struct {
	name   string
	damage int
	speed  float64
}

func (param projectileattackparams) createattack() Attack {
	return &ProjectileAttack{}
}

/**
Add every Attack to this list according to its index
**/
var AttackGetter = []func(e *Entity, x, y float64) Attack{}

type Attack interface {
	GE.Drawable

	/**
	-> Starts and initializes the attack
	-> pl is the player who started the attack
	-> w is the world, on the client w = nil
	-> if w != nil the attack should modifiy other entities (health) that it hits
	-> if w == nil the attack should just be displayed on the client
	**/
	Start(e *Entity, w *World)
	/**
	-> updates the attack
	-> is called every frame
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

	return
}

type ProjectileAttack struct {
	*GE.WObj
	rotation, speed float64
}

func (attack *ProjectileAttack) Start(pl *Player, w *World, x, y float64) {

}

func (attack *ProjectileAttack) Update(pl *Player, w *World) {

}

func (attack *ProjectileAttack) IsFinished() bool {

}

func (attack *ProjectileAttack) ToBytes() []byte {

}
