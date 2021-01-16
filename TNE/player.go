package TNE

import (
	"github.com/mortim-portim/GraphEng/GE"
	//"fmt"
	//"math"
)
const MIN_MOVMENT_DIFF = 0.2

//SHOULD load the player from a file, assigning the necassary race and loading all stats
func GetPlayer() *Player {
	p := &Player{}
	return p
}
type Player struct {
	*Race
}
func (p *Player) Move() {
	p.Entity.Move(0.1,1)
}
func (p *Player) GetCreationData() (bs []byte) {
	
	return
}
func GetPlayerByCreationData(bs []byte) (error, *Player) {
	return nil, &Player{Race:&Race{}}
}
func (p *Player) Copy() (p2 *Player) {
	p2 = &Player{Race:p.Race.Copy()}
	return
}
//updates the player
func (p *Player) Update(w *World) {
	p.Entity.UpdateAll(w)
}
//moves the world to the players position
func (p *Player) MoveWorld(w *GE.WorldStructure) {
	xip, yip := p.IntPos()
	xiw, yiw := w.Middle()	
	ent := p.Race.Entity
	if ent.IsMoving() || xip != int64(xiw) || yip != int64(yiw) {
		pX, pY, _ := p.GetPos()
		w.SetMiddleSmooth(pX-0.5, pY-0.5)
	}
}