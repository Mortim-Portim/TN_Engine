package TNE

import (
	"fmt"

	ws "github.com/gorilla/websocket"
	"github.com/mortim-portim/GameConn/GC"
	"github.com/mortim-portim/GraphEng/GE"
)

// +-+-+-+-+-+-+-+-+-+-+
// |S|y|n|c|E|n|t|i|t|y|
// +-+-+-+-+-+-+-+-+-+-+
type SyncEntity struct {
	ACIDStart int
	ACIDs     []int
	channel   *GC.Channel

	AllChanged bool

	Entity *Entity
	ef     *EntityFactory

	sm *SmallWorld

	OnNewEntity func(se interface{}, oldE, newE GE.Drawable)
}

func (se *SyncEntity) UpdateSyncVars(m GC.Handler) {
	m.UpdateSyncVarsWithACIDsBuffered(se.ACIDs...)
}
func (se *SyncEntity) HasEntity() bool {
	return se.Entity != nil
}

//Sends the msg as message type msgT to the syncchannel
func (se *SyncEntity) SendToChannel(idx int, msg []byte, force bool) bool {
	return se.channel.SendToPipe(idx, msg, force)
}

/**
Server !ONLY!
sets the entity of se if possible and syncronizes it
**/
func (se *SyncEntity) SetEntity(e *Entity) error {
	if e == nil {
		se.SetNilEntity()
		return fmt.Errorf(ERR_ENTITY_IS_NIL, e)
	}
	if !se.ef.HasEntityID(int(e.FactoryCreationID())) {
		se.SetNilEntity()
		return fmt.Errorf(ERR_NO_FACTORY_FOR_ENTITY, e, e.FactoryCreationID())
	}
	oldE := se.Entity
	se.Entity = e

	se.SendToChannel(SYNCENT_CHAN_ENTITY_CREATION, e.GetCreationData(), true)
	se.UpdateChanFromEnt()
	se.AllChanged = true
	if se.OnNewEntity != nil {
		se.OnNewEntity(se, oldE, se.Entity)
	}
	return nil
}
func (se *SyncEntity) SetNilEntity() {
	if se.OnNewEntity != nil {
		se.OnNewEntity(se, se.Entity, nil)
	}
	se.Entity = nil
}
func (se *SyncEntity) OnChannelChange(sv GC.SyncVar, id int) {
	defer se.channel.ResetJustChanged(SYNCENT_CHAN_ENTITY_CREATION, SYNCENT_CHAN_ACTIONS)
	if se.channel.JustChanged(SYNCENT_CHAN_ENTITY_CREATION) {
		se.CreateEntFromChan()
	} else {
		se.UpdateEntFromChan()
	}
}
func (se *SyncEntity) CreateEntFromChan() error {
	oldE := se.Entity
	ent, err := se.ef.LoadEntityFromCreationData(se.channel.Pipes[SYNCENT_CHAN_ENTITY_CREATION])
	if err != nil {
		return err
	}
	se.Entity = ent
	se.UpdateEntFromChan()
	if se.OnNewEntity != nil {
		se.OnNewEntity(se, oldE, se.Entity)
	}
	return nil
}

func (se *SyncEntity) UpdateChanFromEnt() {
	if se.HasEntity() {
		data := se.Entity.Actions().GetAll()
		if len(data) > 0 {
			se.SendToChannel(SYNCENT_CHAN_ACTIONS, data, true)
		}
	}
}
func (se *SyncEntity) UpdateEntFromChan() {
	if se.HasEntity() {
		if se.channel.JustChanged(SYNCENT_CHAN_ACTIONS) {
			se.Entity.Actions().AppendAndApply(se.channel.Pipes[SYNCENT_CHAN_ACTIONS], se.Entity, se.sm)
		}
	}
}

//Returns a new SyncEntity that will use ef as a creature factory
func GetNewSyncEntity(ACIDStart int, ef *EntityFactory, sm *SmallWorld) (se *SyncEntity) {
	se = &SyncEntity{
		ef:        ef,
		ACIDStart: ACIDStart,
		channel:   GC.GetNewChannel(SYNCENT_CHAN_NUM),
		sm:        sm,
	}
	se.ACIDs = make([]int, SYNCVARS_PER_ENTITY)
	for i := range se.ACIDs {
		se.ACIDs[i] = ACIDStart + i
	}
	return
}
func (se *SyncEntity) GetSyncVars(mp map[int]GC.SyncVar) {
	mp[se.ACIDStart] = se.channel
}

//Registers all syncVars to the server
func (se *SyncEntity) RegisterSyncVars(m *GC.ServerManager, clients ...*ws.Conn) {
	m.RegisterSyncVar(false, se.channel, se.ACIDStart+0, clients...)
}

//Gets all syncVars from the Client
func (se *SyncEntity) GetRegisterdSyncVars(m *GC.ClientManager) {
	se.channel = m.SyncvarsByACID[se.ACIDStart+0].(*GC.Channel)
}
func (se *SyncEntity) RegisterOnChange(m GC.Handler) {
	m.RegisterOnChangeFunc(se.ACIDStart+0, se.OnChannelChange)
}
