package TNE

import (
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
	Race
}
func (p *Player) GetCreationData() (bs []byte) {
	
	return
}
func GetPlayerByCreationData(bs []byte) (error, *Player) {
	return nil, nil
}
func (p *Player) Copy() (p2 *Player) {
	p2 = &Player{Race:*p.Race.Copy()}
	return
}
//SHOULD update the player
func (p *Player) Update(w *World) {
	p.Race.UpdateAll(w)
}
//moves the world to the players position
func (p *Player) MoveWorld(w *World) {
	xip, yip := p.IntPos()
	xiw, yiw := w.Structure.Middle()	
	ent := p.Race.Entity
	if ent.IsMoving() || xip != int64(xiw) || yip != int64(yiw) {
		pX, pY, _ := p.GetPos()
		w.Structure.SetMiddleSmooth(pX-0.5, pY-0.5)
	}
}