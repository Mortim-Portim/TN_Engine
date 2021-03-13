package TNE

import (
	"image/color"
	"math"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/mortim-portim/GraphEng/GE"
	cmp "github.com/mortim-portim/GraphEng/compression"
)

type EntityUpdater interface {
	Update(e *Entity, world *SmallWorld)
	Copy() EntityUpdater
}

type Entity struct {
	*Eobj

	ID            int16
	Char          *Character
	ActiveAttacks []Attack

	DeathCause                                    byte
	isDead, justDied                              bool
	Speed                                         float64
	onHealthChange, onStaminaChange, onManaChange func(old, new float64)
	maxHealth, maxStamina, maxMana                float64
	health, stamina, mana                         float64
	showHealth, showStamina, showMana             bool
	UpdateCallBack                                EntityUpdater
}

func (e *Entity) MakeAttackSynced(a Attack, w *SmallWorld) {
	e.Actions().AddAttack(a)
	e.MakeAttackUnSynced(a, w)
}
func (e *Entity) MakeAttackUnSynced(a Attack, w *SmallWorld) {
	e.ActiveAttacks = append(e.ActiveAttacks, a)
	w.Struct.Add_Drawables = w.Struct.Add_Drawables.Add(a)
	a.Start(e, w)
}
func (e *Entity) MoveTiles(tiles float64) {
	e.Eobj.MoveLengthAndFrame(tiles, int(math.Round((tiles/e.Speed)*float64(FPS))))
}
func (e *Entity) Move() {
	e.Eobj.MoveLengthAndFrame(0.017*e.Speed, int(math.Round(0.017*float64(FPS))))
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
		ID:             e.ID,
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
	e.Speed = 6
	return e, nil
}

func (e *Entity) DealDamage(value float64, server bool) {
	if e.IsDead() {
		return
	}
	e.SetHealth(e.Health() - float64(value))
	if server {
		e.CheckDeath(0)
	}
}
func (e *Entity) CheckDeath(cause byte) {
	if e.Health() <= 0 {
		e.setDeadSynced(cause)
	}
}
func (e *Entity) setDead(cause byte) {
	if e.IsDead() {
		return
	}
	e.SetHealth(-1)
	e.isDead = true
	e.justDied = true
	e.DeathCause = cause
}
func (e *Entity) setDeadSynced(cause byte) {
	e.setDead(cause)
	e.actions.AddSetDead(cause)
}
func (e *Entity) IsDead() bool {
	return e.isDead
}
func (e *Entity) JustDied() (b bool) {
	b = e.justDied
	e.justDied = false
	return
}
func (e *Entity) Collides() bool {
	return !e.isDead
}
func (e *Entity) Draw(screen *ebiten.Image, lv int16, leftTopX, leftTopY, xStart, yStart, sqSize float64) {
	if e.IsDead() {
		return
	}
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
	// for _, attack := range e.ActiveAttacks {
	// 	attack.Draw(screen, lv, leftTopX, leftTopY, xStart, yStart, sqSize)
	// }
}

const Health_Stamina_Mana_Bar_Rel = 0.05

func (e *Entity) drawBar(screen *ebiten.Image, idx int, col color.Color, leftTopX, leftTopY, xStart, yStart, sqSize float64, percent float64) {
	bnds := e.Drawbox.Bounds()
	height := Health_Stamina_Mana_Bar_Rel * (bnds.X + bnds.Y)
	width := bnds.X
	y := (e.Drawbox.Min().Y - leftTopY - height*float64(idx+1)) * sqSize
	x := (e.Drawbox.Min().X - leftTopX - (width-bnds.X)/2) * sqSize
	bar := GE.GetHorizontalBar(x+xStart, y+yStart, width*sqSize*float64(percent), height*sqSize)
	bar.Fill(screen, col)
}

