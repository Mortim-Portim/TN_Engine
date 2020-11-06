package TNE

import (
	"github.com/hajimehoshi/ebiten"
	"marvin/GraphEng/GE"
	"math"
)
const MIN_MOVEMENT_DIF = 0.1

type EntityI interface {
	GE.Drawable
}

type Entity struct {
	drawable *GE.WObj
	Pos, Size [2]float64
	frame *int
}
func (e *Entity) Init(d *GE.WObj, frameCounter *int) {
	e.updateSize()
	e.Pos = [2]float64{-1,-1}
	e.frame = frameCounter
}



//Implements EntityI
func (e *Entity) Draw(screen *ebiten.Image, lv int16, leftTopX, leftTopY, xStart, yStart, sqSize float64) {
	e.drawable.Update(*e.frame)
	e.drawable.Draw(screen, lv, leftTopX, leftTopY, xStart, yStart, sqSize)
}
//Implements EntityI
func (e *Entity) GetPos() (float64, float64, int8) {
	return e.drawable.GetPos()
}
//Implements EntityI
func (e *Entity) Height() float64 {
	return e.drawable.Height()
}
func (e *Entity) MoveBy(dx, dy float64) {
	e.SetPosLT(e.Pos[0]+dx, e.Pos[1]+dy)
}
func (e *Entity) SetPosLT(x,y float64) {
	if math.Abs(x - e.Pos[0]) > MIN_MOVEMENT_DIF || math.Abs(y - e.Pos[1]) > MIN_MOVEMENT_DIF {
		e.drawable.SetToXY(x,y)
		e.Pos = [2]float64{x,y}
	}
}
func (e *Entity) SetPosMD(x,y float64) {
	e.SetPosLT(e.Size[0]/2+x, e.Size[1]/2+y)
}
func (e *Entity) updateSize() {
	w,h := e.drawable.Bounds()
	e.Size = [2]float64{w,h}
}