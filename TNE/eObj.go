package TNE

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/mortim-portim/GraphEng/GE"

	"errors"
	"fmt"
	"math"
	"strings"

	cmp "github.com/mortim-portim/GraphEng/Compression"
)

const (
	OBJ_ANIM_STANDARD  = "idle_L"
)
const CREATURE_WOBJ = "#WOBJ"
var ERR_WRONG_BYTE_LENGTH = errors.New("Wrong byte length")
var ERR_UNKNOWN_ACTION = errors.New("Unknown Action")

type Eobj struct {
	*GE.WObj

	anims       []*GE.DayNightAnim
	currentAnim uint8

	xPos, yPos                  int64
	orientation, neworientation *Direction
	isMoving, keepMoving        bool

	movingFrames, movedFrames int
	movingStepSize            float64

	changed bool

	factoryCreationId int16
	
	frame   *int
	UpdateFunc func(eo *Eobj, world *World)
}

const OBJ_CREATION_DATA_LENGTH = 18
func (e *Eobj) GetCreationData() (bs []byte) {
	bs = make([]byte, OBJ_CREATION_DATA_LENGTH)
	copy(bs[0:8], cmp.Int64ToBytes(e.xPos))
	copy(bs[8:16], cmp.Int64ToBytes(e.yPos))
	copy(bs[16:18], cmp.Int16ToBytes(e.factoryCreationId))
	return
}
//Copys the Eobj
func (e *Eobj) Copy() (e2 *Eobj) {
	e2 = &Eobj{e.WObj.Copy(), nil, e.currentAnim, e.xPos, e.yPos, e.orientation.Copy(), e.neworientation.Copy(), e.isMoving, e.keepMoving, 
		e.movingFrames, e.movedFrames, e.movingStepSize, e.changed, e.factoryCreationId, e.frame, e.UpdateFunc}
	e2.anims = make([]*GE.DayNightAnim, len(e.anims))
	for i, anim := range e.anims {
		if anim != nil {
			e2.anims[i] = anim.Copy()
		}
	}
	return
}

//Updates the movement and calls the provided Update func afterwards
func (e *Eobj) UpdateAll(w *World, Collider func(x,y int)bool) {
	if e.isMoving {
		if e.movedFrames >= e.movingFrames {
			e.isMoving = false
			e.orientation = e.neworientation
			if e.keepMoving {
				e.isMoving = true
				e.movedFrames = 0
				e.moveInDirection(e.orientation, Collider)
				e.movedFrames++
			} else {
				e.changed = true
			}
		} else {
			e.moveInDirection(e.orientation, Collider)
			e.movedFrames++
		}
		e.UpdateOrientationAnim()
	}
	if e.UpdateFunc != nil {
		e.UpdateFunc(e, w)
	}
}

//Returns the Bounds of the Eobj
func (e *Eobj) Bounds() (float64, float64) {
	return e.WObj.Bounds()
}

//Initiates a move action with a specific lenght an duration
func (e *Eobj) Move(length float64, frames int) {
	if e.isMoving {
		return
	}
	//e.AppliedActions = append(e.AppliedActions, []byte{ENTITY_START_MOVE, byte(length), byte(frames)}...)
	e.isMoving = true
	e.movingFrames = frames
	e.movedFrames = 0
	e.movingStepSize = length / float64(frames)
	e.neworientation = e.orientation
}

//Sets the middle of the Eobj
func (e *Eobj) SetMiddleTo(x, y float64) {
	e.WObj.SetMiddle(x, y)
	e.setIntPos()
}
//Sets the top left corner of the Eobj
func (e *Eobj) SetTopLeftTo(x, y float64) {
	e.WObj.SetTopLeft(x, y)
	e.setIntPos()
}
//Sets the top left corner of the Eobj
func (e *Eobj) SetBottomRightTo(x, y float64) {
	e.WObj.SetBottomRight(x, y)
	e.setIntPos()
}
func (e *Eobj) setIntPos() {
	xf, yf, _ := e.WObj.GetMiddle()
	x, y := e.FloatPosToIntPos(xf,yf)
	if int64(x) != e.xPos || int64(y) != e.yPos {
		e.xPos, e.yPos = int64(x), int64(y)
		e.changed = true
	}
}

//Changes the orientation
func (e *Eobj) ChangeOrientation(dir *Direction) {
	if dir.IsValid() {
		if e.isMoving {
			e.neworientation = dir
		}else{
			if !dir.Equals(e.orientation) {
				e.orientation = dir
				e.neworientation = dir
			}
		}
	}
}

//Updates the Orientation animation, ONLY call this if really necassary
func (e *Eobj) UpdateOrientationAnim() {
	idx := e.orientation.ID
	if idx >= 0 {
		if idx == ENTITY_ORIENTATION_LU || idx == ENTITY_ORIENTATION_LD {
			idx = ENTITY_ORIENTATION_L
		}
		if idx == ENTITY_ORIENTATION_RU || idx == ENTITY_ORIENTATION_RD {
			idx = ENTITY_ORIENTATION_R
		}
		if e.isMoving {
			idx += 4
		}
		e.SetAnim(int(idx))
	}
}
func (e *Eobj) moveInDirection(dir *Direction, Collider func(x,y int)bool) {
	if dir.IsValid() {
		dx, dy := e.movingStepSize, e.movingStepSize
		dx *= dir.Dx; dy *= dir.Dy
		fx,fy,_ := e.GetPos()
		nX,nY := e.FloatPosToIntPos(fx+dx, fy+dy)
		if Collider(nX, nY) {
			if !Collider(e.FloatPosToIntPos(fx+dx, fy)) {
				e.WObj.MoveBy(dx, 0)
			}else if !Collider(e.FloatPosToIntPos(fx, fy+dy)) {
				e.WObj.MoveBy(0, dy)
			}else{
				return
			}
		}else{
			e.WObj.MoveBy(dx, dy)
		}
		e.setIntPos()
	}
}
func (e *Eobj) FloatPosToIntPos(fx, fy float64) (int, int) {
	return int(math.Round(fx-0.5)), int(math.Round(fy-0.5))
}
//Implements EntityI
func (e *Eobj) GetDrawBox() *GE.Rectangle {
	return e.WObj.GetDrawBox()
}

