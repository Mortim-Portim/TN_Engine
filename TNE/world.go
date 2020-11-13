package TNE

import (
	//"github.com/hajimehoshi/ebiten"
	"marvin/GraphEng/GE"
	"fmt"
)
/**
TODO
write a display-method for the world (sets tW, tH, moves with player)
add a player

syncronize players, entitys, lightlevel
**/
//path = "./res/wrld"
func GetWorld(X,Y,W,H float64, tW, tH, cW,cH int, cf *CreatureFactory, frameCounter *int, path, wrld_name, tile_F, struct_F string) (w *World) {
	if path[len(path)-1:] != "/" {
		path += "/"
	}
	w = &World{Path:path, FrameCounter:frameCounter, CF:cf}
	done := make(chan bool)
	go func() {
		w.ChunkMat = GE.GetMatrix(cW,cH,0)
		w.ChunkMat.InitIdx()
		w.Chunks = make([]*Chunk, cW*cH)
		for x := 0; x < w.ChunkMat.WAbs(); x++ {
			for y := 0; y < w.ChunkMat.HAbs(); y++ {
				idx, err := w.ChunkMat.Get(x,y)
				GE.ShitImDying(err)
				w.Chunks[idx] = GetChunk(x,y, cf, fmt.Sprintf("%stmp/%v-%v.chunk", path, x, y))
			}
		}
		done <- true
	}()
	go func() {
		wS, err := GE.LoadWorldStructure(X,Y,W,H, path+wrld_name+".map", tile_F, struct_F)
		GE.ShitImDying(err)
		wS.SetLightStats(10,255, 0.3)
		wS.SetLightLevel(15)
		wS.SetDisplayWH(tW, tH)
		w.Structure = wS
		
		done <- true
	}()
	
	<- done
	<- done
	return
}
type World struct {	
	Structure *GE.WorldStructure
	
	//Stores references to all chunks
	ChunkMat *GE.Matrix
	//Stores all chunks
	Chunks []*Chunk
	
	CF *CreatureFactory
	
	Path string
	FrameCounter *int
}

func (w *World) UpdateChunks(plXYChunkR ...[3]int) {
	allRems := make([]*chunkEntity, 0)
	done := make(chan bool)
	for _,posTs := range(plXYChunkR) {
		x,y := GetChunkOfTile(posTs[0], posTs[1])
		for _,delta := range(CHUNK_DELTAS[posTs[2]]) {
			go func() {
				idx, err := w.ChunkMat.Get(x+delta[0], y+delta[1])
				if err == nil {
					if w.Chunks[idx].LastUpdateFrame != *w.FrameCounter {
						w.Chunks[idx].LastUpdateFrame = *w.FrameCounter
						allRems = append(allRems, w.Chunks[idx].Update(w)...)
					}
				}
				done <- true
			}()
		}
	}
	for i := 0; i < len(CHUNK_DELTAS)*len(plXYChunkR); i++ {
		<- done
	}
	
}
func (w *World) AddDrawablesToWorldStructure(dws *GE.Drawables, x, y int) {
	done := make(chan bool)
	for _,delta := range(CHUNK_DELTAS[0]) {
		idx, err := w.ChunkMat.Get(x+delta[0], y+delta[1])
		if err == nil {
			go func() {
				w.Chunks[idx].AddToDrawables(dws)
				done <- true
			}()
		}
	}
	for i := 0; i < len(CHUNK_DELTAS[0]); i++ {
		<- done
	}
}
func (w *World) ReAssignEntities(ents []*chunkEntity) {
	for _,ent := range(ents) {
		x,y := ent.IntPos()
		cX,cY := GetChunkOfTile(int(x),int(y))
		idx,err := w.ChunkMat.Get(cX,cY)
		if err == nil {
			w.Chunks[idx].AddEntity(ent.EntityI)
		}
	}
}
