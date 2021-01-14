package TNE

import (
	ws "github.com/gorilla/websocket"
	"github.com/mortim-portim/GameConn/GC"
	"github.com/mortim-portim/GraphEng/GE"
	"fmt"
)

const (
	ERR_CHAN_IS_EMPTY = 			"Channel is empty: %v"
	ERR_SYNCPLAYER_CREATION = 		"Cannot create syncplayer with mt: %v, data: %v"
)

const (
	SYNCPLAYER_CREATION = iota
)

// +-+-+-+-+-+-+-+-+-+-+
// |S|y|n|c|P|l|a|y|e|r|
// +-+-+-+-+-+-+-+-+-+-+
type SyncPlayer struct {
	*Player
	se *SyncEntity
	ACIDStart int
	
	Channel *GC.SyncString
	
	OnNewPlayer func(se interface{}, oldE, newE GE.Drawable)
}
func (sp *SyncPlayer) HasPlayer() bool {
	return sp.Player != nil
}
//Sets the player and syncronizes it, i
func (sp *SyncPlayer) SetPlayer(pl *Player) error {
	oldE := sp.Player
	if pl == nil {
		sp.SetNilPlayer()
		return fmt.Errorf(ERR_ENTITY_IS_NIL, pl)
	}
	sp.Player = pl
	sp.CreateVarsFromPlayer()
	if sp.OnNewPlayer != nil {
		sp.OnNewPlayer(sp, oldE, sp.Player)
	}
	return nil
}
func (sp *SyncPlayer) SetNilPlayer() {
	oldE := sp.Player
	sp.Player = nil
	if sp.OnNewPlayer != nil {
		sp.OnNewPlayer(sp, oldE, sp.Player)
	}
}
//Is called when the channel receives
func (sp *SyncPlayer) OnChannelChange(sv GC.SyncVar, id int) {
	err, mt, _ := sp.GetFromChannel()
	if err != nil {
		return
	}
	switch mt {
		case SYNCPLAYER_CREATION:
			sp.CreatePlayerFromVars()
			break;
	}
}
//tries to build the entity and the player from the creation data that should be in the channel
func (sp *SyncPlayer) CreatePlayerFromVars() error {
	//oldE := sp.Player
	err := sp.se.CreateEntFromVars()
	if err != nil {return err}
	err, mt, data := sp.GetFromChannel()
	if err != nil {return err}
	if mt != SYNCPLAYER_CREATION {
		return fmt.Errorf(ERR_SYNCPLAYER_CREATION, mt, data)
	}
	err, sp.Player = GetPlayerByCreationData(data)
	//sp.OnNewPlayer(sp, oldE, sp.Player)
	return err
}
//tries to transfer the entity and send the creation data to the channel
func (sp *SyncPlayer) CreateVarsFromPlayer() error {
	err := sp.se.SetEntity(&sp.Player.Race.Entity)
	if err != nil {return err}
	sp.SendToChannel(SYNCPLAYER_CREATION, sp.Player.GetCreationData())
	return nil
}
func (sp *SyncPlayer) UpdatePlayerFromVars() {
	sp.se.UpdateEntFromVars()
}
func (sp *SyncPlayer) UpdateVarsFromPlayer() {
	sp.se.UpdateVarsFromEnt()
}
//Sends the msg as message type msgT to the syncchannel
func (sp *SyncPlayer) SendToChannel(msgT byte, msg []byte) {
	sp.Channel.SetBs(append([]byte{msgT}, msg...))
}
//Returns the message type, message int the syncchannel
func (sp *SyncPlayer) GetFromChannel() (error, byte, []byte) {
	data := sp.Channel.GetBs()
	if len(data) < 1 {
		return fmt.Errorf(ERR_CHAN_IS_EMPTY, data), 0, nil
	}
	return nil, data[0], data[1:]
}
//Returns an emtpy new SyncPlayer struct, that can build its own player the EntityFactory and Creation data
func GetNewSyncPlayer(ACIDStart int, ef *EntityFactory) (sp *SyncPlayer) {
	sp = &SyncPlayer{
		ACIDStart:ACIDStart,
		se:GetNewSyncEntity(ACIDStart+1, ef),
		Channel:GC.CreateSyncString(""),
	}
	return
}
func (sp *SyncPlayer) GetSyncVars(mp map[int]GC.SyncVar) {
	sp.se.GetSyncVars(mp)
	mp[sp.ACIDStart] = sp.Channel
}
func (sp *SyncPlayer) RegisterSyncVars(m *GC.ServerManager, clients ...*ws.Conn) {
	sp.se.RegisterSyncVars(m, clients...)
	m.RegisterSyncVar(sp.Channel, sp.ACIDStart, clients...)
}
func (sp *SyncPlayer) GetRegisterdSyncVars(m *GC.ClientManager) {
	sp.se.GetRegisterdSyncVars(m)
	sp.Channel = 			m.SyncvarsByACID[sp.ACIDStart].(*GC.SyncString)
}
func (sp *SyncPlayer) RegisterOnChange(m GC.Handler) {
	sp.se.RegisterOnChange(m)
	m.RegisterOnChangeFunc(sp.ACIDStart, sp.OnChannelChange)
}