//Implements EntityI
func (e *Eobj) GetPos() (float64, float64, int8) {
	return e.WObj.GetMiddle()
}
func (e *Eobj) GetTopLeft() (float64, float64) {
	return e.WObj.GetTopLeft()
}
func (e *Eobj) GetBottomRight() (float64, float64) {
	return e.WObj.GetBottomRight()
}
func (e *Eobj) FactoryCreationID() int16 {
	return e.factoryCreationId
}
func (e *Eobj) GetWObj() *GE.WObj {
	return e.WObj
}
func (e *Eobj) Changed() bool {
	return e.changed
}
func (e *Eobj) NotChangedAnymore() {
	e.changed = false
}
func (e *Eobj) IsMoving() bool {
	return e.isMoving
}
func (e *Eobj) KeepMoving(mv bool) {
	if mv != e.keepMoving {
//		ac := ENTITY_KEEP_MOVING
//		if !mv {
//			ac = ENTITY_STOP_KEEP_MOVING
//		}
//		e.AppliedActions = append(e.AppliedActions, ac)
		e.keepMoving = mv
	}
}
func (e *Eobj) KeepsMoving() bool {
	return e.keepMoving
}
func (e *Eobj) IntPos() (int64, int64) {
	return e.xPos, e.yPos
}
//[6]byte
func (e *Eobj) PosToBytes() []byte {
	x,y,dx,dy := e.GetPosIntPBytes()
	return append(append(cmp.Int16ToBytes(int16(x)), cmp.Int16ToBytes(int16(y))...), byte(dx), byte(dy))
}
//[6]byte
func (e *Eobj) PosFromBytes(bs []byte) {
	x := cmp.BytesToInt16(bs[0:2]); y := cmp.BytesToInt16(bs[2:4])
	e.SetPosIntPBytes(int(x),int(y), bs[4], bs[5])
}
func (e *Eobj) GetPosIntPBytes() (int, int, byte, byte) {
	fx,fy := e.GetBottomRight()
	x := math.Floor(fx); y := math.Floor(fy)
	dx := fx-x; dy := fy-y
	return int(x), int(y), byte(dx*255), byte(dy*255)
}
func (e *Eobj) SetPosIntPBytes(x,y int, bdx, bdy byte) {
	dx := float64(bdx)/255;dy := float64(bdy)/255
	e.SetBottomRightTo(float64(x)+dx, float64(y)+dy)
}
func (e *Eobj) RegisterUpdateFunc(u func(eo *Eobj, world *World)) {
	e.UpdateFunc = u
}
func (e *Eobj) SetAnim(idx int) {
	if e.currentAnim == uint8(idx) || idx < 0 || idx >= len(e.anims) || e.anims[idx] == nil {
		return
	}
	e.currentAnim = uint8(idx)
	e.WObj.SetAnim(e.anims[idx])
	e.changed = true
}
func (e *Eobj) GetAnim() uint8 {
	return e.currentAnim
}

//Implements EntityI
func (e *Eobj) Draw(screen *ebiten.Image, lv int16, leftTopX, leftTopY, xStart, yStart, sqSize float64) {
	e.WObj.Update(*e.frame)
	e.WObj.Draw(screen, lv, leftTopX, leftTopY, xStart, yStart, sqSize)
}

func (e *Eobj) Print() (out string) {
	out = fmt.Sprintf("Eobj: FcID: %v, X: %v, Y: %v, CurrO: %v, NextO: %v, moves: %v, keepsMoving: %v, isDirty: %v, WObj: %v", e.factoryCreationId, e.xPos, e.yPos, e.orientation, e.neworientation, e.isMoving, e.keepMoving, e.changed, e.GetWObj())
	return
}

/**
Loads a Eobj from a directory that contains DayNightAnims with the names listed in CREATURE_ANIM_NAMES
Also a WOBJ.txt file is needed, that describes the Attributes of the Eobj
**/
func LoadEobj(path string, frameCounter *int) (*Eobj, error) {
	if path[len(path)-1:] != "/" {
		path += "/"
	}
	pathS := strings.Split(path, "/")
	name := pathS[len(pathS)-2]
	e := &Eobj{frame: frameCounter, anims: make([]*GE.DayNightAnim, 0), changed: true}
	wobj, err := GE.GetWObjFromPath(name, path+OBJ_ANIM_STANDARD, path+CREATURE_WOBJ)
	if err != nil {
		return e, err
	}
	e.WObj = wobj

	idx := &GE.List{}
	idx.LoadFromFile(path + INDEX_FILE_NAME)
	names := idx.GetSlice()
	e.anims = make([]*GE.DayNightAnim, 0)
	for _, anim_n := range names {
		anim, _ := GE.GetDayNightAnimFromParams(1, 1, 1, 1, path+anim_n+".txt", path+anim_n+".png")
		e.anims = append(e.anims, anim)
	}
	e.setIntPos()
	e.orientation = GetNewDirection()
	e.neworientation = GetNewDirection()
	//e.ResetAppliedActions()
	return e, nil
}