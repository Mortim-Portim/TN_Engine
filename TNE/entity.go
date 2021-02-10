package TNE

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/mortim-portim/GraphEng/GE"
	"image/color"
)
type EntityUpdater interface {
	Update(e *Entity, world *World)
	Copy() EntityUpdater
}

type Entity struct {
	*Eobj
	
	maxHealth, maxStamina, maxMana float32
	health, stamina, mana float32
	showHealth, showStamina, showMana bool
	UpdateCallBack EntityUpdater
}
func (e *Entity) Init() {
	if e.Eobj != nil {
		e.Eobj.RegisterUpdateFunc(e.OnEobjUpdate)
	}
}
func (e *Entity) Copy() (e2 *Entity) {
	e2 = &Entity{
		Eobj:e.Eobj.Copy(),
		UpdateCallBack:e.UpdateCallBack,
		maxHealth: e.maxHealth,
		maxStamina: e.maxStamina,
		maxMana: e.maxMana,
		health: e.health,
		stamina: e.stamina,
		mana: e.mana,
		showHealth: e.showHealth,
		showStamina: e.showStamina,
		showMana: e.showMana,
	}
	e2.Init()
	return 
}
func (e *Entity) SetMaxHealth(v float32) {
	e.maxHealth = v
}
func (e *Entity) SetMaxStamina(v float32) {
	e.maxStamina = v
}
func (e *Entity) SetMaxMana(v float32) {
	e.maxMana = v
}
func (e *Entity) ShowHealth(v bool) {e.showHealth = v}
func (e *Entity) ShowStamina(v bool) {e.showStamina = v}
func (e *Entity) ShowMana(v bool) {e.showMana = v}
func (e *Entity) MaxHealth() float32 {return e.maxHealth}
func (e *Entity) MaxStamina() float32 {return e.maxStamina}
func (e *Entity) MaxMana() float32 {return e.maxMana}
func (e *Entity) DoesShowHealth() bool {return e.showHealth}
func (e *Entity) DoesShowStamina() bool {return e.showStamina}
func (e *Entity) DoesShowMana() bool {return e.showMana}

const Health_Stamina_Mana_Bar_Rel = 0.05
func (e *Entity) drawBar(screen *ebiten.Image, idx int, col color.Color, leftTopX, leftTopY, xStart, yStart, sqSize float64) {
	bnds := e.Drawbox.Bounds()
	height := Health_Stamina_Mana_Bar_Rel*(bnds.X+bnds.Y)
	width := bnds.X
	y := (e.Drawbox.Min().Y-leftTopY-height*float64(idx+1))*sqSize
	x := (e.Drawbox.Min().X-leftTopX-(width-bnds.X)/2)*sqSize
	bar := GE.GetHorizontalBar(x+xStart,y+yStart, width*sqSize, height*sqSize)
	bar.Fill(screen, col)
}
func (e *Entity) Draw(screen *ebiten.Image, lv int16, leftTopX, leftTopY, xStart, yStart, sqSize float64) {
	e.Eobj.Draw(screen, lv, leftTopX, leftTopY, xStart, yStart, sqSize)
	if e.showHealth {
		e.drawBar(screen, 0, color.RGBA{255,0,0,255}, leftTopX, leftTopY, xStart, yStart, sqSize)
	}
	if e.showStamina {
		e.drawBar(screen, 1, color.RGBA{0,255,0,255}, leftTopX, leftTopY, xStart, yStart, sqSize)
	}
	if e.showMana {
		e.drawBar(screen, 2, color.RGBA{0,0,255,255}, leftTopX, leftTopY, xStart, yStart, sqSize)
	}
}
func (e *Entity) OnEobjUpdate(eo *Eobj, w *World) {
	if e.UpdateCallBack != nil {
		e.UpdateCallBack.Update(e, w)
	}
}
func (e *Entity) RegisterUpdateCallback(u EntityUpdater) {
	e.UpdateCallBack = u
}
func LoadEntity(path string, frameCounter *int, c *chan bool) (*Entity, error) {
	eo, err := LoadEobj(path, frameCounter, c)
	if err != nil {return nil,err}
	e := &Entity{Eobj:eo}
	e.Init()
	
	//Set test values
	
	
	return e, nil
}