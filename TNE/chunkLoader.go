package TNE

import (
	"strings"
)

func init() {
	CHUNK_DELTAS = make([][][2]int, 5)
	CHUNK_DELTAS[0] = ChunkStringToDeltaFromMiddle(chunk1x1)
	CHUNK_DELTAS[1] = ChunkStringToDeltaFromMiddle(chunk3x3)
	CHUNK_DELTAS[2] = ChunkStringToDeltaFromMiddle(chunk5x5)
	CHUNK_DELTAS[3] = ChunkStringToDeltaFromMiddle(chunk7x7)
	CHUNK_DELTAS[4] = ChunkStringToDeltaFromMiddle(chunk9x9)
}

const (
//	1
	chunk1x1 = "1"
	
//	1 1 1+
//	1 1 1+
//	1 1 1+
	chunk3x3 = "1 1 1+1 1 1+1 1 1"
	
//	0 1 1 1 0+
//	1 1 1 1 1+
//	1 1 1 1 1+
//	1 1 1 1 1+
//	0 1 1 1 0
	chunk5x5 = "0 1 1 1 0+1 1 1 1 1+1 1 1 1 1+1 1 1 1 1+0 1 1 1 0"
	
//	0 0 1 1 1 0 0+
//	0 1 1 1 1 1 0+
//	1 1 1 1 1 1 1+
//	1 1 1 1 1 1 1+
//	1 1 1 1 1 1 1+
//	0 1 1 1 1 1 0+
//	0 0 1 1 1 0 0
	chunk7x7 = "0 0 1 1 1 0 0+0 1 1 1 1 1 0+1 1 1 1 1 1 1+1 1 1 1 1 1 1+1 1 1 1 1 1 1+0 1 1 1 1 1 0+0 0 1 1 1 0 0"
	
//	0 0 0 1 1 1 0 0 0+
//	0 0 1 1 1 1 1 0 0+
//	0 1 1 1 1 1 1 1 0+
//	1 1 1 1 1 1 1 1 1+
//	1 1 1 1 1 1 1 1 1+
//	1 1 1 1 1 1 1 1 1+
//	0 1 1 1 1 1 1 1 0+
//	0 0 1 1 1 1 1 0 0+
//	0 0 0 1 1 1 0 0 0
	chunk9x9 = "0 0 0 1 1 1 0 0 0+0 0 1 1 1 1 1 0 0+0 1 1 1 1 1 1 1 0+1 1 1 1 1 1 1 1 1+1 1 1 1 1 1 1 1 1+1 1 1 1 1 1 1 1 1+0 1 1 1 1 1 1 1 0+0 0 1 1 1 1 1 0 0+0 0 0 1 1 1 0 0 0"
)

var CHUNK_DELTAS [][][2]int

func ChunkStringToDeltaFromMiddle(cs string) (d [][2]int) {
	cs = strings.ReplaceAll(cs, " ", "")
	lines := strings.Split(cs, "+")
	middle := (len(lines)-1)/2
	d = make([][2]int, 0)
	
	for y,line := range(lines) {
		for x,r := range(line) {
			if string(r) == "1" {
				d = append(d, [2]int{x-middle, y-middle})
			}
		}
	}
	return
}
func GetChunkOfTile(x,y int) (xC, yC int) {
	xWithoutDx := float64(x-x%int(CHUNK_SIZE)); yWithoutDy := float64(y-y%int(CHUNK_SIZE))
	tilesDX := xWithoutDx/CHUNK_SIZE; tilesDY := yWithoutDy/CHUNK_SIZE
	return int(tilesDX), int(tilesDY)
}