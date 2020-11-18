package TNE

import (
	"errors"
	"fmt"

	"github.com/mortim-portim/GraphEng/GE"
)

const INDEX_FILE_NAME = "#index.txt"

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
func (cf *EntityFactory) SetUpdateFunction(fncs map[string]func(e *Entity, world *World)) {
	for name, idx := range cf.mapper {
		fnc, ok := fncs[name]
		if ok {
			cf.entities[idx].RegiserUpdateFunc(fnc)
		}
	}
	return
}

//Loads the entities
func (cf *EntityFactory) Load() error {
	cf.entities = make([]*Entity, len(cf.crNames))
	for i, name := range cf.crNames {
		ent, err := LoadEntity(cf.rootPath+name, cf.frameCounter)
		if err != nil {
			return err
		}
		ent.factoryCreationId = int16(i)
		cf.entities[i] = ent
	}
	return nil
}

//Should be run on a new goroutine
//Use prepare in order to save time during runtime
//Takes ~30.000ns
func (cf *EntityFactory) Prepare() {
	//done := make(chan bool)
	for i, cr := range cf.entities {
		func() {
			if len(cf.prepared[i]) != cf.prepare {
				cf.prepared[i] = make([]*Entity, cf.prepare)
			}
			for idx, _ := range cf.prepared[i] {
				if cf.prepared[i][idx] == nil {
					cf.prepared[i][idx] = cr.Copy()
				}
			}
			//done <- true
		}()
	}
	/**
	for range(cf.entities) {
		<- done
	}
	**/
}

//Slower than Get ~1000ns
func (cf *EntityFactory) GetByName(name string) (*Entity, error) {
	idx, ok := cf.mapper[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("No entity with name: %s", name))
	}
	return cf.Get(idx), nil
}

//Returns a new entity, using a prepared one if possible ~500ns
func (cf *EntityFactory) Get(idx int) (cr *Entity) {
	crs := cf.prepared[idx]
	for i, lcr := range crs {
		if lcr != nil {
			cr = lcr
			cf.prepared[idx][i] = nil
			break
		}
	}
	if cr == nil {
		cr = cf.entities[idx].Copy()
		if cr == nil {
			panic(fmt.Sprintf("There is no entity with index: ", idx))
		}
	}
	return
}

//Returns a slice containing the names of all entities
func (cf *EntityFactory) EntityNames() []string {
	return cf.crNames
}
