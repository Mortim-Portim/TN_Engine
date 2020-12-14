package TNE

import (
	"errors"
	"fmt"
	"time"
	"github.com/hajimehoshi/ebiten"
	"github.com/mortim-portim/GraphEng/GE"
	"github.com/mortim-portim/GameConn/GC"
)

var ERR_UNKNOWN_PLAYER = errors.New("Unknown player")

//X,Y,W,H float64, tW, tH, cW,cH, ChunkUpdateRange int, CF *EntityFactory, frameCounter *int, path, wrld_name, tile_F, struct_F string
func GetWorld(X, Y, W, H float64, tW, tH, cW, cH, ChunkUpdateRange int, CF *EntityFactory, frameCounter *int, path, wrld_name, tile_F, struct_F string) (w *World) {
	if path[len(path)-1:] != "/" {
		path += "/"
	}
	w = &World{Path: path, FrameCounter: frameCounter, CF: CF, ActivePlayerChunk: -1}
	w.Players = make([]*Player, 0)

	done := make(chan bool)
	go func() {
		w.ChunkMat = GE.GetMatrix(cW, cH, 0)
		w.ChunkMat.InitIdx()
		w.Chunks = make([]*Chunk, cW*cH)
		for x := 0; x < w.ChunkMat.WAbs(); x++ {
			for y := 0; y < w.ChunkMat.HAbs(); y++ {
				idx, err := w.ChunkMat.Get(x, y)
				GE.ShitImDying(err)
				w.Chunks[idx] = GetChunk(x, y, CF, fmt.Sprintf("%stmp/%v-%v.chunk", path, x, y))
			}
		}
		done <- true
	}()
	go func() {
		wS, err := GE.LoadWorldStructure(X, Y, W, H, path+wrld_name+".map", tile_F, struct_F)
		GE.ShitImDying(err)
		wS.SetLightStats(10, 255, 0.3)
		wS.SetLightLevel(15)
		wS.SetDisplayWH(tW, tH)
		w.Structure = wS
		
		w.ChunkUpdateRange =   GC.CreateSyncByte(byte(ChunkUpdateRange))
		w.Actions =			   GC.CreateSyncString("")
		w.OtherPlayerChanges = GC.CreateSyncString("")
		w.LocalPlayerChanges = GC.CreateSyncString("")
		w.SyncedFrame =		   GC.CreateSyncInt64(0)
		w.SyncVarsHaveUpdate = make(chan bool)
		
		done <- true
	}()
	<-done
	<-done
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

	CF *EntityFactory

	Path         string
	FrameCounter *int
	
	SyncVarsHaveUpdate chan bool
	LastSyncedTime time.Time
	LastSyncedFrame int
	Ping time.Duration
	
	SyncedFrame *GC.SyncInt64
	ChunkUpdateRange *GC.SyncByte
	Actions, OtherPlayerChanges, LocalPlayerChanges *GC.SyncString
}

func (w *World) Print() (out string) {
	out = fmt.Sprintf("Chunks: %v, Players: %v, ActivePlayer: %v, AP_Chunk: %v, AP_Change: %v, Path: %v, frame: %v\nCF:%s\n",
		len(w.Chunks), len(w.Players), w.ActivePlayer, w.ActivePlayerChunk, w.ActivePlayerChange, w.Path, *w.FrameCounter, w.CF.Print())
	out += fmt.Sprintf("ChunkMat:\n%s\n", w.ChunkMat.Print())
	out += w.Structure.Print()
	return
}

/**
  _
 /  | o  _  ._ _|_
 \_ | | (/_ | | |_

**/
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
Draws the surroundings of the active player
CALL on !!Client!! on !!every frame!!
**/
func (w *World) Draw(screen *ebiten.Image) {
	w.Structure.Draw(screen)
}

//-------------------------------------------------------------------------------------------------------------------------------------
/**
Updates the active player
CALL on !!Client!! on !!every frame!!
**/
func (w *World) UpdateActivePlayer() {
	w.ActivePlayerChange = false
	activePlayer := w.Players[w.ActivePlayer]
	activePlayer.Update(w)
	if activePlayer.Changed() {
		w.ActivePlayerChange = true
		activePlayer.NotChangedAnymore()
	}
}

