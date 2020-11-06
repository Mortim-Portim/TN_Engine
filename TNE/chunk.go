package TNE

import (
	//"marvin/GraphEng/GE"
)
const CHUNK_SIZE = 32

type Chunk struct {
	x1,y1,x2,y2 int
	entities []EntityI
}