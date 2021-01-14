package TNE

import (
	//cmp "github.com/mortim-portim/GraphEng/Compression"
	ws "github.com/gorilla/websocket"
	"github.com/mortim-portim/GameConn/GC"
	"github.com/mortim-portim/GraphEng/GE"
	"fmt"
)

const (
	//amount of syncVars needed by one Entity
	SYNCVARS_PER_ENTITY = 6
	//count of SyncEntities to be prepared
	SYNCENTITIES_PREP = 100
	//amount of syncVars needed by one Player
	SYNCVARS_PER_PLAYER = 7
	//count of SyncPlayer to be prepared besides the own player
	SYNCPLAYER_PREP =	 10
)
const (
	OTHERPLAYERS = SYNCPLAYER_PREP
	//SyncVars that are registered by the world
	WorldStructChan_ACID = iota
	WorldFrameChan_ACID
	WorldLightLevelChan_ACID
	WORLD_SYNCVARS
)
//Returns the startACID for the own player
func GetSVACID_Start_OwnPlayer() int {
	return WORLD_SYNCVARS
}
//Returns the startACID for the other player
func GetSVACID_Start_OtherPlayer(idx int) int {
	return GetSVACID_Start_OwnPlayer()+1+idx*SYNCVARS_PER_PLAYER
}
//Returns the startACID for the entity
func GetSVACID_Start_Entities(idx int) int {
	return GetSVACID_Start_OtherPlayer(OTHERPLAYERS)+idx*SYNCVARS_PER_ENTITY
}
//Initializes the SyncClient
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
	
	OnNewEntity func(se interface{}, oldE, newE GE.Drawable)
}
func (se *SyncEntity) HasEntity() bool {
	return se.Entity != nil
}
/**
Server !ONLY!
sets the entity of se if possible and syncronizes it
**/
func (se *SyncEntity) SetEntity(e *Entity) error {
	oldE := se.Entity
	if e == nil {
		se.SetNilEntity()
		return fmt.Errorf(ERR_ENTITY_IS_NIL, e)
	}
	if !se.ef.HasEntityID(int(e.FactoryCreationID())) {
		se.SetNilEntity()
		return fmt.Errorf(ERR_NO_FACTORY_FOR_ENTITY, e, e.FactoryCreationID())
	}
	se.Entity = e
	se.UpdateVarsFromEnt()
	se.fcID.SetInt(uint16(e.FactoryCreationID()))
	se.OnNewEntity(se, oldE, se.Entity)
	return nil
}
func (se *SyncEntity) SetNilEntity() {
	oldE := se.Entity
	se.Entity = nil
	se.OnNewEntity(se, oldE, nil)
}
//Called when the x-position syncVar changes
func (se *SyncEntity) OnXChange(sv GC.SyncVar, id int) {
	se.UpdateEntFromVars()
}
//Called when the y-position syncVar changes
func (se *SyncEntity) OnYChange(sv GC.SyncVar, id int) {
	se.UpdateEntFromVars()
}
//Called when the x-difference syncVar changes
func (se *SyncEntity) OnxdChange(sv GC.SyncVar, id int) {
	se.UpdateEntFromVars()
}
//Called when the y-difference syncVar changes
func (se *SyncEntity) OnydChange(sv GC.SyncVar, id int) {
	se.UpdateEntFromVars()
}
//Called when the extra data syncVar changes
func (se *SyncEntity) OnextraDataChange(sv GC.SyncVar, id int) {
	
}
//Called when the fcID syncVar changes
func (se *SyncEntity) OnfcIDChange(sv GC.SyncVar, id int) {
	se.CreateEntFromVars()
}
//Tries to create the Entity from the SyncVars
func (se *SyncEntity) CreateEntFromVars() error {
	oldE := se.Entity
	ent, err := se.ef.Get(int(se.fcID.GetInt()))
	if err != nil {return err}
	se.Entity = ent
	se.UpdateEntFromVars()
	se.OnNewEntity(se, oldE, se.Entity)
	return nil
}
//Updates the SyncVars from the Entity if possible
func (se *SyncEntity) UpdateVarsFromEnt() {
	if se.HasEntity() {
		x,y := se.Entity.IntPos()
		se.X.SetInt(uint16(x));se.X.SetInt(uint16(y))
		dx,dy := se.Entity.GetPosDelta()
		se.Dx.SetByte(dx);se.Dy.SetByte(dy)
	}
}
//Updates the Entity from the SyncVars if possible
func (se *SyncEntity) UpdateEntFromVars() {
	if se.HasEntity() {
		se.Entity.SetPosDelta(int(se.X.GetInt()), int(se.Y.GetInt()), se.Dx.GetByte(), se.Dy.GetByte())
	}
}
//Returns a new SyncEntity that will use ef as a creature factory
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
//Registers all syncVars to the server
func (se *SyncEntity) RegisterSyncVars(m *GC.ServerManager, clients ...*ws.Conn) {
	m.RegisterSyncVar(se.X, 		se.ACIDStart+0, clients...)
	m.RegisterSyncVar(se.Y, 		se.ACIDStart+1, clients...)
	m.RegisterSyncVar(se.fcID, 		se.ACIDStart+2, clients...)
	m.RegisterSyncVar(se.Dx, 		se.ACIDStart+3, clients...)
	m.RegisterSyncVar(se.Dy, 		se.ACIDStart+4, clients...)
	m.RegisterSyncVar(se.extraData, se.ACIDStart+5, clients...)
}
//Gets all syncVars from the Client
func (se *SyncEntity) GetRegisterdSyncVars(m *GC.ClientManager) {
	se.X = 			m.SyncvarsByACID[se.ACIDStart+0].(*GC.SyncUInt16)
	se.Y = 			m.SyncvarsByACID[se.ACIDStart+1].(*GC.SyncUInt16)
	se.fcID = 		m.SyncvarsByACID[se.ACIDStart+2].(*GC.SyncUInt16)
	se.Dx = 		m.SyncvarsByACID[se.ACIDStart+3].(*GC.SyncByte)
	se.Dy = 		m.SyncvarsByACID[se.ACIDStart+4].(*GC.SyncByte)
	se.extraData =	m.SyncvarsByACID[se.ACIDStart+5].(*GC.SyncString)
}
func (se *SyncEntity) RegisterOnChange(m GC.Handler) {
	m.RegisterOnChangeFunc(se.ACIDStart+0, se.OnXChange)
	m.RegisterOnChangeFunc(se.ACIDStart+1, se.OnYChange)
	m.RegisterOnChangeFunc(se.ACIDStart+2, se.OnfcIDChange)
	m.RegisterOnChangeFunc(se.ACIDStart+3, se.OnxdChange)
	m.RegisterOnChangeFunc(se.ACIDStart+4, se.OnydChange)
	m.RegisterOnChangeFunc(se.ACIDStart+5, se.OnextraDataChange)
}