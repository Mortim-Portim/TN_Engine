package TNE

import (
	"fmt"

	cmp "github.com/mortim-portim/GraphEng/compression"
)

const (
	SyncAction_StartMove = byte(iota)
	SyncAction_Orientation
	SyncAction_NextOrientation
	SyncAction_KeepMoving
	SyncAction_ManualAnimationChange
	SyncAction_Position
)

func NewActionStack(data ...byte) *ActionStack {
	if len(data) == 0 {
		data = []byte{}
	}
	as := &ActionStack{data, false}
	return as
}

type ActionStack struct {
	data        []byte
	ManualReset bool
}

func (as *ActionStack) Print() (out string) {
	as.iterate(func(t byte, data []byte) int {
		switch t {
		case SyncAction_StartMove:
			out += "|StartMove"
			return 10
		case SyncAction_Orientation:
			out += "|Orientation"
			return 1
		case SyncAction_NextOrientation:
			out += "|NextOrientation"
			return 1
		case SyncAction_KeepMoving:
			out += "|KeepMoving"
			return 1
		case SyncAction_ManualAnimationChange:
			out += "|ManualAnimationChange"
			return 1
		case SyncAction_Position:
			out += "|Position"
			return 6
		}
		return 0
	})
	out += fmt.Sprintf("|")
	return
}
func (as *ActionStack) Copy() (as2 *ActionStack) {
	as2 = &ActionStack{make([]byte, len(as.data)), as.ManualReset}
	copy(as2.data, as.data)
	return
}
func (as *ActionStack) Reset() {
	as.data = []byte{}
}
func (as *ActionStack) ApplyOnEobj(e *Entity) {
	as.iterate(func(t byte, data []byte) int {
		switch t {
		case SyncAction_StartMove:
			length := cmp.BytesToFloat64(data[0:8])
			frames := int(cmp.BytesToUInt16(data[8:10]))
			e.isMoving = true
			e.movingFrames = frames
			e.movedFrames = 0
			e.movingStepSize = length / float64(frames)
			e.neworientation = e.orientation
			return 10
		case SyncAction_Orientation:
			e.orientation.FromByte(data[0])
			return 1
		case SyncAction_NextOrientation:
			e.neworientation.FromByte(data[0])
			return 1
		case SyncAction_KeepMoving:
			e.keepMoving = cmp.BytesToBools(data[0:1])[0]
			return 1
		case SyncAction_ManualAnimationChange:
			e.setAnim(uint8(data[0]))
			return 1
		case SyncAction_Position:
			e.PosFromBytes(data)
			return 6
		}
		return 0
	})
	//as.WaitForReset()
}
func (as *ActionStack) AppendAndApply(bs []byte, e *Entity) {
	old_data := as.data
	as.SetAll(bs)
	as.ApplyOnEobj(e)
	as.SetAll(append(old_data, bs...))
}
func (as *ActionStack) SetAll(bs []byte) {
	as.data = bs
}
func (as *ActionStack) GetAll() (bs []byte) {
	bs = as.data
	if !as.ManualReset {
		as.Reset()
	}
	return
}
func (as *ActionStack) iterate(fnc func(byte, []byte) int) {
	data := as.data
	for len(data) > 0 {
		t := data[0]
		l := fnc(t, data[1:])
		data = data[l+1:]
	}
}
func (as *ActionStack) Add(t byte, data ...byte) {
	as.data = append(as.data, t)
	as.data = append(as.data, data...)
}
func (as *ActionStack) AddStartMove(length float64, frames int) {
	as.Add(SyncAction_StartMove, append(cmp.Float64ToBytes(length), cmp.UInt16ToBytes(uint16(frames))...)...)
}
func (as *ActionStack) AddOrientation(o *Direction) {
	as.Add(SyncAction_Orientation, o.ToByte())
}
func (as *ActionStack) AddNextOrientation(o *Direction) {
	as.Add(SyncAction_NextOrientation, o.ToByte())
}
func (as *ActionStack) AddKeepMoving(b bool) {
	as.Add(SyncAction_KeepMoving, cmp.BoolsToBytes(b)...)
}
func (as *ActionStack) AddManualAnimationChange(idx uint8) {
	as.Add(SyncAction_ManualAnimationChange, byte(idx))
}
func (as *ActionStack) AddPosition(bs []byte) {
	as.Add(SyncAction_Position, bs...)
}
