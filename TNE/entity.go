package TNE

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten"
	cmp "github.com/mortim-portim/GraphEng/Compression"
	"github.com/mortim-portim/GraphEng/GE"
)

type EntityUpdater interface {
	Update(e *Entity, world *World)
	Copy() EntityUpdater
}

type Entity struct {
	*Eobj

	Char *Character

	Speed                             float64
	maxHealth, maxStamina, maxMana    float32
	health, stamina, mana             float32
	showHealth, showStamina, showMana bool
	UpdateCallBack                    EntityUpdater
}

func (e *Entity) Move(tiles float64) {
	e.Eobj.MoveLengthAndFrame(tiles, int(math.Round((tiles/e.Speed)*float64(FPS))))
}
func (e *Entity) Copy() (e2 *Entity) {
	e2 = &Entity{
		Eobj:           e.Eobj.Copy(),
		UpdateCallBack: e.UpdateCallBack,
		maxHealth:      e.maxHealth,
		maxStamina:     e.maxStamina,
		maxMana:        e.maxMana,
		health:         e.health,
		stamina:        e.stamina,
		mana:           e.mana,
		showHealth:     e.showHealth,
		showStamina:    e.showStamina,
		showMana:       e.showMana,
		Char:           e.Char.Copy(),
		Speed:          e.Speed,
	}
	e2.Init()
	return
}
func LoadEntity(path string, frameCounter *int) (*Entity, error) {
	eo, err := LoadEobj(path, frameCounter)
	if err != nil {
		return nil, err
	}
	e := &Entity{Eobj: eo}
	e.Init()
	e.SetMaxHealth(80)
	e.SetMaxStamina(100)
	e.SetMaxMana(130)
	e.ResetHSM()
	e.Speed = 3
	return e, nil
}

func (e *Entity) Draw(screen *ebiten.Image, lv int16, leftTopX, leftTopY, xStart, yStart, sqSize float64) {
	e.Eobj.Draw(screen, lv, leftTopX, leftTopY, xStart, yStart, sqSize)
	if e.showHealth {
		e.drawBar(screen, 0, color.RGBA{255, 0, 0, 255}, leftTopX, leftTopY, xStart, yStart, sqSize, e.HealthPercent())
	}
	if e.showStamina {
		e.drawBar(screen, 1, color.RGBA{0, 255, 0, 255}, leftTopX, leftTopY, xStart, yStart, sqSize, e.StaminaPercent())
	}
	if e.showMana {
		e.drawBar(screen, 2, color.RGBA{0, 0, 255, 255}, leftTopX, leftTopY, xStart, yStart, sqSize, e.ManaPercent())
	}
}

const Health_Stamina_Mana_Bar_Rel = 0.05

func (e *Entity) drawBar(screen *ebiten.Image, idx int, col color.Color, leftTopX, leftTopY, xStart, yStart, sqSize float64, percent float32) {
	bnds := e.Drawbox.Bounds()
	height := Health_Stamina_Mana_Bar_Rel * (bnds.X + bnds.Y)
	width := bnds.X
	y := (e.Drawbox.Min().Y - leftTopY - height*float64(idx+1)) * sqSize
	x := (e.Drawbox.Min().X - leftTopX - (width-bnds.X)/2) * sqSize
	bar := GE.GetHorizontalBar(x+xStart, y+yStart, width*sqSize*float64(percent), height*sqSize)
	bar.Fill(screen, col)
}

const ENTITY_CREATION_DATA_LENGTH = 48 + CHARACTER_BYTES_LENGTH

func (e *Entity) GetCreationData() (bs []byte) {
	bs = e.Eobj.GetCreationData()
	copy(bs[23:47], cmp.Float32sToBytes(e.MaxHealth(), e.MaxStamina(), e.MaxMana(), e.Health(), e.Stamina(), e.Mana()))
	bs[47] = cmp.BoolsToBytes(e.DoesShowHealth(), e.DoesShowStamina(), e.DoesShowMana())[0]
	if e.Char != nil {
		copy(bs[48:48+CHARACTER_BYTES_LENGTH], e.Char.ToByte())
	}
	return
}
func (e *Entity) OnEobjUpdate(eo *Eobj, w *World) {
	if e.UpdateCallBack != nil {
		e.UpdateCallBack.Update(e, w)
	}
}
func (e *Entity) RegisterUpdateCallback(u EntityUpdater) {
	e.UpdateCallBack = u
}
func (e *Entity) Init() {
	if e.Eobj != nil {
		e.Eobj.RegisterUpdateFunc(e.OnEobjUpdate)
	}
}
func (e *Entity) SetMaxHealth(v float32)  { e.maxHealth = v }
func (e *Entity) SetMaxStamina(v float32) { e.maxStamina = v }
func (e *Entity) SetMaxMana(v float32)    { e.maxMana = v }
func (e *Entity) SetHealth(v float32)     { e.health = v }
func (e *Entity) SetStamina(v float32)    { e.stamina = v }
func (e *Entity) SetMana(v float32)       { e.mana = v }
func (e *Entity) ResetHealth()            { e.health = e.maxHealth }
func (e *Entity) ResetStamina()           { e.stamina = e.maxStamina }
func (e *Entity) ResetMana()              { e.mana = e.maxMana }
func (e *Entity) ResetHSM()               { e.ResetHealth(); e.ResetStamina(); e.ResetMana() }
func (e *Entity) MaxHealth() float32      { return e.maxHealth }
func (e *Entity) MaxStamina() float32     { return e.maxStamina }
func (e *Entity) MaxMana() float32        { return e.maxMana }
func (e *Entity) Health() float32         { return e.health }
func (e *Entity) Stamina() float32        { return e.stamina }
func (e *Entity) Mana() float32           { return e.mana }
func (e *Entity) HealthPercent() float32  { return e.health / e.maxHealth }
func (e *Entity) StaminaPercent() float32 { return e.stamina / e.maxStamina }
func (e *Entity) ManaPercent() float32    { return e.mana / e.maxMana }
func (e *Entity) ShowHealth(v bool)       { e.showHealth = v }
func (e *Entity) ShowStamina(v bool)      { e.showStamina = v }
func (e *Entity) ShowMana(v bool)         { e.showMana = v }
func (e *Entity) DoesShowHealth() bool    { return e.showHealth }
func (e *Entity) DoesShowStamina() bool   { return e.showStamina }
func (e *Entity) DoesShowMana() bool      { return e.showMana }
