package TNE

import (
	"github.com/hajimehoshi/ebiten"
	//cmp "marvin/GraphEng/Compression"
	"marvin/GraphEng/GE"
	//"math"
)
const MIN_MOVEMENT_DIF = 0.01

type EntityI interface {
	GE.Drawable
	Update(e EntityI, world *World)
	GetData() []byte
	SetData([]byte)
}

func GetEntity(fc *int, u func(e EntityI, world *World)) (e *Entity) {
	e = &Entity{}
	e.Update = u
	e.frame = fc
	return
}
type Entity struct {
	GE.WObj
	frame *int
	Update func(e EntityI, world *World)
}
func (e *Entity) Copy() (e2 *Entity) {
	e2 = &Entity{WObj:*e.WObj.Copy(), frame:e.frame, Update:e.Update}
	return
}
func (e *Entity) GetData() (bs []byte) {
	return
}
func (e *Entity) SetData([]byte) {
	
}
//Implements EntityI
func (e *Entity) Draw(screen *ebiten.Image, lv int16, leftTopX, leftTopY, xStart, yStart, sqSize float64) {
	e.WObj.Update(*e.frame)
	e.WObj.Draw(screen, lv, leftTopX, leftTopY, xStart, yStart, sqSize)
}