package TNE

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/mortim-portim/GraphEng/GE"

	//cmp "marvin/GraphEng/Compression"
	"errors"
	"fmt"
	"math"
	"strings"

	cmp "github.com/mortim-portim/GraphEng/Compression"
)

const (
	CREATURE_ANIM_STANDARD  = "idle_L"
	CREATURE_ANIM_IDLE_L    = 0
	CREATURE_ANIM_IDLE_R    = 1
	CREATURE_ANIM_IDLE_U    = 2
	CREATURE_ANIM_IDLE_D    = 3
	CREATURE_ANIM_RUNNING_L = 4
	CREATURE_ANIM_RUNNING_R = 5
	CREATURE_ANIM_RUNNING_U = 6
	CREATURE_ANIM_RUNNING_D = 7
)
const CREATURE_WOBJ = "#WOBJ"

var ERR_WRONG_BYTE_LENGTH = errors.New("Wrong byte length")

//changes should consist out of a list slices with the following meanings
const (
	ENTITY_START_MOVE         = byte(0)
	ENTITY_KEEP_MOVING        = byte(1)
	ENTITY_CHANGE_ORIENTATION = byte(2)
)

type EntityUpdater func(e *Entity, world *World)

type Entity struct {
	GE.WObj

	anims       []*GE.DayNightAnim
	currentAnim uint8

	xPos, yPos                  int64
	orientation, neworientation uint8
	isMoving, keepMoving        bool

	movingFrames, movedFrames int
	movingStepSize            float64

	//IsDirty should only change when an action is started (Move, KeepMoving, ChangeOrientation)
	changed, isDirty bool

	factoryCreationId int16

	frame   *int
	Updater EntityUpdater
}

func (e *Entity) GetDelta() []byte {
	return []byte{byte(e.orientation), byte(e.currentAnim)}
}
func (e *Entity) SetDelta(bs []byte) {
	e.orientation = uint8(bs[0])
	e.currentAnim = uint8(bs[1])
}

//Loads an entity from full data len(data) = 37
func (cf *EntityFactory) LoadEntityFromFullData(data []byte) (*Entity, error) {
	if len(data) != 37 {
		return nil, ERR_WRONG_BYTE_LENGTH
	}
	fcID := cmp.BytesToInt16(data[35:36])
	e := cf.Get(int(fcID))
	e.currentAnim = uint8(data[0])
	e.xPos = cmp.BytesToInt64(data[1:8])
	e.yPos = cmp.BytesToInt64(data[9:16])
	e.orientation = uint8(data[17])
	e.neworientation = uint8(data[18])
	e.isMoving = cmp.ByteToBool(data[19])
	e.keepMoving = cmp.ByteToBool(data[20])
	e.movingFrames = int(cmp.BytesToInt16(data[21:22]))
	e.movedFrames = int(cmp.BytesToInt16(data[23:24]))
	e.movingStepSize = cmp.BytesToFloat64(data[25:32])
	e.changed = cmp.ByteToBool(data[33])
	e.isDirty = cmp.ByteToBool(data[34])
	return e, nil
}

//(1)currentAnim| (8)xPos| (8)yPos| (1)orientation| (1)neworientation| (1)isMoving| (1)keepMoving| (2)movingFrames|
//(2)movedFrames| (8)movingStepSize| (1)changed| (1)isDirty| (2)factoryCreationId| len() = 37
func (e *Entity) FullData() (data []byte) {
	data = make([]byte, 37)
	data[0] = byte(e.currentAnim)
	copy(data[1:8], cmp.Int64ToBytes(e.xPos))
	copy(data[9:16], cmp.Int64ToBytes(e.yPos))
	data[17] = byte(e.orientation)
	data[18] = byte(e.neworientation)
	data[19] = cmp.BoolToByte(e.isMoving)
	data[20] = cmp.BoolToByte(e.keepMoving)
	copy(data[21:22], cmp.Int16ToBytes(int16(e.movingFrames)))
	copy(data[23:24], cmp.Int16ToBytes(int16(e.movedFrames)))
	copy(data[25:32], cmp.Float64ToBytes(e.movingStepSize))
	data[33] = cmp.BoolToByte(e.changed)
	data[34] = cmp.BoolToByte(e.isDirty)
	copy(data[35:36], cmp.Int16ToBytes(e.factoryCreationId))
	return
}

//Copys the Entity
func (e *Entity) Copy() (e2 *Entity) {
	e2 = &Entity{*e.WObj.Copy(), nil, e.currentAnim, e.xPos, e.yPos, e.orientation, e.neworientation, e.isMoving,
		e.keepMoving, e.movingFrames, e.movedFrames, e.movingStepSize, e.changed, e.isDirty, e.factoryCreationId, e.frame, e.Updater}
	e2.anims = make([]*GE.DayNightAnim, len(e.anims))
	for i, anim := range e.anims {
		if anim != nil {
			e2.anims[i] = anim.Copy()
		}
	}
	return
}

//Updates the movement and calls the provided Update func afterwards
func (e *Entity) UpdateAll(w *World) {
	if e.isMoving {
		if e.movedFrames >= e.movingFrames {
			e.isMoving = false
			e.orientation = e.neworientation
			if e.keepMoving {
				e.isMoving = true
				e.movedFrames = 0
				e.moveInDirection(e.orientation)
				e.movedFrames++
			} else {
				e.changed = true
			}
		} else {
			e.moveInDirection(e.orientation)
			e.movedFrames++
		}
		e.UpdateOrientationAnim()
	}
	if e.Updater != nil {
		e.Updater(e, w)
	}
}