func (e *Entity) GetCreationData() (bs []byte) {
	bs = e.Eobj.GetCreationData()
	bs = append(bs, cmp.Float32sToBytes(float32(e.MaxHealth()), float32(e.MaxStamina()), float32(e.MaxMana()), float32(e.Health()), float32(e.Stamina()), float32(e.Mana()))...)
	bs = append(bs, cmp.BoolsToBytes(e.DoesShowHealth(), e.DoesShowStamina(), e.DoesShowMana())[0])
	bs = append(bs, cmp.Int16ToBytes(e.ID)...)
	if e.Char != nil {
		charData := e.Char.ToByte()
		bs = append(bs, charData...)
	}
	return
}
func (e *Entity) OnEobjUpdate(eo *Eobj, w *SmallWorld) {
	if e.UpdateCallBack != nil {
		e.UpdateCallBack.Update(e, w)
	}
	for i, attack := range e.ActiveAttacks {
		attack.Update(e, w)
		if attack.IsFinished() {
			if w != nil && w.HasWorldStruct() {
				err, dws := w.Struct.Add_Drawables.Remove(e.ActiveAttacks[i])
				GE.ShitImDying(err)
				w.Struct.Add_Drawables = dws
			}
			e.ActiveAttacks[i] = nil
		}
	}
	e.RemoveNilAttacks()
}
func (e *Entity) RemoveNilAttacks() {
	rems := 0
	for idx := range e.ActiveAttacks {
		if e.ActiveAttacks[idx-rems] == nil {
			e.RemoveAttackByIdx(idx - rems)
			rems++
		}
	}
}
func (e *Entity) RemoveAttackByIdx(i int) {
	e.ActiveAttacks[i] = e.ActiveAttacks[len(e.ActiveAttacks)-1]
	e.ActiveAttacks = e.ActiveAttacks[:len(e.ActiveAttacks)-1]
}
func (e *Entity) RegisterUpdateCallback(u EntityUpdater) {
	e.UpdateCallBack = u
}
func (e *Entity) Init() {
	if e.Eobj != nil {
		e.Eobj.RegisterUpdateFunc(e.OnEobjUpdate)
	}
}
func (e *Entity) SetMaxHealth(v float64)  { e.maxHealth = v }
func (e *Entity) SetMaxStamina(v float64) { e.maxStamina = v }
func (e *Entity) SetMaxMana(v float64)    { e.maxMana = v }
func (e *Entity) SetHealth(v float64) {
	if v != e.health {
		e.health = v
		if v < 0 {
			e.health = 0.0
		}
		if e.onHealthChange != nil {
			e.onHealthChange(e.health, v)
		}
	}
}
func (e *Entity) SetStamina(v float64) {
	if v != e.stamina {
		e.stamina = v
		if v < 0 {
			e.stamina = 0.0
		}
		if e.onStaminaChange != nil {
			e.onStaminaChange(e.stamina, v)
		}
	}
}
func (e *Entity) SetMana(v float64) {
	if v != e.mana {
		e.mana = v
		if v < 0 {
			e.mana = 0.0
		}
		if e.onManaChange != nil {
			e.onManaChange(e.mana, v)
		}
	}
}
func (e *Entity) ResetHealth()        { e.SetHealth(e.maxHealth) }
func (e *Entity) ResetStamina()       { e.SetStamina(e.maxStamina) }
func (e *Entity) ResetMana()          { e.SetMana(e.maxMana) }
func (e *Entity) ResetHSM()           { e.ResetHealth(); e.ResetStamina(); e.ResetMana() }
func (e *Entity) MaxHealth() float64  { return e.maxHealth }
func (e *Entity) MaxStamina() float64 { return e.maxStamina }
func (e *Entity) MaxMana() float64    { return e.maxMana }
func (e *Entity) Health() float64     { return e.health }
func (e *Entity) Stamina() float64    { return e.stamina }
func (e *Entity) Mana() float64       { return e.mana }
func (e *Entity) HealthPercent() float64 {
	if e.health <= 0 {
		return 0
	}
	return e.health / e.maxHealth
}
func (e *Entity) StaminaPercent() float64 {
	if e.stamina <= 0 {
		return 0
	}
	return e.stamina / e.maxStamina
}
func (e *Entity) ManaPercent() float64 {
	if e.mana <= 0 {
		return 0
	}
	return e.mana / e.maxMana
}
func (e *Entity) ShowHealth(v bool)     { e.showHealth = v }
func (e *Entity) ShowStamina(v bool)    { e.showStamina = v }
func (e *Entity) ShowMana(v bool)       { e.showMana = v }
func (e *Entity) DoesShowHealth() bool  { return e.showHealth }
func (e *Entity) DoesShowStamina() bool { return e.showStamina }
func (e *Entity) DoesShowMana() bool    { return e.showMana }

func (e *Entity) SetOnHealthChange(fnc func(old, new float64))  { e.onHealthChange = fnc }
func (e *Entity) SetOnStaminaChange(fnc func(old, new float64)) { e.onStaminaChange = fnc }
func (e *Entity) SetOnManaChange(fnc func(old, new float64))    { e.onManaChange = fnc }
