package TNE

import (
	"github.com/hajimehoshi/ebiten"
	"marvin/GraphEng/GE"
	"math"
)
const (
	CREATURE_ANIM_IDLE_L = 0
	CREATURE_ANIM_IDLE_R = 1
	CREATURE_ANIM_IDLE_U = 2
	CREATURE_ANIM_IDLE_D = 3
	CREATURE_ANIM_RUNNING_L = 4
	CREATURE_ANIM_RUNNING_R = 5
	CREATURE_ANIM_RUNNING_U = 6
	CREATURE_ANIM_RUNNING_D = 7
)
var CREATURE_ANIM_NAMES = []string{"idle_L","idle_R","idle_U","idle_D", "running_L","running_R","running_U","running_D"}
const CREATURE_WOBJ = "WOBJ"

type Creature struct {
	Entity
	anims []*GE.DayNightAnim
	orientation int
	IsMoving bool
}
func (c *Creature) Copy() (c2 *Creature) {
	c2 = &Creature{Entity:*c.Entity.Copy()}
	c2.anims = make([]*GE.DayNightAnim, len(c.anims))
	for i,anim := range(c.anims) {
		c2.anims[i] = anim.Copy()
	}
	return
}

func (c *Creature) ChangeOrientation(newO int) {
	if newO != c.orientation {
		c.orientation = newO
	}
}
func (c *Creature) UpdateOrientation() {
	idx := c.orientation
	if c.IsMoving {
		idx += 4
	}
	c.SetAnim(idx)
}
func (c *Creature) SetTo(x, y float64) {
	oX, oY,_ := c.GetPos()
	if oX != x || oY != y {
		c.SetPosMD(x,y)
		dx, dy := x-oX, y-oY
		if math.Abs(dx) > math.Abs(dy) {
			if dx < 0 {
				c.ChangeOrientation(0)
			}else{
				c.ChangeOrientation(1)
			}
		}else{
			if dy < 0 {
				c.ChangeOrientation(2)
			}else{
				c.ChangeOrientation(3)
			}
		}
	}
}
func (c *Creature) Move(dx, dy float64) {
	oX, oY,_ := c.GetPos()
	c.SetTo(oX+dx, oY+dx)
}

/**
Loads a Creature from a directory that contains DayNightAnims with the names listed in CREATURE_ANIM_NAMES
Also a WOBJ.txt file is needed, that describes the Attributes of the Creature
**/
func LoadCreature(path string, frameCounter *int) (*Creature, error) {
	if path[len(path)-1:] != "/" {
		path += "/"
	}
	c := &Creature{anims:make([]*GE.DayNightAnim, len(CREATURE_ANIM_NAMES))}
	wobj, err := GE.GetWObjFromPath(path+CREATURE_ANIM_NAMES[CREATURE_ANIM_IDLE_D] , path+CREATURE_WOBJ)
	if err != nil {return c, err}
	c.Entity.drawable = wobj
	
	for i,anim_n := range(CREATURE_ANIM_NAMES) {
		anim, err := GE.GetDayNightAnimFromParams(1,1,1,1, path+anim_n+".txt", path+anim_n+".png")
		if err != nil {return c, err}
		c.anims[i] = anim
	}
	c.Entity.Init(frameCounter)
	return c, nil
}
func (c *Creature) SetAnim(idx int) {
	c.Entity.drawable.SetAnim(c.anims[idx])
}
//Implements EntityI
func (c *Creature) Draw(screen *ebiten.Image, lv int16, leftTopX, leftTopY, xStart, yStart, sqSize float64) {
	c.drawable.Update(*c.frame)
	c.drawable.Draw(screen, lv, leftTopX, leftTopY, xStart, yStart, sqSize)
}