//Returns the Bounds of the Entity
func (e *Entity) Bounds() (float64, float64) {
	return e.WObj.Bounds()
}

//Initiates a move action with a specific lenght an duration
func (e *Entity) Move(length, frames int) {
	if e.isMoving {
		return
	}
	//THIS SHOULD CAUSE A CHANGE
	e.isMoving = true
	e.movingFrames = frames
	e.movedFrames = 0
	e.movingStepSize = float64(length) / float64(frames)
	e.neworientation = e.orientation
}

//Sets the middle of the Entity
func (e *Entity) SetMiddleTo(x, y float64) {
	e.WObj.SetPos(x, y)
	e.setIntPos()
}

//Sets the top left corner of the Entity
func (e *Entity) SetTopLeftTo(x, y float64) {
	e.WObj.SetToXY(x, y)
	e.setIntPos()
}
func (e *Entity) setIntPos() {
	xf, yf, _ := e.WObj.GetPos()
	x, y := int64(math.Round(xf-0.5)), int64(math.Round(yf-0.5))
	if x != e.xPos || y != e.yPos {
		e.xPos, e.yPos = x, y
		e.changed = true
	}
}

//Changes the orientation
func (e *Entity) ChangeOrientation(newO uint8) {
	if newO != e.orientation {
		if e.isMoving {
			e.neworientation = newO
		} else {
			e.orientation = newO
		}
	}
}

//Updates the Orientation animation, ONLY call this if really necassary
func (e *Entity) UpdateOrientationAnim() {
	idx := e.orientation
	if e.isMoving {
		idx += 4
	}
	e.SetAnim(int(idx))
}
func (e *Entity) moveInDirection(dir uint8) {
	dx, dy := 0.0, 0.0
	switch dir {
	case 0:
		dx = -e.movingStepSize
		break
	case 1:
		dx = e.movingStepSize
		break
	case 2:
		dy = -e.movingStepSize
		break
	case 3:
		dy = e.movingStepSize
		break
	}
	e.WObj.MoveBy(dx, dy)
	e.setIntPos()
}

//Implements EntityI
func (e *Entity) GetDrawBox() *GE.Rectangle {
	return e.WObj.GetDrawBox()
}

//Implements EntityI
func (e *Entity) GetPos() (float64, float64, int8) {
	return e.WObj.GetPos()
}
func (e *Entity) FactoryCreationID() int16 {
	return e.factoryCreationId
}
func (e *Entity) GetWObj() *GE.WObj {
	return &e.WObj
}
func (e *Entity) Changed() bool {
	return e.changed
}
func (e *Entity) NotChangedAnymore() {
	e.changed = false
}
func (e *Entity) IsMoving() bool {
	return e.isMoving
}
func (e *Entity) KeepMoving(mv bool) {
	if mv != e.keepMoving {
		//THIS SHOULD CAUSE A CHANGE
		e.keepMoving = mv
	}
}
func (e *Entity) KeepsMoving() bool {
	return e.keepMoving
}
func (e *Entity) IntPos() (int64, int64) {
	return e.xPos, e.yPos
}
func (e *Entity) RegiserUpdateFunc(u EntityUpdater) {
	e.Updater = u
}
func (e *Entity) SetAnim(idx int) {
	if e.currentAnim == uint8(idx) || idx < 0 || idx >= len(e.anims) || e.anims[idx] == nil {
		return
	}
	e.currentAnim = uint8(idx)
	e.WObj.SetAnim(e.anims[idx])
	e.changed = true
}
func (e *Entity) GetAnim() uint8 {
	return e.currentAnim
}

//Implements EntityI
func (e *Entity) Draw(screen *ebiten.Image, lv int16, leftTopX, leftTopY, xStart, yStart, sqSize float64) {
	e.WObj.Update(*e.frame)
	e.WObj.Draw(screen, lv, leftTopX, leftTopY, xStart, yStart, sqSize)
}

func (e *Entity) Print() (out string) {
	out = fmt.Sprintf("Entity: FcID: %v, X: %v, Y: %v, CurrO: %v, NextO: %v, moves: %v, keepsMoving: %v, isDirty: %v, WObj: %v", e.factoryCreationId, e.xPos, e.yPos, e.orientation, e.neworientation, e.isMoving, e.keepMoving, e.changed, e.GetWObj())
	return
}

/**
Loads a Entity from a directory that contains DayNightAnims with the names listed in CREATURE_ANIM_NAMES
Also a WOBJ.txt file is needed, that describes the Attributes of the Entity
**/
func LoadEntity(path string, frameCounter *int) (*Entity, error) {
	if path[len(path)-1:] != "/" {
		path += "/"
	}
	pathS := strings.Split(path, "/")
	name := pathS[len(pathS)-2]
	e := &Entity{frame: frameCounter, anims: make([]*GE.DayNightAnim, 0), changed: true}
	wobj, err := GE.GetWObjFromPath(name, path+CREATURE_ANIM_STANDARD, path+CREATURE_WOBJ)
	if err != nil {
		return e, err
	}
	e.WObj = *wobj

	idx := &GE.List{}
	idx.LoadFromFile(path + INDEX_FILE_NAME)
	names := idx.GetSlice()
	e.anims = make([]*GE.DayNightAnim, 0)
	for _, anim_n := range names {
		anim, _ := GE.GetDayNightAnimFromParams(1, 1, 1, 1, path+anim_n+".txt", path+anim_n+".png")
		e.anims = append(e.anims, anim)
	}
	e.setIntPos()
	return e, nil
}
