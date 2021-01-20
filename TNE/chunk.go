package TNE


import (
	"github.com/mortim-portim/GraphEng/GE"
	//"github.com/hajimehoshi/ebiten"
	"errors"

	//cmp "github.com/mortim-portim/GraphEng/Compression"
)

const CHUNK_SIZE = 16

var ERR_UNKNOWN_ENTITY_ID = errors.New("Unknown Entity ID")
var ERR_ENTITY_NOT_IN_THIS_CHUNK = errors.New("Entity not in this chunk")
var ERR_ENTITY_DOES_NOT_EXIST = errors.New("Entity does not exist")

//tmpPath is a path to a temporary file used for saving the chunk
func GetChunk(x, y int) (c *Chunk) {
	c = &Chunk{pos: [2]int16{int16(x), int16(y)}}
	c.tileLT = [2]int16{CHUNK_SIZE * c.pos[0], CHUNK_SIZE * c.pos[1]}
	c.tileRB = [2]int16{c.tileLT[0] + CHUNK_SIZE, c.tileLT[1] + CHUNK_SIZE}
	c.entities = make([]*Entity, 0)
	return
}

type Chunk struct {
	pos, tileLT, tileRB            [2]int16
	entities                       []*Entity
	
	LastUpdateFrame int
}
func (c *Chunk) GetEntities() []*Entity {
	return c.entities
}
func (c *Chunk) Add(e *Entity) error {
	_, _, err := c.RelPosOfEntity(e)
	if err != nil {
		return err
	} else {
		c.entities = append(c.entities, e)
	}
	return nil
}
func (c *Chunk) RemoveEntity(e *Entity) {
	idx := -1
	for i, e2 := range c.entities {
		if e2 == e {idx = i}
	}
	if idx >= 0 {
		c.RemoveEntityByIdx(idx)
	}
}
func (c *Chunk) UpdateEntities(w *World, Collider func(x,y int)bool) (removed []*Entity) {
	removed = make([]*Entity, 0)
	for idx, entity := range c.entities {
		if entity != nil {
			entity.UpdateAll(w, Collider)
			_, _, err := c.RelPosOfEntity(entity)
			if err != nil {
				//Creature is not in this chunk anymore
				removed = append(removed, entity)
				c.entities[idx] = nil
			}
		}
	}
	return
}
func (c *Chunk) RemoveNilEntities() {
	rems := 0
	for idx, _ := range c.entities {
		if c.entities[idx-rems] == nil {
			c.RemoveEntityByIdx(idx-rems)
			rems ++
		}
	}
}
func (c *Chunk) RemoveEntityByIdx(i int) {
	c.entities[i] = c.entities[len(c.entities)-1]
	c.entities = c.entities[:len(c.entities)-1]
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
			dws.Add(ent)
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
//DEPRECATED
func (c *Chunk) GetDelta() (bs []byte) {
	bs = append([]byte{byte(len(c.createdEntities)), byte(len(c.removedEntities))}, c.removedEntities...)
	entChngs := make([][]byte, 0)
	for i,ent := range(c.entities) {
		chngs := ent.GetDelta()
		if len(chngs) > 0 {
			entChngs = append(entChngs, append([]byte{byte(i)}, chngs...))
			ent.ResetAppliedActions()
		}
	}
	bs = append(bs, cmp.Merge(entChngs, c.createdEntities...)...)
	
	for _,crBs := range(c.createdEntities) {
		c.AddEntityFromCreationData(crBs)
	}
	for _,entI := range(c.removedEntities) {
		c.entities[int(entI)] = nil
	}
	c.RemoveNilLocal()
	
	c.changes = make([]int, 0)
	c.createdEntities = make([][]byte, 0)
	c.removedEntities = make([]byte, 0)
	return
}
func (c *Chunk) SetDelta(bs []byte) {
	createdL := int(bs[0])
	rems := int(bs[1])
	lengths := GetSliceOfVal(createdL, ENTITY_CREATION_DATA_LENGTH)
	bs = bs[2:]
	removed := bs[:rems]; bs = bs[rems:]
	createdAndChanges := cmp.Demerge(bs, lengths)
	created := createdAndChanges[:createdL]
	changes := createdAndChanges[createdL:]
	
	for _,chngs := range(changes) {
		idx := int(chngs[0])
		c.entities[idx].SetDelta(chngs[1:])
	}
	
	for _,crBs := range(created) {
		c.AddEntityFromCreationData(crBs)
	}
	for _,entI := range(removed) {
		c.entities[int(entI)] = nil
	}
	c.RemoveNilLocal()
}

func getNewChunkEntityFromBytes(bs []byte, cf *EntityFactory, fcID int) (e *chunkEntity) {
	ent := cf.Get(fcID)
	e = &chunkEntity{ent, [2]byte{}, 0}
	e.FromBytes(bs)
	return e
}
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