/**
The active player moves the worldstructure if necassary
and if the active player changes the chunk the drawn entities are updated
CALL on !!Client!! on !!every frame!!
**/
func (w *World) UpdateDrawables() {
	activePlayer := w.Players[w.ActivePlayer]
	activePlayer.MoveWorld(w)
	if w.ActivePlayerChange {
		x, y := activePlayer.IntPos()
		cX, cY := GetChunkOfTile(int(x), int(y))
		idx, err := w.ChunkMat.Get(cX, cY)
		if err == nil {
			if int(idx) != w.ActivePlayerChunk {
				w.Structure.Add_Drawables.Clear()
				w.Structure.AddDrawable(w.Players[w.ActivePlayer])
				w.AddEntitiesToDrawables(w.Structure.Add_Drawables, cX, cY)
				w.ActivePlayerChunk = int(idx)
			}
		}
	}
}
/**
Updates the world structures obj, if the player moved
CALL on !!Client!! on !!every frame!!
**/
func (w *World) UpdateWorldStructure() {
	if w.ActivePlayerChange {
		w.Structure.UpdateObjDrawables()
	}
}

//-------------------------------------------------------------------------------------------------------------------------------------
/**
  __
 (_   _  ._    _  ._
 __) (/_ | \/ (/_ |

**/
/**
Adds a entity to a chunk given by the coords cX, cY
CALL on !!Server!!
**/
func (w *World) AddEntity(cX, cY int, e *Entity) error {
	idx, err := w.ChunkMat.Get(cX, cY)
	if err != nil {
		return err
	}
	return w.Chunks[idx].AddEntity(e)
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

//-------------------------------------------------------------------------------------------------------------------------------------
/**
Updates the lightlevel and applies raycasting if necassary
CALL on !!Server!! on !!every frame!!
**/
func (w *World) UpdateLights() {
	w.Structure.UpdateLightLevel(1)
	w.Structure.UpdateAllLightsIfNecassary()
}
/**
Downdates the lightlevel and applies raycasting if necassary
CALL on !!Server!! on !!every frame!!
**/
func (w *World) DowndateLights(count int) {
	w.Structure.DowndateLightLevel(1, count)
	w.Structure.UpdateAllLightsIfNecassary()
}

/**
Updates all chunks around all players with the specified delta of the world
calls UpdateChunks and ReAssignEntities
CALL on !!Server!! on !!every frame!!
**/
func (w *World) UpdatePlayerChunks(Players []*Player) {
	pos := make([][2]int, len(Players))
	for i, player := range Players {
		x, y := player.IntPos()
		pos[i] = [2]int{int(x), int(y)}
	}
	changedEntities := w.UpdateChunks(int(w.ChunkUpdateRange.GetByte()), pos...)
	w.ReAssignEntities(changedEntities)
}

//-------------------------------------------------------------------------------------------------------------------------------------
/**
Adds all entities of all chunks around the given chunk to drawables (including the given one)
called by !!UpdateDrawables!!
**/
func (w *World) AddEntitiesToDrawables(dws *GE.Drawables, x, y int) {
	for _, delta := range CHUNK_DELTAS[1] {
		idx, err := w.ChunkMat.Get(x+delta[0], y+delta[1])
		if err == nil {
			w.Chunks[idx].AddToDrawables(dws)
		}
	}
}

/**
Reassigns all entities in ents to the chunk that they fit in
This should only be called when the entities move to a different chunk
called by !!UpdatePlayerChunks!!
**/
func (w *World) ReAssignEntities(ents []*Entity) {
	for _, ent := range ents {
		x, y := ent.IntPos()
		cX, cY := GetChunkOfTile(int(x), int(y))
		idx, err := w.ChunkMat.Get(cX, cY)
		if err == nil {
			w.Chunks[idx].AddEntityLocal(ent)
		}
	}
}

/**
Updates all chunks in a pattern given by CHUNK_DELTAS[chunkRange] that lie around plXYChunkR[0], plXYChunkR[1]
called by !!UpdatePlayerChunks!!
**/
func (w *World) UpdateChunks(chunkRange int, plXYChunkR ...[2]int) (allRems []*Entity) {
	allRems = make([]*Entity, 0)
	done := make(chan bool)
	for _, posTs := range plXYChunkR {
		go func() {
			x, y := GetChunkOfTile(posTs[0], posTs[1])
			for _, delta := range CHUNK_DELTAS[chunkRange] {
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
		<-done
	}
	return
}