package TNE

import (
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

	DialogEntity      *Entity
	ShowsDialogSymbol bool
	DialogSymbol      *GE.ImageObj
}

func (p *Player) Move() {
	if p.Entity.isMoving {
		return
	}
	p.Entity.Move(0.1, 1)
}
func (p *Player) GetCreationData() (bs []byte) {
	return []byte{24, 46, 88, 24, 66}
}
func GetPlayerByCreationData(bs []byte) (*Player, error) {
	return &Player{Entity: &Entity{}}, nil
}
func (p *Player) Copy() (p2 *Player) {
	p2 = &Player{Entity: p.Entity.Copy()}
	return
}

const INTERACTION_DISTANCE = 1

func (p *Player) CheckNearbyDialogs(syncEnts ...*SyncEntity) {
	if p.ShowsDialogSymbol {
		if p.DialogEntity.Hitbox.GetMiddle().DistanceTo(p.Hitbox.GetMiddle()) > INTERACTION_DISTANCE {
			p.ShowsDialogSymbol = false
		} else {
			return
		}
	}
	for _, syncEnt := range syncEnts {
		if syncEnt.HasEntity() {
			ent := syncEnt.Entity
			dis := ent.Hitbox.GetMiddle().DistanceTo(p.Hitbox.GetMiddle())
			if dis <= INTERACTION_DISTANCE {
				p.ShowsDialogSymbol = true
				p.DialogEntity = ent
				break
			}
		}
	}
}

//updates the player
func (p *Player) Update(w *World, server bool, Collider func(x, y, w, h float64) bool) {
	p.UpdateAll(w, server, Collider)
}
func (p *Player) Draw(screen *ebiten.Image, lv int16, leftTopX, leftTopY, xStart, yStart, sqSize float64) {
	p.Entity.Draw(screen, lv, leftTopX, leftTopY, xStart, yStart, sqSize)
	if p.ShowsDialogSymbol && p.DialogSymbol != nil {
		p.DialogSymbol.ScaleToX(p.Drawbox.Bounds().X * sqSize)
		p.DialogSymbol.Y = (p.Drawbox.Min().Y-leftTopY)*sqSize - p.DialogSymbol.H + yStart
		p.DialogSymbol.X = (p.Drawbox.Min().X-leftTopX)*sqSize + xStart
		p.DialogSymbol.Draw(screen)
	}
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
