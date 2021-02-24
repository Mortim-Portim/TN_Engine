package TNE

import (
	"fmt"
	"runtime"

	"github.com/hajimehoshi/ebiten"
)

var FPS = 30.0

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
