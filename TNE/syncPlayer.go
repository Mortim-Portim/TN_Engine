package TNE

import (
	ws "github.com/gorilla/websocket"
	"github.com/mortim-portim/GameConn/GC"
	"github.com/mortim-portim/GraphEng/GE"
)

const (
	ERR_CHAN_IS_EMPTY       = "Channel is empty: %v"
	ERR_SYNCPLAYER_CREATION = "Cannot create syncplayer with mt: %v, data: %v"
)

// +-+-+-+-+-+-+-+-+-+-+
// |S|y|n|c|P|l|a|y|e|r|
// +-+-+-+-+-+-+-+-+-+-+
type SyncPlayer struct {
	*Player
	Se        *SyncEntity
	ACIDStart int
	ACIDs     []int

	OnNewPlayer func(se interface{}, oldE, newE GE.Drawable)
}

func (sp *SyncPlayer) UpdateSyncVars(m GC.Handler) {
	m.UpdateSyncVarsWithACIDsBuffered(sp.ACIDs...)
}
func (sp *SyncPlayer) HasPlayer() bool {
	return sp.Player != nil
}

//Sets the player and syncronizes it, i
func (sp *SyncPlayer) SetPlayer(pl *Player) {
	oldE := sp.Player
	sp.Player = pl
	if pl == nil {
		sp.SendToChannel(SYNCENT_CHAN_PLAYER_CREATION, []byte{0}, true)
	} else {
		err := sp.CreateVarsFromPlayer()
		if err != nil {
			panic(err)
		}
	}
	if sp.OnNewPlayer != nil {
		sp.OnNewPlayer(sp, oldE, sp.Player)
	}
}
func (sp *SyncPlayer) SetNilPlayer() {
	if sp.OnNewPlayer != nil {
		sp.OnNewPlayer(sp, sp.Player, nil)
	}
	sp.Player = nil
}

//Is called when the channel receives
func (sp *SyncPlayer) OnChannelChange(sv GC.SyncVar, id int) {
	defer sp.Se.channel.ResetJustChanged(SYNCENT_CHAN_PLAYER_CREATION, SYNCENT_CHAN_PLAYER_CREATION)
	sp.Se.OnChannelChange(sv, id)
	if sp.Se.channel.JustChanged(SYNCENT_CHAN_PLAYER_CREATION) {
		data := sp.Se.channel.Pipes[SYNCENT_CHAN_PLAYER_CREATION]
		if data[0] == 0 {
			sp.SetNilPlayer()
		} else if data[0] == 1 {
			err := sp.CreatePlayerFromVars(data[1:])
			if err != nil {
				panic(err)
			}
		}
	}
}
func (sp *SyncPlayer) CreatePlayerFromVars(data []byte) error {
	oldE := sp.Player
	err := sp.Se.CreateEntFromChan()
	if err != nil {
		return err
	}
	pl, err := GetPlayerByCreationData(data)
	if err != nil {
		return err
	}
	sp.Player = pl
	sp.Player.Entity = sp.Se.Entity
	if sp.OnNewPlayer != nil {
		sp.OnNewPlayer(sp, oldE, sp.Player)
	}
	return err
}
func (sp *SyncPlayer) CreateVarsFromPlayer() error {
	err := sp.Se.SetEntity(sp.Player.Entity)
	if err != nil {
		return err
	}
	data := sp.Player.GetCreationData()
	sp.SendToChannel(SYNCENT_CHAN_PLAYER_CREATION, append([]byte{1}, data...), true)
	return nil
}
func (sp *SyncPlayer) UpdatePlayerFromChan() {
	sp.Se.UpdateEntFromChan()
}
func (sp *SyncPlayer) UpdateChanFromPlayer() {
	sp.Se.UpdateChanFromEnt()
}

//Sends the msg as message type msgT to the syncchannel
func (sp *SyncPlayer) SendToChannel(idx int, msg []byte, force bool) bool {
	return sp.Se.SendToChannel(idx, msg, force)
}

//Returns an emtpy new SyncPlayer struct, that can build its own player the EntityFactory and Creation data
func GetNewSyncPlayer(ACIDStart int, ef *EntityFactory) (sp *SyncPlayer) {
	sp = &SyncPlayer{
		ACIDStart: ACIDStart,
		Se:        GetNewSyncEntity(ACIDStart, ef),
	}
	sp.ACIDs = make([]int, SYNCVARS_PER_PLAYER)
	for i := range sp.ACIDs {
		sp.ACIDs[i] = ACIDStart + i
	}
	return
}
func (sp *SyncPlayer) GetSyncVars(mp map[int]GC.SyncVar) {
	sp.Se.GetSyncVars(mp)
}
func (sp *SyncPlayer) RegisterSyncVars(m *GC.ServerManager, clients ...*ws.Conn) {
	sp.Se.RegisterSyncVars(m, clients...)
}
func (sp *SyncPlayer) GetRegisterdSyncVars(m *GC.ClientManager) {
	sp.Se.GetRegisterdSyncVars(m)
}
func (sp *SyncPlayer) RegisterOnChange(m GC.Handler) {
	//sp.Se.RegisterOnChange(m)
	m.RegisterOnChangeFunc(sp.ACIDStart, sp.OnChannelChange)
}
