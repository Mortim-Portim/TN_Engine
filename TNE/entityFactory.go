package TNE

import (
	"fmt"
	"github.com/mortim-portim/GraphEng/GE"
)

const INDEX_FILE_NAME = "#index.txt"
const(
	ERR_NO_FACTORY_FOR_ENTITY_BY_NAME = 	"No factory for Entity: %v, with fcID: %v and Name: %s"
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
		}else{
			idx, _ := cf.mapper[name]
			cf.entities[idx].RegiserUpdateFunc(fnc)
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
		cf.entities[i].RegiserUpdateFunc(fnc)
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
		for idx, _ := range cf.prepared[i] {
			if cf.prepared[i][idx] == nil {
				cf.prepared[i][idx] = cr.Copy()
			}
		}
	}
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
