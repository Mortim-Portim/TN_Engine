package TNE

import (
	//"marvin/GraphEng/GE"
)
const CHUNK_SIZE = 16

type Chunk struct {
	pos [2]int16
	entities [][]EntityI
}


func (c *Chunk) GetData() (bs []byte) {
	
	return
}
func (c *Chunk) SetData(bs []byte) {
	
}

func ChunkCoord2DtoIdx(x, y int) byte {
	if x >= 16 || y >= 16 {
		panic("NEVER call ChunkCoord2DtoIdx with coords >= 16")
	}
	return byte(x+CHUNK_SIZE*y)
}
func IdxtoChunkCoord2D(idx byte) (x,y int) {
	csm1 := byte(CHUNK_SIZE -1)
	x = int(idx%CHUNK_SIZE)
	y = int((idx-(idx%csm1))/csm1)
	return
}