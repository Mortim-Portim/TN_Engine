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
