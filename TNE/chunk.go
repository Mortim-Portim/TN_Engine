package TNE

import (
	"marvin/GraphEng/GE"
	//"github.com/hajimehoshi/ebiten"
	cmp "marvin/GraphEng/Compression"
	"io/ioutil"
	"errors"
)
const CHUNK_SIZE = 16

var ERR_UNKNOWN_ENTITY_ID = errors.New("Unknown Entity ID")
var ERR_ENTITY_NOT_IN_THIS_CHUNK = errors.New("Entity not in this chunk")
var ERR_ENTITY_DOES_NOT_EXIST = errors.New("Entity does not exist")

//tmpPath is a path to a temporary file used for saving the chunk
func GetChunk(x,y int, cf *EntityFactory, tmpPath string) (c *Chunk) {
	c = &Chunk{pos:[2]int16{int16(x),int16(y)}, cf:cf, tmpPath:tmpPath}
	c.tileLT = [2]int16{CHUNK_SIZE*c.pos[0], CHUNK_SIZE*c.pos[1]}
	c.tileRB = [2]int16{c.tileLT[0]+CHUNK_SIZE, c.tileLT[1]+CHUNK_SIZE}
	c.entities = make([][]*chunkEntity, len(cf.EntityNames()))
	c.changes = make([][2]int, 0)
	c.removed = make([][2]int, 0)
	return
}
type Chunk struct {
	pos, tileLT, tileRB [2]int16
	entities [][]*chunkEntity
	changes [][2]int
	removed [][2]int
	cf *EntityFactory
	tmpPath string
	changed bool
	LastUpdateFrame, LastDrawFrame int
}
//Adds all entities of the chunk to drawables
func (c *Chunk) AddToDrawables(dws *GE.Drawables) {
	for _,l := range(c.entities) {
		for _,ent := range(l) {
			dws.Add(ent)
		}
	}
}
//Adds an Entity to the chunk
func (c *Chunk) AddEntity(e EntityI) error {
	id := int(e.FactoryCreationID())
	rx,ry,err := c.RelPosOfEntity(e)
	if err != nil {
		return err
	}else if id < 0 || id >= len(c.entities) {
		return ERR_UNKNOWN_ENTITY_ID
	}else {
		newEnt := getNewChunkEntity(e, rx, ry)
		if len(c.entities[id]) == 0 {
			c.entities[id] = []*chunkEntity{newEnt}
		}else{
			c.entities[id] = append(c.entities[id], newEnt)
		}
		c.changes = append(c.changes, [2]int{id, len(c.entities[id])-1})
		newEnt.SaveChanges()
	}
	return nil
}
//Updates the chunk, returning removed entities
func (c *Chunk) Update(w *World) (removed []*chunkEntity) {
	removed = make([]*chunkEntity, 0)
	for fcID,l := range(c.entities) {
		for idx,entity := range(l) {
			if entity != nil {
				err := entity.Update(c, w)
				if err != nil {
					//Creature is not in this chunk anymore
					removed = append(removed, entity)
					c.entities[fcID][idx] = nil
				}else{
					if entity.Changed() {
						c.changed = true
						c.changes = append(c.changes, [2]int{fcID, idx})
						entity.SaveChanges()
					}
				}
			}
		}
	}
	c.RemoveNil()
	return
}
//Resmoves all entities that are nil
func (c *Chunk) RemoveNil() error {
	for fcID, l := range(c.entities) {
		for i := len(l)-1; i >= 0; i -- {
			if l[i] == nil {
				c.Remove(fcID, i)
				c.removed = append(c.removed, [2]int{fcID, i})
			}
		}
	}
	return nil
}
//Removes a entity with speciefied index
func (c *Chunk) Remove(fcID, i int) {
	c.entities[fcID][i] = c.entities[fcID][len(c.entities[fcID])-1]
	c.entities[fcID] = c.entities[fcID][:len(c.entities[fcID])-1]
}
//returns changes as []byte
func (c *Chunk) GetDelta() (bs []byte) {
	//[1]byte
	bs = []byte{byte(len(c.removed))}
	for _,rem := range(c.removed) {
		//[3]byte
		bs = append(bs, cmp.Int16ToBytes(int16(rem[0]))...)
		bs = append(bs, byte(rem[1]))
	}
	for _,chng := range(c.changes) {
		//[6]byte
		bs = append(bs, cmp.Int16ToBytes(int16(chng[0]))...)
		bs = append(bs, byte(chng[1]))
		bs = append(bs, c.entities[chng[0]][chng[1]].changes...)
	}
	c.changes = make([][2]int, 0)
	c.removed = make([][2]int, 0)
	c.changed = false
	return
}
//sets changes
func (c *Chunk) SetDelta(bs []byte) (removed []*chunkEntity) {
	removed = make([]*chunkEntity, 0)
	rems := int(bs[0]); bs = bs[1:]
	for i := 0; i < rems; i ++ {
		fcID := int(cmp.BytesToInt16(bs[0:1]))
		idx :=	int(bs[2])
		removed = append(removed, c.entities[fcID][idx])
		c.Remove(fcID, idx)
		bs = bs[3:]
	}
	for i := 0; i < len(bs)/6; i++ {
		fcID := int(cmp.BytesToInt16(bs[0:1]))
		idx :=	int(bs[2])
		if idx >= len(c.entities[fcID]) {
			c.entities[fcID] = append(c.entities[fcID], getNewChunkEntityFromBytes(bs[3:5], c.cf, fcID))
		}else{
			c.entities[fcID][idx].FromBytes(bs[3:5])
		}
		bs = bs[6:]
	}
	return
}
//Writes the chunk to the disk
func (c *Chunk) ToDisk() error {
	bs := make([]byte, 0)
	for fcID,l := range(c.entities) {
		bs = append(bs, cmp.Int16ToBytes(int16(fcID))...)
		for idx,entity := range(l) {
			bs = append(bs, append(entity.ToBytes(), byte(idx))...)
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
			c.entities[fcID][idx].EntityI = c.cf.Get(fcID)
			c.entities[fcID][idx].FromBytes(data[:3])
			data = data[4:]
		}
	}
	return nil
}
//Returns the relative position of a entity in a chunk
func (c *Chunk) RelPosOfEntity(e EntityI) (byte, byte, error) {
	eX, eY := e.IntPos()
	relX, relY := eX-int64(c.tileLT[0]), eY-int64(c.tileLT[1])
	if relX < 0 || relY < 0 || relX >= CHUNK_SIZE || relY >= CHUNK_SIZE {
		return 0,0, ERR_ENTITY_NOT_IN_THIS_CHUNK
	}
	return byte(relX), byte(relY), nil
}
//converts 2d coords in a chunk to a index
func ChunkCoord2DtoIdx(x, y int) byte {
	if x >= CHUNK_SIZE || y >= CHUNK_SIZE {
		panic("NEVER call ChunkCoord2DtoIdx with coords >= 16")
	}
	return byte(x+CHUNK_SIZE*y)
}
//converts a index in a chunk to 2d coords
func IdxtoChunkCoord2D(idx byte) (x,y int) {
	csm1 := byte(CHUNK_SIZE -1)
	x = int(idx%CHUNK_SIZE)
	y = int((idx-(idx%csm1))/csm1)
	return
}
func getNewChunkEntityFromBytes(bs []byte, cf *EntityFactory, fcID int) (e *chunkEntity) {
	ent := cf.Get(fcID)
	e = &chunkEntity{ent, [2]byte{}, 0, nil}
	e.FromBytes(bs)
	return e
}
func getNewChunkEntity(e EntityI, rx, ry byte) *chunkEntity {
	return &chunkEntity{e, [2]byte{rx,ry}, ChunkCoord2DtoIdx(int(rx), int(ry)), nil}
}
type chunkEntity struct {
	EntityI
	chunkPos [2]byte
	chunkPosIdx byte
	//[3]byte
	changes []byte
}
func (ce *chunkEntity) SaveChanges() {
	ce.changes = ce.ToBytes()
}
func (ce *chunkEntity) FromBytes(bs []byte) {
	ce.EntityI.SetData(bs[0:1])
	ce.chunkPosIdx = bs[2]
	x,y := IdxtoChunkCoord2D(bs[2])
	ce.chunkPos = [2]byte{byte(x),byte(y)}
	ce.EntityI.SetTopLeftTo(float64(x), float64(y))
}
func (ce *chunkEntity) ToBytes() (bs []byte) {
	defer func(){ce.EntityI = nil}()
	bs = make([]byte, 3)
	copy(bs[0:1], ce.EntityI.GetData())
	bs[2] = ce.chunkPosIdx
	return
}
func (ce *chunkEntity) Update(c *Chunk, w *World) error {
	ce.EntityI.Update(w)
	rx,ry,err := c.RelPosOfEntity(ce.EntityI)
	if err != nil {return err}
	ce.chunkPos[0] = rx
	ce.chunkPos[1] = ry
	ce.chunkPosIdx = ChunkCoord2DtoIdx(int(rx), int(ry))
	return nil
}