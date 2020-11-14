package TNE

import (
	"runtime"
	"github.com/hajimehoshi/ebiten"
	"fmt"
)

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
    return b / 1024 / 1024 / 8
}