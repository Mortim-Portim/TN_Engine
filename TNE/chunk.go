package TNE

import (
	"github.com/mortim-portim/GraphEng/GE"
	//"github.com/hajimehoshi/ebiten"
	"errors"

	//cmp "github.com/mortim-portim/GraphEng/Compression"
)
/**
TODO
delete chunkEntity
add an Add method that syncronizes entities using CreationData

implement a getDelta and setDelta method
**/

const CHUNK_SIZE = 16

var ERR_UNKNOWN_ENTITY_ID = errors.New("Unknown Entity ID")
var ERR_ENTITY_NOT_IN_THIS_CHUNK = errors.New("Entity not in this chunk")
var ERR_ENTITY_DOES_NOT_EXIST = errors.New("Entity does not exist")

//tmpPath is a path to a temporary file used for saving the chunk
func GetChunk(x, y int, cf *EntityFactory, tmpPath string) (c *Chunk) {
	c = &Chunk{pos: [2]int16{int16(x), int16(y)}, cf: cf, tmpPath: tmpPath}
	c.tileLT = [2]int16{CHUNK_SIZE * c.pos[0], CHUNK_SIZE * c.pos[1]}
	c.tileRB = [2]int16{c.tileLT[0] + CHUNK_SIZE, c.tileLT[1] + CHUNK_SIZE}
	c.entities = make([]*chunkEntity, 0)
	c.changes = make([]int, 0)
	return
}

type Chunk struct {
	pos, tileLT, tileRB            [2]int16
	entities                       []*chunkEntity
	changes 					   []int
	cf                             *EntityFactory
	tmpPath                        string
	LastUpdateFrame, LastDrawFrame int
}

/**
Adds an Entity to the chunk
called by world to move entities to this chunk
**/
func (c *Chunk) AddEntity(e *Entity) error {
	rx, ry, err := c.RelPosOfEntity(e)
	if err != nil {
		return err
	} else {
		newEnt := getNewChunkEntity(e, rx, ry)
		c.entities = append(c.entities, newEnt)
	}
	return nil
}
/**
Updates the chunk, returning removed entities
sets the removed entities to nil
RemoveNil should be called afterwards
**/
func (c *Chunk) Update(w *World) (removed []*chunkEntity) {
	removed = make([]*chunkEntity, 0)
	for idx, entity := range c.entities {
		if entity != nil {
			err := entity.Update(c, w)
			if err != nil {
				//Creature is not in this chunk anymore
				removed = append(removed, entity)
				c.entities[idx] = nil
			} else {
				if entity.Changed() {
					c.changes = append(c.changes, idx)
				}
			}
		}
	}
	return
}
/**
Removes all entities that are nil
**/
func (c *Chunk) RemoveNil() error {
	for idx, _ := range c.entities {
		if c.entities[idx] == nil {
			c.Remove(idx)
		}
	}
	return nil
}

/**
Removes a entity with speciefied index
**/
func (c *Chunk) Remove(i int) {
	c.entities[i] = c.entities[len(c.entities)-1]
	c.entities = c.entities[:len(c.entities)-1]
}

//returns changes as []byte
func (c *Chunk) GetDelta() (bs []byte) {
//	//[1]byte
//	bs = []byte{byte(len(c.removed))}
//	for _, rem := range c.removed {
//		//[3]byte
//		bs = append(bs, cmp.Int16ToBytes(int16(rem[0]))...)
//		bs = append(bs, byte(rem[1]))
//	}
//	for _, chng := range c.changes {
//		//[6]byte
//		bs = append(bs, cmp.Int16ToBytes(int16(chng[0]))...)
//		bs = append(bs, byte(chng[1]))
//		bs = append(bs, c.entities[chng[0]][chng[1]].changes...)
//	}
	c.changes = make([]int, 0)
	return
}

