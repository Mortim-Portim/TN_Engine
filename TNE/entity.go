package TNE

import (
	
)
type EntityUpdater interface {
	Update(e *Entity, world *World)
	Copy() EntityUpdater
}

type Entity struct {
	*Eobj
	
	UpdateCallBack EntityUpdater
}
func (e *Entity) Init() {
	if e.Eobj != nil {
		e.Eobj.RegisterUpdateFunc(e.OnEobjUpdate)
	}
}
func (e *Entity) Copy() (e2 *Entity) {
	e2 = &Entity{Eobj:e.Eobj.Copy(), UpdateCallBack:e.UpdateCallBack}
	e2.Init()
	return 
}


func (e *Entity) OnEobjUpdate(eo *Eobj, w *World) {
	if e.UpdateCallBack != nil {
		e.UpdateCallBack.Update(e, w)
	}
}
func (e *Entity) RegisterUpdateCallback(u EntityUpdater) {
	e.UpdateCallBack = u
}
func LoadEntity(path string, frameCounter *int, c *chan bool) (*Entity, error) {
	eo, err := LoadEobj(path, frameCounter, c)
	if err != nil {return nil,err}
	e := &Entity{Eobj:eo}
	e.Init()
	return e, nil
}