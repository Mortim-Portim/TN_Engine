package TNE

import (
	//cmp "github.com/mortim-portim/GraphEng/Compression"
	"github.com/mortim-portim/GameConn/GC"
	"fmt"
)
/**
TODO
Implement SyncPlayer
**/
const(
	SYNCENTITY_ENT_CHANGE = iota
)
const (
	SYNCVARS_PER_ENTITY = 6
	SYNCENTITIES_PREP = 100
	SYNCVARS_PER_PLAYER = 6
	SYNCPLAYER_PREP =	 10
	
	ADDITIONAL_SYNCVARS = 0
)
func GetACIDStart(ents, plrs int) int {return ents*SYNCVARS_PER_ENTITY+plrs*SYNCVARS_PER_PLAYER}
func InitialSyncClient() {
	GC.InitSyncVarStandardTypes()
}
const (
	ERR_ENTITY_IS_NIL = 			"Entity: %v is nil"
	ERR_NO_FACTORY_FOR_ENTITY = 	"No factory for Entity: %v, with fcID: %v"
)

// +-+-+-+-+-+-+-+-+-+-+
// |S|y|n|c|E|n|t|i|t|y|
// +-+-+-+-+-+-+-+-+-+-+
type SyncEntity struct {
	ACIDStart int
	X,Y, fcID *GC.SyncUInt16
	Dx,Dy *GC.SyncByte
	extraData *GC.SyncString
	
	Entity *Entity
	ef *EntityFactory
}
func (se *SyncEntity) HasEntity() bool {
	return se.Entity != nil
}

func (se *SyncEntity) Update(w *World) {
	if se.HasEntity() {
		se.Entity.UpdateAll(w)
		se.UpdateVarsFromEnt()
	}
}
func (se *SyncEntity) SetEntity(e *Entity) error {
	if e == nil {
		return fmt.Errorf(ERR_ENTITY_IS_NIL, e)
	}
	if !se.ef.HasEntityID(int(e.FactoryCreationID())) {
		return fmt.Errorf(ERR_NO_FACTORY_FOR_ENTITY, e, e.FactoryCreationID())
	}
	se.Entity = e
	se.UpdateVarsFromEnt()
	se.fcID.SetInt(uint16(e.FactoryCreationID()))
	return nil
}

func (se *SyncEntity) OnXChange(sv GC.SyncVar, id int) {
	se.UpdateEntFromVars()
}
func (se *SyncEntity) OnYChange(sv GC.SyncVar, id int) {
	se.UpdateEntFromVars()
}
func (se *SyncEntity) OnxdChange(sv GC.SyncVar, id int) {
	se.UpdateEntFromVars()
}
func (se *SyncEntity) OnydChange(sv GC.SyncVar, id int) {
	se.UpdateEntFromVars()
}
func (se *SyncEntity) OnextraDataChange(sv GC.SyncVar, id int) {
	
}
func (se *SyncEntity) OnfcIDChange(sv GC.SyncVar, id int) {
	se.CreateEntFromVars()
}
func (se *SyncEntity) CreateEntFromVars() error {
	ent, err := se.ef.Get(int(se.fcID.GetInt()))
	if err != nil {return err}
	se.Entity = ent
	se.UpdateEntFromVars()
	return nil
}
func (se *SyncEntity) UpdateVarsFromEnt() {
	x,y := se.Entity.IntPos()
	se.X.SetInt(uint16(x));se.X.SetInt(uint16(y))
	dx,dy := se.Entity.GetPosDelta()
	se.Dx.SetByte(dx);se.Dy.SetByte(dy)
}
func (se *SyncEntity) UpdateEntFromVars() {
	se.Entity.SetPosDelta(int(se.X.GetInt()), int(se.Y.GetInt()), se.Dx.GetByte(), se.Dy.GetByte())
}
func GetNewSyncEntity(ACIDStart int, ef *EntityFactory) (se *SyncEntity) {
	se = &SyncEntity{
		ef:ef,
		ACIDStart:ACIDStart,
		X:GC.CreateSyncUInt16(0),
		Y:GC.CreateSyncUInt16(0),
		fcID:GC.CreateSyncUInt16(0),
		Dx:GC.CreateSyncByte(0),
		Dy:GC.CreateSyncByte(0),
		extraData:GC.CreateSyncString(""),
	}
	return
}
func (se *SyncEntity) RegisterSyncVars(m *GC.ServerManager) {
	m.RegisterSyncVarToAllClients(se.X, 		se.ACIDStart+0)
	m.RegisterSyncVarToAllClients(se.Y, 		se.ACIDStart+1)
	m.RegisterSyncVarToAllClients(se.fcID, 		se.ACIDStart+2)
	m.RegisterSyncVarToAllClients(se.Dx, 		se.ACIDStart+3)
	m.RegisterSyncVarToAllClients(se.Dy, 		se.ACIDStart+4)
	m.RegisterSyncVarToAllClients(se.extraData, se.ACIDStart+5)
}
func (se *SyncEntity) GetRegisterdSyncVars(m *GC.ClientHandler) {
	se.X = 			m.SyncvarsByACID[se.ACIDStart+0].(*GC.SyncUInt16)
	se.Y = 			m.SyncvarsByACID[se.ACIDStart+1].(*GC.SyncUInt16)
	se.fcID = 		m.SyncvarsByACID[se.ACIDStart+2].(*GC.SyncUInt16)
	se.Dx = 		m.SyncvarsByACID[se.ACIDStart+3].(*GC.SyncByte)
	se.Dy = 		m.SyncvarsByACID[se.ACIDStart+4].(*GC.SyncByte)
	se.extraData =	m.SyncvarsByACID[se.ACIDStart+5].(*GC.SyncString)
}
func (se *SyncEntity) RegisterOnChange(m *GC.ClientHandler) {
	m.RegisterOnChangeFunc(se.ACIDStart+0, se.OnXChange)
	m.RegisterOnChangeFunc(se.ACIDStart+1, se.OnYChange)
	m.RegisterOnChangeFunc(se.ACIDStart+2, se.OnfcIDChange)
	m.RegisterOnChangeFunc(se.ACIDStart+3, se.OnxdChange)
	m.RegisterOnChangeFunc(se.ACIDStart+4, se.OnydChange)
	m.RegisterOnChangeFunc(se.ACIDStart+5, se.OnextraDataChange)
}