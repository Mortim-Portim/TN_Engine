package TNE

import "github.com/mortim-portim/GraphEng/GE"

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

type Attack interface {
	GE.Drawable

	/**
	-> Starts and initializes the attack
	-> pl is the player who started the attack
	-> w is the world, on the client w = nil
	-> if w != nil the attack should modifiy other entities (health) that it hits
	**/
	Start(pl *Player, w *World, x, y float64)
	/**
	-> updates the attack
	-> is called every frame
	**/
	Update(pl *Player, w *World)
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
