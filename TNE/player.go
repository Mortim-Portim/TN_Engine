package TNE

import (
	"github.com/mortim-portim/GraphEng/GE"
)
const MIN_MOVMENT_DIFF = 0.2

//SHOULD load the player from a file, assigning the necassary race and loading all stats
func GetPlayer() *Player {
	p := &Player{}
	return p
}
type Player struct {
	*Entity
}
func (p *Player) Move() {
	if p.Entity.isMoving {
		return
	}
	p.Entity.Move(0.1,1)
}
func (p *Player) GetCreationData() (bs []byte) {
	return []byte{24,46,88,24,66}
}
func GetPlayerByCreationData(bs []byte) (*Player, error) {
	return &Player{Entity:&Entity{}}, nil
}
func (p *Player) Copy() (p2 *Player) {
	p2 = &Player{Entity:p.Entity.Copy()}
	return
}
//updates the player
func (p *Player) Update(w *World, server bool, Collider func(x,y,w,h float64)bool) {
	p.UpdateAll(w, server, Collider)
}
//moves the world to the players position
func (p *Player) MoveWorld(w *GE.WorldStructure) {
	xip, yip := p.IntPos()
	xiw, yiw := w.Middle()	
	ent := p.Entity
	if ent.IsMoving() || xip != int64(xiw) || yip != int64(yiw) {
		pX, pY, _ := p.GetPos()
		w.SetMiddleSmooth(pX-0.5, pY-0.5)
	}
}