//sets changes
func (c *Chunk) SetDelta(bs []byte) {
	
//	removed = make([]*chunkEntity, 0)
//	rems := int(bs[0])
//	bs = bs[1:]
//	for i := 0; i < rems; i++ {
//		fcID := int(cmp.BytesToInt16(bs[0:1]))
//		idx := int(bs[2])
//		removed = append(removed, c.entities[fcID][idx])
//		c.Remove(fcID, idx)
//		bs = bs[3:]
//	}
//	for i := 0; i < len(bs)/6; i++ {
//		fcID := int(cmp.BytesToInt16(bs[0:1]))
//		idx := int(bs[2])
//		if idx >= len(c.entities[fcID]) {
//			c.entities[fcID] = append(c.entities[fcID], getNewChunkEntityFromBytes(bs[3:5], c.cf, fcID))
//		} else {
//			c.entities[fcID][idx].FromBytes(bs[3:5])
//		}
//		bs = bs[6:]
//	}
//	return
}

//Returns the relative position of a entity in a chunk
func (c *Chunk) RelPosOfEntity(e *Entity) (byte, byte, error) {
	eX, eY := e.IntPos()
	relX, relY := eX-int64(c.tileLT[0]), eY-int64(c.tileLT[1])
	if relX < 0 || relY < 0 || relX >= CHUNK_SIZE || relY >= CHUNK_SIZE {
		return 0, 0, ERR_ENTITY_NOT_IN_THIS_CHUNK
	}
	return byte(relX), byte(relY), nil
}

//Adds all entities of the chunk to drawables
func (c *Chunk) AddToDrawables(dws *GE.Drawables) {
	for _, ent := range c.entities {
		if ent != nil {
			dws.Add(ent.Entity)
		}
	}
}

//converts 2d coords in a chunk to a index
func ChunkCoord2DtoIdx(x, y int) byte {
	if x >= CHUNK_SIZE || y >= CHUNK_SIZE {
		panic("NEVER call ChunkCoord2DtoIdx with coords >= 16")
	}
	return byte(x + CHUNK_SIZE*y)
}

//converts a index in a chunk to 2d coords
func IdxtoChunkCoord2D(idx byte) (x, y int) {
	csm1 := byte(CHUNK_SIZE - 1)
	x = int(idx % CHUNK_SIZE)
	y = int((idx - (idx % csm1)) / csm1)
	return
}
/**
func getNewChunkEntityFromBytes(bs []byte, cf *EntityFactory, fcID int) (e *chunkEntity) {
	ent := cf.Get(fcID)
	e = &chunkEntity{ent, [2]byte{}, 0}
	e.FromBytes(bs)
	return e
}
**/
func getNewChunkEntity(e *Entity, rx, ry byte) *chunkEntity {
	return &chunkEntity{e, [2]byte{rx, ry}, ChunkCoord2DtoIdx(int(rx), int(ry))}
}

type chunkEntity struct {
	*Entity
	chunkPos    [2]byte
	chunkPosIdx byte
}

func (ce *chunkEntity) ApplyDelta(bs []byte) {
	ce.SetDelta(bs)
}
func (ce *chunkEntity) ReturnDelta() (bs []byte) {
	return ce.GetDelta()
}
func (ce *chunkEntity) Update(c *Chunk, w *World) error {
	ce.Entity.UpdateAll(w)
	rx, ry, err := c.RelPosOfEntity(ce.Entity)
	if err != nil {
		return err
	}
	ce.chunkPos[0] = rx
	ce.chunkPos[1] = ry
	ce.chunkPosIdx = ChunkCoord2DtoIdx(int(rx), int(ry))
	return nil
}

/**
//DEPRECATED
//Writes the chunk to the disk
func (c *Chunk) ToDisk() error {
	bs := make([]byte, 0)
	for fcID,l := range(c.entities) {
		bs = append(bs, cmp.Int16ToBytes(int16(fcID))...)
		for idx,entity := range(l) {
			bs = append(bs, append(entity.ToBytes(), byte(idx))...)
			entity.Entity = nil
		}
	}
	return ioutil.WriteFile(c.tmpPath, bs, 0644)
}
//Loads the chunk from the disk
func (c *Chunk) ToRAM() error {
	data, err := ioutil.ReadFile(c.tmpPath)
	if err != nil {return err}
	for _,l := range(c.entities) {
		fcID := int(cmp.BytesToInt16(data[0:1]))
		data = data[2:]
		for range(l) {
			idx := int(data[3])
			c.entities[fcID][idx].Entity = c.cf.Get(fcID)
			c.entities[fcID][idx].FromBytes(data[:3])
			data = data[4:]
		}
	}
	return nil
}
**/
