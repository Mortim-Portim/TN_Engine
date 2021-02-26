package TNE

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
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

	CurrentAttack     int
	DialogEntity      *SyncEntity
	ShowsDialogSymbol bool
	DialogSymbol      *GE.ImageObj
}

func (p *Player) ChangeToAttack(idx int) {
	i := p.Char.GetAttack(idx)
	if i >= 0 {
		p.CurrentAttack = i
	}
}

func (p *Player) StartAttack(x, y float64, sm *SmallWorld) {
	px, py, _ := p.GetMiddle()
	fmt.Printf("Starting Attack %v to (%0.2f|%0.2f), from player at (%0.2f|%0.2f)\n", p.CurrentAttack, x, y, px, py)
	a := Attacks[p.CurrentAttack].Createattack(p.Entity, x, y, nil)

	for _, att := range a {
		p.MakeAttackSynced(att, sm)
	}
}

//Move (tiles float64) - moves the player at his own speed
func (p *Player) Move(tiles float64) {
	if p.Entity.isMoving {
		return
	}
	p.Entity.Move()
}

//GetCreationData () (bs []byte) - returns the players creation data
func (p *Player) GetCreationData() (bs []byte) {
	return []byte{24, 46, 88, 24, 66}
}

//GetPlayerByCreationData (bs []byte) (*Player, error) - tries to create an player from the given data
func GetPlayerByCreationData(bs []byte) (*Player, error) {
	return &Player{Entity: &Entity{}}, nil
}

//Copy () (p2 *Player)
func (p *Player) Copy() (p2 *Player) {
	p2 = &Player{Entity: p.Entity.Copy()}
	return
}

//INTERACTION_DISTANCE - maximum distance to interact with other entities
const INTERACTION_DISTANCE = 1

//CheckNearbyDialogs (syncEnts ...*SyncEntity)
func (p *Player) CheckNearbyDialogs(syncEnts ...*SyncEntity) bool {
	oldE := p.DialogEntity
	min := float64(INTERACTION_DISTANCE)
	for _, syncEnt := range syncEnts {
		if syncEnt.HasEntity() {
			dis := syncEnt.Entity.Hitbox.GetMiddle().DistanceTo(p.Hitbox.GetMiddle())
			if dis <= min {
				min = dis
				p.ShowsDialogSymbol = true
				p.DialogEntity = syncEnt
			}
		}
	}
	if min == float64(INTERACTION_DISTANCE) {
		p.ShowsDialogSymbol = false
	} else if oldE != p.DialogEntity {
		return true
	}
	return false
}

//Update (w *World, server bool, Collider func(x, y, w, h float64) updates the player
func (p *Player) Update(w *SmallWorld, server bool, Collider func(x, y, w, h float64) bool) {
	p.UpdateAll(w, server, Collider)
}

//Draw (screen *ebiten.Image, lv int16, leftTopX, leftTopY, xStart, yStart, sqSize float64) - implements GE.Drawable
func (p *Player) Draw(screen *ebiten.Image, lv int16, leftTopX, leftTopY, xStart, yStart, sqSize float64) {
	p.Entity.Draw(screen, lv, leftTopX, leftTopY, xStart, yStart, sqSize)
	if p.ShowsDialogSymbol && p.DialogSymbol != nil {
		p.DialogSymbol.ScaleToX(p.DialogEntity.Entity.Drawbox.Bounds().X * sqSize)
		p.DialogSymbol.Y = (p.DialogEntity.Entity.Drawbox.Min().Y-leftTopY)*sqSize - p.DialogSymbol.H + yStart
		p.DialogSymbol.X = (p.DialogEntity.Entity.Drawbox.Min().X-leftTopX)*sqSize + xStart
		p.DialogSymbol.Draw(screen)
	}
}

//MoveWorld (w *GE.WorldStructure) - moves the world to the players position
func (p *Player) MoveWorld(w *GE.WorldStructure) {
	xip, yip := p.IntPos()
	xiw, yiw := w.Middle()
	ent := p.Entity
	if ent.IsMoving() || xip != int64(xiw) || yip != int64(yiw) {
		pX, pY, _ := p.GetPos()
		w.SetMiddleSmooth(pX-0.5, pY-0.5)
	}
}
