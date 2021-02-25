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

/**
Add every Attack to this list according to its index
**/
var Attacks = []Attackparams{
	&Projectileattparam{"Fireball", ATTACK_FIREBALL, 5, 0.2, nil},
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
	Start(e *Entity, w *SmallWorld)
	/**
	-> updates the attack
	-> is called every frame
	-> e is the player who started the attack
	-> if w != nil the attack should modifiy other entities (health) that it hits
	-> if w == nil the attack should just be displayed on the client
	**/
	Update(e *Entity, w *SmallWorld)
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
