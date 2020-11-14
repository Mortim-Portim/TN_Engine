package TNE

import (
	"github.com/hajimehoshi/ebiten"
	"marvin/GraphEng/GE"
	"errors"
	"fmt"
)
var ERR_UNKNOWN_PLAYER = errors.New("Unknown player")
/**
TODO
Test entity movement

syncronize players, entitys, lightlevel
**/
//path = "./res/wrld"
func GetWorld(X,Y,W,H float64, tW, tH, cW,cH int, CF *CreatureFactory, frameCounter *int, path, wrld_name, tile_F, struct_F string) (w *World) {
	if path[len(path)-1:] != "/" {
		path += "/"
	}
	w = &World{Path:path, FrameCounter:frameCounter, CF:CF, ActivePlayerChunk:-1}
	w.Players = make([]*Player, 0)
	
	done := make(chan bool)
	go func() {
		w.ChunkMat = GE.GetMatrix(cW,cH,0)
		w.ChunkMat.InitIdx()
		w.Chunks = make([]*Chunk, cW*cH)
		for x := 0; x < w.ChunkMat.WAbs(); x++ {
			for y := 0; y < w.ChunkMat.HAbs(); y++ {
				idx, err := w.ChunkMat.Get(x,y)
				GE.ShitImDying(err)
				w.Chunks[idx] = GetChunk(x,y, CF, fmt.Sprintf("%stmp/%v-%v.chunk", path, x, y))
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
	
	//Stores all Players
	Players []*Player
	//Stores reference to the active player
	ActivePlayer, ActivePlayerChunk int
	//If the current player changes position or is replaced
	ActivePlayerChange bool
	
	CF *CreatureFactory
	
	Path string
	FrameCounter *int
}
func (w *World) Print() (out string) {
	out = fmt.Sprintf("Chunks: %v, Players: %v, ActivePlayer: %v, AP_Chunk: %v, AP_Change: %v, Path: %v, frame: %v\nCF:%s\n",
		len(w.Chunks), len(w.Players), w.ActivePlayer, w.ActivePlayerChunk, w.ActivePlayerChange, w.Path, *w.FrameCounter, w.CF.Print())
	out += fmt.Sprintf("ChunkMat:\n%s\n", w.ChunkMat.Print())
	out += w.Structure.Print()
	return
}

/**
Draws the surroundings of the active player
CALL on Client
**/
func (w *World) Draw(screen *ebiten.Image) {
	w.Structure.Draw(screen)
}
/**
Updates the world structures obj, if the player moved
Updates the lightlevel and possible changes of a lights position
CALL on Client
**/
func (w *World) UpdateWorldStructure() {
	if w.ActivePlayerChange {
		//fmt.Println("Updating ObjDrawables")
		w.Structure.UpdateObjDrawables()
	}
	//fmt.Println("Updating Lights")
	w.Structure.UpdateLightLevel(1)
	w.Structure.UpdateAllLightsIfNecassary()
}
/**
The active player moves the worldstructure if necassary
and if the active player changes the chunk the drawn entities are updated
CALL on Client
**/
func (w *World) UpdateDrawables() {
	//fmt.Println("Player moving world")
	activePlayer := w.Players[w.ActivePlayer]
	activePlayer.MoveWorld(w)
	if w.ActivePlayerChange {
		x,y := activePlayer.IntPos()
		cX,cY := GetChunkOfTile(int(x), int(y))
		idx,err := w.ChunkMat.Get(cX,cY)
		//fmt.Printf("Checking if player changed chunk: x:%v, y:%v, cx:%v, cy:%v, idx:%v\n", x,y,cX,cY,idx)
		if err == nil {
			if int(idx) != w.ActivePlayerChunk {
				//fmt.Printf("Changing from chunk %v to %v\n", w.ActivePlayerChunk, idx)
				w.Structure.Add_Drawables.Clear()
				w.Structure.AddDrawable(w.Players[w.ActivePlayer])
				w.AddEntitiesToDrawables(w.Structure.Add_Drawables, cX, cY)
				w.ActivePlayerChunk = int(idx)
			}
		}
	}
}
/**
Updates the active player
CALL on !!Client!! on !!every frame!!
**/
func (w *World) UpdateActivePlayer() {
	w.ActivePlayerChange = false
	activePlayer := w.Players[w.ActivePlayer]
	activePlayer.Update(w)
	if activePlayer.Changed() {
		fmt.Println("Player changed")
		w.ActivePlayerChange = true
		activePlayer.NotChangedAnymore()
	}
}
/**
Sets the player that the world is drawn for and that is Updated
CALL on !!Client!!
**/
func (w *World) SetActivePlayer(playerIdx int) error {
	//fmt.Printf("Setting new active player: %v\n", playerIdx)
	if playerIdx < 0 || playerIdx >= len(w.Players) {
		return ERR_UNKNOWN_PLAYER
	}
	w.Structure.Add_Drawables.Remove(w.Players[w.ActivePlayer])
	w.Structure.AddDrawable(w.Players[playerIdx])
	w.ActivePlayer = playerIdx
	w.ActivePlayerChange = true
	return nil
}
/**
SHOULD remove a player by his name
CALL on !!Server!!
**/
func (w *World) RemovePlayer(name string) {
	
}
/**
Adds a player
CALL on !!Server!!
**/
func (w *World) AddPlayer(p *Player) {
	w.Players = append(w.Players, p)
}
/**
Adds all entities of all chunks around the given chunk to drawables (including the given one)
CALL on !!Client!!
**/
func (w *World) AddEntitiesToDrawables(dws *GE.Drawables, x, y int) {
	for _,delta := range(CHUNK_DELTAS[0]) {
		idx, err := w.ChunkMat.Get(x+delta[0], y+delta[1])
		if err == nil {
			w.Chunks[idx].AddToDrawables(dws)
		}
	}
}
/**
Reassigns all entities in ents to the chunk that they fit in
This should only be called when the entities move to a different chunk
CALL on !!Server!! on !!every frame!! for objs returned by !!UpdateChunks!!
**/
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
/**
Updates all chunks in a pattern given by CHUNK_DELTAS[plXYChunkR[2]] that lie around plXYChunkR[0], plXYChunkR[1]
CALL on !!Server!! on !!every frame!!
**/
func (w *World) UpdateChunks(plXYChunkR ...[3]int) (allRems []*chunkEntity) {
	allRems = make([]*chunkEntity, 0)
	done := make(chan bool)
	for _,posTs := range(plXYChunkR) {
		go func() {
			x,y := GetChunkOfTile(posTs[0], posTs[1])
			for _,delta := range(CHUNK_DELTAS[posTs[2]]) {
				idx, err := w.ChunkMat.Get(x+delta[0], y+delta[1])
				if err == nil {
					if w.Chunks[idx].LastUpdateFrame != *w.FrameCounter {
						w.Chunks[idx].LastUpdateFrame = *w.FrameCounter
						allRems = append(allRems, w.Chunks[idx].Update(w)...)
					}
				}
			}
			done <- true
		}()
	}
	for i := 0; i < len(plXYChunkR); i++ {
		<- done
	}
	return
}