package TNE

import (
	"marvin/GraphEng/GE"
	"errors"
	"fmt"
)
const INDEX_FILE_NAME = "#index.txt"

func GetCreatureFactory(path string, frameCounter *int, prepare int) (*CreatureFactory, error) {
	if path[len(path)-1:] != "/" {
		path += "/"
	}
	CF := &CreatureFactory{frameCounter:frameCounter, prepare:prepare}
	CF.rootPath = path
	idx := &GE.List{}
	idx.LoadFromFile(path+INDEX_FILE_NAME)
	CF.crNames = idx.GetSlice()
	CF.mapper = make(map[string]int)
	for i,name := range(CF.crNames) {
		CF.mapper[name] = i
	}
	err := CF.Load()
	CF.prepared = make([][]*Creature, len(CF.creatures))
	return CF, err
}
type CreatureFactory struct {
	rootPath 		string
	crNames 		[]string
	mapper 			map[string]int
	creatures 		[]*Creature
	prepared 		[][]*Creature
	prepare			int
	frameCounter 	*int
}
func (cf *CreatureFactory) Load() error {
	cf.creatures = make([]*Creature, len(cf.crNames))
	for i,name := range(cf.crNames) {
		creature, err := LoadCreature(cf.rootPath+name, cf.frameCounter)
		if err != nil {
			return err
		}
		creature.factoryCreationId = int16(i)
		cf.creatures[i] = creature
	}
	return nil
}
//Should be run on a new goroutine
//Use prepare in order to save time during runtime
//Takes ~30.000ns
func (cf *CreatureFactory) Prepare() {
	for i,cr := range(cf.creatures) {
		if len(cf.prepared[i]) != cf.prepare {
			cf.prepared[i] = make([]*Creature, cf.prepare)
		}
		for idx,_ := range(cf.prepared[i]) {
			if cf.prepared[i][idx] == nil {
				cf.prepared[i][idx] = cr.Copy()
			}
		}
	}
}
//Slower than Get ~1000ns
func (cf *CreatureFactory) GetByName(name string) (*Creature, error) {
	idx, ok := cf.mapper[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("No creature with name: %s", name))
	}
	return cf.Get(idx), nil
}
//Returns a new creature, using a prepared one if possible ~500ns
func (cf *CreatureFactory) Get(idx int) (cr *Creature) {
	crs := cf.prepared[idx]
	for i,lcr := range(crs) {
		if lcr != nil {
			cr = lcr
			cf.prepared[idx][i] = nil
			break
		}
	}
	if cr == nil {
		cr = cf.creatures[idx].Copy()
		if cr == nil {
			panic(fmt.Sprintf("There is no creature with index: ", idx))
		}
	}
	return
}