package TNE

import (
	"fmt"
	"math"
	"time"

	"github.com/mortim-portim/GraphEng/GE"
)

/**
TODO
save world with worldstruct, entities, players
**/

type WorldParams struct {
	ChunkUpdateRange int
	Ef               *EntityFactory
	FrameCounter     *int
	Struct           *GE.WorldStructure
}

//
func GetWorld(wp *WorldParams, path string) (w *World) {
	if path[len(path)-1:] != "/" {
		path += "/"
	}
	w = &World{
		ChunkRange:   wp.ChunkUpdateRange,
		Path:         path,
		FrameCounter: wp.FrameCounter,
		Ef:           wp.Ef,
		Players:      make([]*Player, 0),
		Entities:     make([]*Entity, 0),
	}

	w.Structure = wp.Struct
	//	w.Structure.SetLightStats(10, 255, 0.3)
	//	w.Structure.SetLightLevel(15)
	//	w.Structure.SetDisplayWH(wp.TilesX, wp.TilesY)

	wX, wY := w.Structure.Size()
	cW := int(math.Ceil(float64(wX) / CHUNK_SIZE))
	cH := int(math.Ceil(float64(wY) / CHUNK_SIZE))

	w.ChunkMat = GE.GetMatrix(cW, cH, 0)
	w.ChunkMat.InitIdx()
	w.Chunks = make([]*Chunk, cW*cH)
	for x := 0; x < w.ChunkMat.WAbs(); x++ {
		for y := 0; y < w.ChunkMat.HAbs(); y++ {
			idx, err := w.ChunkMat.Get(x, y)
			GE.ShitImDying(err)
			w.Chunks[idx] = GetChunk(x, y)
		}
	}
	return
}

type World struct {
	Structure *GE.WorldStructure

	//Stores references to all chunks
	ChunkMat *GE.Matrix
	//Stores all chunks
	Chunks []*Chunk
	//Range of chunks around players
	ChunkRange int

	//Stores all Entities
	Entities []*Entity

	//Stores all Players
	Players []*Player

	Ef *EntityFactory

	ResetConfirm *chan bool

	Path         string
	FrameCounter *int
}

func (w *World) ResetActions() {
	for _, pl := range w.Players {
		pl.Actions().Reset()
	}
	for _, ent := range w.Entities {
		ent.Actions().Reset()
	}
}
func (w *World) Print(ents bool) (out string, c int) {
	out = fmt.Sprintf("%v: ", *w.FrameCounter)
	for _, pl := range w.Players {
		x, y, _ := pl.GetPos()
		out += fmt.Sprintf("(P)(%p)(%v)|%0.2f, %0.2f, %s|", pl, pl.ID, x, y, pl.Entity.Actions().Print())
		c++
	}
	if ents {
		for _, ent := range w.Entities {
			x, y, _ := ent.GetPos()
			out += fmt.Sprintf("(E)(%p)(%v)|%0.2f, %0.2f, %s|", ent, ent.ID, x, y, ent.Actions().Print())
			c++
		}
	}
	return
}
func (w *World) UpdateAllPos() {
	for _, pl := range w.Players {
		pl.AddPos()
	}
	for _, ent := range w.Entities {
		ent.AddPos()
	}
	return
}

/**
Updates the lightlevel and applies raycasting if necassary
**/
func (w *World) UpdateLights(t time.Duration) {
	w.Structure.UpdateTime(t)
	w.Structure.UpdateAllLightsIfNecassary()
}

/**
Updates all chunks around all players with the specified delta of the world
**/
func (w *World) UpdateAllPlayer() {
	for _, pl := range w.Players {
		pl.Update(w, true, w.Structure.Collides)
	}
}

/**
Updates all chunks around all players with the specified delta of the world
**/
func (w *World) UpdatePlayerChunks(Players ...*Player) []int {
	return w.UpdateChunks(w.GetPlayerChunks(Players...))
}

/**
Returns a list of indexes refering to the chunks around the players
**/
func (w *World) GetPlayerChunks(Players ...*Player) (idxs []int) {
	idxs = make([]int, 0)
	for _, player := range Players {
		cx, cy := GetChunkOfEntity(player.Entity)
		for _, delta := range CHUNK_DELTAS[w.ChunkRange] {
			idx, err := w.ChunkMat.Get(cx+delta[0], cy+delta[1])
			if err == nil && containsI(idxs, int(idx)) == -1 {
				idxs = append(idxs, int(idx))
			}
		}
	}
	return
}

/**
Updates the given chunks
**/
func (w *World) UpdateChunks(idxs []int) (chnged []int) {
	allRems := make([]*Entity, 0)
	chnged = make([]int, 0)
	for _, idx := range idxs {
		if w.Chunks[idx].LastUpdateFrame != *w.FrameCounter {
			w.Chunks[idx].LastUpdateFrame = *w.FrameCounter
			rems := w.Chunks[idx].UpdateEntities(w, false, w.Structure.Collides)
			if len(rems) > 0 {
				allRems = append(allRems, rems...)
				chnged = append(chnged, idx)
			}
		}
	}
	w.ReAssignEntities(allRems)
	return
}

