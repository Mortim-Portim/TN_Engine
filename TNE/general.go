package TNE

import (
	"fmt"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/mortim-portim/GraphEng/GE"
)

var FPS = 30.0

func OnRectWithWorldStructObjCollision(rect *GE.Rectangle, ws *GE.WorldStructure, onColl func(so *GE.StructureObj, ent *Entity, ply *Player)) {
	var Hitbox *GE.Rectangle
	var so *GE.StructureObj
	var ent *Entity
	var ply *Player
	ws.IterateOverCollidablesInRect(rect.BoundingRect(ws.TileMat.Focus()), func(dw GE.Drawable) {
		Hitbox = nil
		so = nil
		ent = nil
		ply = nil
		SO, ok := dw.(*GE.StructureObj)
		if ok && SO.Collides {
			Hitbox = SO.Hitbox
			so = SO
		} else {
			ENT, ok := dw.(*Entity)
			if ok && ENT.Collides() {
				Hitbox = ENT.Hitbox
				ent = ENT
			} else {
				PLY, ok := dw.(*Player)
				if ok && PLY.Collides() {
					Hitbox = PLY.Hitbox
					ply = PLY
				}
			}
		}
		if Hitbox != nil && Hitbox.Overlaps(rect) {
			onColl(so, ent, ply)
		}
	})
}

func CPUs() int {
	return runtime.NumCPU()
}
func PrintPerformance(frame, timeTaken int) (out string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	out = fmt.Sprintf("frame %v, TPS: %0.3f, Updating took: %v, Alloc: %v, TotalAlloc: %v, Sys: %v, NumGC: %v",
		frame, ebiten.CurrentTPS(), timeTaken, bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys), m.NumGC)
	return
}
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
func GetSliceOfVal(length, val int) (is []int) {
	is = make([]int, length)
	for i := range is {
		is[i] = val
	}
	return
}

type UniqueIDFactory []int

func (uf *UniqueIDFactory) GetID(min, max int) int {
	for id := min; id <= max; id++ {
		if containsI(*uf, id) == -1 {
			*uf = append(*uf, id)
			return id
		}
	}
	return -1
}
func (uf *UniqueIDFactory) AddID(id int) {
	if idx := containsI(*uf, id); idx >= 0 {
		(*uf)[idx] = (*uf)[len(*uf)-1]
		*uf = (*uf)[:len(*uf)-1]
	}
}

//Returns true if e is in s
func containsI(s []int, e int) int {
	for i, a := range s {
		if a == e {
			return i
		}
	}
	return -1
}
