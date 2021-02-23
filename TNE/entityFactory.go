package TNE

import (
	"fmt"

	"github.com/mortim-portim/GraphEng/GE"
	cmp "github.com/mortim-portim/GraphEng/compression"
)

const INDEX_FILE_NAME = "#index.txt"
const (
	ERR_NO_FACTORY_FOR_ENTITY_BY_NAME = "No factory for Entity: %v, with fcID: %v and Name: %s"
)

//returns a entity factory that loads all entities from a specific path, prepare specifies the number of entities to be prepared
func GetEntityFactory(path string, frameCounter *int, prepare int) (*EntityFactory, error) {
	if path[len(path)-1:] != "/" {
		path += "/"
	}
	CF := &EntityFactory{frameCounter: frameCounter, prepare: prepare}
	CF.rootPath = path
	idx := &GE.List{}
	idx.LoadFromFile(path + INDEX_FILE_NAME)
	CF.crNames = idx.GetSlice()
	CF.mapper = make(map[string]int)
	for i, name := range CF.crNames {
		CF.mapper[name] = i
	}
	err := CF.Load()
	CF.prepared = make([][]*Entity, len(CF.entities))
	return CF, err
}

type EntityFactory struct {
	rootPath     string
	crNames      []string
	mapper       map[string]int
	entities     []*Entity
	prepared     [][]*Entity
	prepare      int
	frameCounter *int
}

func (cf *EntityFactory) LoadEntityFromCreationData(data []byte) (*Entity, error) {
	fcID := int(cmp.BytesToInt16(data[0:2]))
	e, err := cf.Get(fcID)
	if err != nil {
		return nil, err
	}
	e.PosFromBytes(data[2:8])
	e.orientation.FromByte(data[8])
	e.neworientation.FromByte(data[9])
	e.movingFrames = int(cmp.BytesToInt16(data[10:12]))
	e.movedFrames = int(cmp.BytesToInt16(data[12:14]))
	e.movingStepSize = cmp.BytesToFloat64(data[14:22])
	e.setAnim(uint8(data[22]))
	vals := cmp.BytesToFloat32s(data[23:47])
	e.SetMaxHealth(vals[0])
	e.SetMaxStamina(vals[1])
	e.SetMaxMana(vals[2])
	e.SetHealth(vals[3])
	e.SetStamina(vals[4])
	e.SetMana(vals[5])
	shows := cmp.BytesToBools(data[47:48])
	e.ShowHealth(shows[0])
	e.ShowStamina(shows[1])
	e.ShowMana(shows[2])
	e.ID = cmp.BytesToInt16(data[48:50])
	e.Char, err = LoadChar(data[50 : 50+CHARACTER_BYTES_LENGTH])
	return e, err
}
func (cf *EntityFactory) Print() (out string) {
	out = fmt.Sprintf("Path: %v, crNames: %v, entities: %v, prepare: %v, frame: %v",
		cf.rootPath, cf.crNames, len(cf.entities), cf.prepare, *cf.frameCounter)
	return
}
func (cf *EntityFactory) SetUpdateFunctionMap(fncs map[string]EntityUpdater) error {
	err := error(nil)
	for name, fnc := range fncs {
		if !cf.HasEntityName(name) {
			err = fmt.Errorf(ERR_NO_FACTORY_FOR_ENTITY_BY_NAME, nil, -1, name)
		} else {
			idx, _ := cf.mapper[name]
			cf.entities[idx].RegisterUpdateCallback(fnc)
		}
	}
	return err
}
func (cf *EntityFactory) SetUpdateFunctionList(fncs []EntityUpdater) error {
	err := error(nil)
	if len(fncs) > len(cf.entities) {
		err = fmt.Errorf(ERR_NO_FACTORY_FOR_ENTITY, nil, len(cf.entities))
		fncs = fncs[:len(cf.entities)]
	}

	for i, fnc := range fncs {
		cf.entities[i].RegisterUpdateCallback(fnc)
	}
	return err
}

//Loads the entities
func (cf *EntityFactory) Load() error {
	cf.entities = make([]*Entity, len(cf.crNames))
	for i, name := range cf.crNames {
		ent, err := LoadEntity(cf.rootPath+name, cf.frameCounter)
		if err != nil {
			cf.entities = cf.entities[:i]
			return err
		}
		ent.factoryCreationId = int16(i)
		cf.entities[i] = ent
	}
	return nil
}

//Should be run on a new goroutine
//Use prepare in order to save time during runtime
//Takes some time
func (cf *EntityFactory) Prepare() {
	for i, cr := range cf.entities {
		if len(cf.prepared[i]) != cf.prepare {
			cf.prepared[i] = make([]*Entity, cf.prepare)
		}
		for idx := range cf.prepared[i] {
			if cf.prepared[i][idx] == nil {
				cf.prepared[i][idx] = cr.Copy()
			}
		}
	}
}
func (cf *EntityFactory) GetFromCharacter(char *Character) (*Entity, error) {
	ent, err := cf.GetByName(char.Race.Name)
	if err != nil {
		return ent, err
	}
	ent.Char = char
	return ent, nil
}

//Slower than Get ~1000ns
func (cf *EntityFactory) GetByName(name string) (*Entity, error) {
	if !cf.HasEntityName(name) {
		return nil, fmt.Errorf(ERR_NO_FACTORY_FOR_ENTITY_BY_NAME, nil, -1, name)
	}
	idx, _ := cf.mapper[name]
	return cf.Get(idx)
}

//Returns a new entity, using a prepared one if possible ~500ns
func (cf *EntityFactory) Get(idx int) (cr *Entity, err error) {
	if !cf.HasEntityID(idx) {
		return nil, fmt.Errorf(ERR_NO_FACTORY_FOR_ENTITY, nil, idx)
	}
	for i, lcr := range cf.prepared[idx] {
		if lcr != nil {
			cr = lcr
			cf.prepared[idx][i] = nil
			break
		}
	}
	if cr == nil {
		cr = cf.entities[idx].Copy()
	}
	return
}
func (cf *EntityFactory) HasEntityID(fcID int) bool {
	if fcID >= 0 && fcID < len(cf.entities) {
		return true
	}
	return false
}
func (cf *EntityFactory) HasEntityName(name string) bool {
	_, ok := cf.mapper[name]
	return ok
}

//Returns a slice containing the names of all entities
func (cf *EntityFactory) EntityNames() []string {
	return cf.crNames
}