/**
Reassigns all entities in ents to the chunk that they fit in
This should only be called when the entities move to a different chunk
**/
func (w *World) ReAssignEntities(ents []*Entity) {
	for _, ent := range ents {
		x, y := ent.IntPos()
		cX, cY := GetChunkOfTile(int(x), int(y))
		idx, err := w.ChunkMat.Get(cX, cY)
		if err == nil {
			w.Chunks[idx].Add(ent)
		}
	}
}

/**
Adds a entity to the chunk it belongs to
**/
func (w *World) AddEntity(e *Entity) error {
	cX, cY := GetChunkOfEntity(e)
	idx, err := w.ChunkMat.Get(cX, cY)
	if err != nil {
		return err
	}
	e.Actions().ManualReset = true
	w.Entities = append(w.Entities, e)
	return w.Chunks[idx].Add(e)
}
func (w *World) RemoveEntity(e *Entity) {
	idx := -1
	for i, e2 := range w.Entities {
		if e2 == e {
			idx = i
		}
	}
	if idx >= 0 {
		w.Entities[idx] = w.Entities[len(w.Entities)-1]
		w.Entities = w.Entities[:len(w.Entities)-1]
		cX, cY := GetChunkOfEntity(e)
		idx, err := w.ChunkMat.Get(cX, cY)
		if err == nil {
			w.Chunks[idx].RemoveEntity(e)
		}
	}
}

/**
Adds a player
**/
func (w *World) AddPlayer(p *Player) {
	idx := w.indexOfPlayer(p)
	if idx < 0 {
		p.Actions().ManualReset = true
		w.Players = append(w.Players, p)
	}
}

/**
Removes the player p if possible
**/
func (w *World) RemovePlayer(p *Player) {
	idx := w.indexOfPlayer(p)
	if idx >= 0 {
		w.Players[idx] = w.Players[len(w.Players)-1]
		w.Players = w.Players[:len(w.Players)-1]
	}
}
func (w *World) indexOfPlayer(p *Player) int {
	idx := -1
	for i, p2 := range w.Players {
		if p2 == p {
			idx = i
		}
	}
	return idx
}

//DEPRECATED
///**
//Sets the player that the world is drawn for and that is Updated
//**/
//func (w *World) SetActivePlayer(playerIdx int) error {
//	//fmt.Printf("Setting new active player: %v\n", playerIdx)
//	if playerIdx < 0 || playerIdx >= len(w.Players) {
//		return ERR_UNKNOWN_PLAYER
//	}
//	w.Structure.Add_Drawables.Remove(w.Players[w.ActivePlayer])
//	w.Structure.Add_Drawables.Add(w.Players[playerIdx])
//	w.ActivePlayer = playerIdx
//	w.ActivePlayerChange = true
//	return nil
//}
///**
//Draws the surroundings of the active player
//**/
//func (w *World) Draw(screen *ebiten.Image) {
//	w.Structure.Draw(screen)
//}
///**
//Updates the active player
//**/
//func (w *World) UpdateActivePlayer() {
//	w.ActivePlayerChange = false
//	activePlayer := w.Players[w.ActivePlayer]
//	activePlayer.Update(w)
//	if activePlayer.Changed() {
//		w.ActivePlayerChange = true
//		activePlayer.NotChangedAnymore()
//	}
//}
//
///**
//The active player moves the worldstructure if necassary
//and if the active player or a nearby entity changed the chunk the drawn entities are updated
//**/
//func (w *World) UpdateDrawables() {
//	activePlayer := w.Players[w.ActivePlayer]
//	activePlayer.MoveWorld(w.Structure)
//	if w.ActivePlayerChange || w.EntityChunkChange {
//		x, y := activePlayer.IntPos()
//		cX, cY := GetChunkOfTile(int(x), int(y))
//		idx, err := w.ChunkMat.Get(cX, cY)
//		if err == nil {
//			if int(idx) != w.ActivePlayerChunk || w.EntityChunkChange {
//				w.Structure.Add_Drawables.Clear()
//				w.Structure.AddDrawable(w.Players[w.ActivePlayer])
//				w.AddEntitiesToDrawables(w.Structure.Add_Drawables, cX, cY)
//				w.ActivePlayerChunk = int(idx)
//			}
//		}
//	}
//}
///**
//Updates the world structures obj, if the player moved
//**/
//func (w *World) UpdateWorldStructure() {
//	if w.ActivePlayerChange || w.EntityChunkChange {
//		w.Structure.UpdateObjDrawables()
//	}
//}
///**
//Adds all entities of all chunks around the given chunk to drawables (including the given one)
//called by !!UpdateDrawables!!
//**/
//func (w *World) AddEntitiesToDrawables(dws *GE.Drawables, x, y int) {
//	for _, delta := range CHUNK_DELTAS[1] {
//		idx, err := w.ChunkMat.Get(x+delta[0], y+delta[1])
//		if err == nil {
//			w.Chunks[idx].AddToDrawables(dws)
//		}
//	}
//}
