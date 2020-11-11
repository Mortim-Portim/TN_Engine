package TNE

import (
	"github.com/hajimehoshi/ebiten"
	cmp "marvin/GraphEng/Compression"
	"marvin/GraphEng/GE"
	"strings"
	"math"
)
const (
	CREATURE_ANIM_STANDARD = "idle_L"
	CREATURE_ANIM_IDLE_L = 0
	CREATURE_ANIM_IDLE_R = 1
	CREATURE_ANIM_IDLE_U = 2
	CREATURE_ANIM_IDLE_D = 3
	CREATURE_ANIM_RUNNING_L = 4
	CREATURE_ANIM_RUNNING_R = 5
	CREATURE_ANIM_RUNNING_U = 6
	CREATURE_ANIM_RUNNING_D = 7
)
const CREATURE_WOBJ = "#WOBJ"

type Creature struct {
	e *Entity
	anims []*GE.DayNightAnim
	currentAnim uint8
	
	xPos, yPos int64
	orientation, neworientation uint8
	isMoving bool
	
	movingFrames, movedFrames int
	movingStepSize float64
	
	isDirty, dataNeedsUpdate bool
	data []byte
	
	factoryCreationId int16
}
//Copys the Creature
func (c *Creature) Copy() (c2 *Creature) {
	c2 = &Creature{e:c.e.Copy(), currentAnim:c.currentAnim, xPos:c.xPos, yPos:c.yPos, orientation:c.orientation,
		neworientation:c.neworientation, isMoving:c.isMoving, movingFrames:c.movingFrames, movedFrames:c.movedFrames,
		isDirty:c.isDirty, dataNeedsUpdate:c.dataNeedsUpdate, factoryCreationId:c.factoryCreationId}
	c2.anims = make([]*GE.DayNightAnim, len(c.anims))
	for i,anim := range(c.anims) {
		if anim != nil {
			c2.anims[i] = anim.Copy()
		}
	}
	copy(c2.data, c.data)
	return
}
//Updates the movement and calls the provided Update func afterwards
func (c *Creature) UpdateAll(w *World) {
	if c.isMoving {
		if c.movedFrames >= c.movingFrames {
			c.isMoving = false
			c.orientation = c.neworientation
			c.isDirty = true
		}else{
			c.moveInDirection(c.orientation)
			c.movedFrames ++
		}
		c.UpdateOrientationAnim()
	}
	c.e.Update(c, w)
}
//Returns the Bounds of the Creature
func (c *Creature) Bounds() (float64, float64) {
	return c.e.Bounds()
}
//Initiates a move action with a specific lenght an duration
func (c *Creature) Move(length, frames int) {
	c.isMoving = true
	c.isDirty = true
	c.movingFrames = frames
	c.movedFrames = 0
	c.movingStepSize = float64(length)/float64(frames)
	c.neworientation = c.orientation
}
//Sets the middle of the Creature
func (c *Creature) SetMiddleTo(x, y float64) {
	c.e.SetPos(x,y)
	c.setIntPos()
}
//Sets the top left corner of the Creature
func (c *Creature) SetTopLeftTo(x, y float64) {
	c.e.SetToXY(x,y)
	c.setIntPos()
}
func (c *Creature) setIntPos() {
	xf,yf,_ := c.e.GetPos()
	x,y := int64(math.Round(xf-0.5)), int64(math.Round(yf-0.5))
	if x != c.xPos || y != c.yPos {
		c.xPos, c.yPos = x,y
		c.isDirty = true
	}
}
//Changes the orientation
func (c *Creature) ChangeOrientation(newO uint8) {
	if newO != c.orientation {
		if c.isMoving {
			c.neworientation = newO
		}else{
			c.orientation = newO
		}
		c.isDirty = true
	}
}
//Updates the Orientation animation, ONLY call this if really necassary
func (c *Creature) UpdateOrientationAnim() {
	idx := c.orientation
	if c.isMoving {
		idx += 4
	}
	c.SetAnim(int(idx))
}
func (c *Creature) moveInDirection(dir uint8) {
	dx, dy := 0.0,0.0
	switch dir {
		case 0:
			dx = -c.movingStepSize
			break
		case 1:
			dx = c.movingStepSize
			break
		case 2:
			dy = -c.movingStepSize
			break
		case 3:
			dy = c.movingStepSize
			break
	}
	c.e.MoveBy(dx, dy)
	c.setIntPos()
}
//Returns the creatures data
func (c *Creature) GetData() []byte {
	c.UpdateData(false)
	c.isDirty = false
	return c.data
}
//Sets the creatures data
func (c *Creature) SetData(data []byte) {
	c.data = data
	c.isMoving = cmp.ByteToBool(data[0])
	c.orientation = uint8(data[1])
	c.currentAnim = uint8(data[2])
}
//Computes the creatures data
func (c *Creature) UpdateData(force bool) {
	if c.dataNeedsUpdate || force {
		c.data = []byte{cmp.BoolToByte(c.isMoving), byte(c.orientation), byte(c.currentAnim)}
		c.dataNeedsUpdate = false
	}
}

//Implements EntityI
func (c *Creature) Update(_ EntityI, w *World) {
	c.UpdateAll(w)
}
//Implements EntityI
func (c *Creature) Height() float64 {
	return c.e.Height()
}
//Implements EntityI
func (c *Creature) GetPos() (float64, float64, int8) {
	return c.e.GetPos()
}
func (c *Creature) Entity() *Entity {
	return c.e
}
func (c *Creature) IsMoving() bool {
	return c.isMoving
}
func (c *Creature) IntPos() (int64, int64) {
	return c.xPos, c.yPos
}
func (c *Creature) RegiserUpdateFunc(u func(e EntityI, w *World)) {
	c.e.Update = u
}

/**
Loads a Creature from a directory that contains DayNightAnims with the names listed in CREATURE_ANIM_NAMES
Also a WOBJ.txt file is needed, that describes the Attributes of the Creature
**/
func LoadCreature(path string, frameCounter *int) (*Creature, error) {
	if path[len(path)-1:] != "/" {
		path += "/"
	}
	pathS := strings.Split(path, "/"); name := pathS[len(pathS)-2]
	c := &Creature{e:GetEntity(frameCounter, nil), anims:make([]*GE.DayNightAnim, 0)}
	wobj, err := GE.GetWObjFromPath(name, path+CREATURE_ANIM_STANDARD , path+CREATURE_WOBJ)
	if err != nil {return c, err}
	c.e.WObj = *wobj
	
	idx := &GE.List{}
	idx.LoadFromFile(path+INDEX_FILE_NAME)
	names := idx.GetSlice()
	c.anims = make([]*GE.DayNightAnim, 0)
	for _,anim_n := range(names) {
		anim, _ := GE.GetDayNightAnimFromParams(1,1,1,1, path+anim_n+".txt", path+anim_n+".png")
		c.anims = append(c.anims, anim)
	}
	return c, nil
}
func (c *Creature) SetAnim(idx int) {
	if idx < 0 || idx >= len(c.anims) || c.anims[idx] == nil {
		return
	}
	c.currentAnim = uint8(idx)
	c.isDirty = true
	c.e.WObj.SetAnim(c.anims[idx])
}
func (c *Creature) GetAnim() uint8 {
	return c.currentAnim
}
//Implements EntityI
func (c *Creature) Draw(screen *ebiten.Image, lv int16, leftTopX, leftTopY, xStart, yStart, sqSize float64) {
	c.e.WObj.Update(*c.e.frame)
	c.e.WObj.Draw(screen, lv, leftTopX, leftTopY, xStart, yStart, sqSize)
}