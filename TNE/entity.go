package TNE

import (
	
)
type EntityUpdater interface {
	Update(e *Entity, world *World)
}

type Entity struct {
	*Eobj
	
	UpdateCallBack EntityUpdater
}
func (e *Entity) Copy() (*Entity) {
	return &Entity{Eobj:e.Eobj.Copy()}
}


func (e *Entity) OnEobjUpdate(eo *Eobj, w *World) {
	if e.UpdateCallBack != nil {
		e.UpdateCallBack.Update(e, w)
	}
}
func (e *Entity) RegisterUpdateCallback(u EntityUpdater) {
	e.UpdateCallBack = u
}
func LoadEntity(path string, frameCounter *int) (*Entity, error) {
	eo, err := LoadEobj(path, frameCounter)
	if err != nil {return nil,err}
	e := &Entity{Eobj:eo}
	e.Eobj.RegisterUpdateFunc(e.OnEobjUpdate)
	return e, nil
}