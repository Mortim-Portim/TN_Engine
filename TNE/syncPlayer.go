package TNE

import (
	"github.com/mortim-portim/GameConn/GC"
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
	
	Channel *GC.SyncString
}

func (sp *SyncPlayer) OnChannelChange(sv GC.SyncVar, id int) {
	err, mt, _ := sp.GetFromChannel()
	if err != nil {panic(err)}
	switch mt {
		case SYNCPLAYER_CREATION:
			sp.CreatePlayerFromVars()
	}
}
func (sp *SyncPlayer) CreatePlayerFromVars() error {
	err := sp.se.CreateEntFromVars()
	if err != nil {return err}
	err, mt, data := sp.GetFromChannel()
	if err != nil {return err}
	if mt != SYNCPLAYER_CREATION {
		return fmt.Errorf(ERR_SYNCPLAYER_CREATION, mt, data)
	}
	err, sp.Player = GetPlayerByCreationData(data)
	return err
}
func (sp *SyncPlayer) CreateVarsFromPlayer() error {
	err := sp.se.SetEntity(&sp.Player.Race.Entity)
	if err != nil {return err}
	sp.SendToChannel(SYNCPLAYER_CREATION, sp.Player.GetCreationData())
	return nil
}
func (sp *SyncPlayer) SetPlayer(pl *Player) error {
	if pl == nil {
		return fmt.Errorf(ERR_ENTITY_IS_NIL, pl)
	}
	sp.Player = pl
	sp.CreateVarsFromPlayer()
	return nil
}
func (sp *SyncPlayer) SendToChannel(msgT byte, msg []byte) {
	sp.Channel.SetBs(append([]byte{msgT}, msg...))
}
func (sp *SyncPlayer) GetFromChannel() (error, byte, []byte) {
	data := sp.Channel.GetBs()
	if len(data) < 1 {
		return fmt.Errorf(ERR_CHAN_IS_EMPTY, data), 0, nil
	}
	return nil, data[0], data[1:]
}
//func GetNewSyncEntity(ACIDStart int, ef *EntityFactory) (se *SyncEntity) {
//	se = &SyncEntity{
//		ef:ef,
//		ACIDStart:ACIDStart,
//		X:GC.CreateSyncUInt16(0),
//		Y:GC.CreateSyncUInt16(0),
//		fcID:GC.CreateSyncUInt16(0),
//		Dx:GC.CreateSyncByte(0),
//		Dy:GC.CreateSyncByte(0),
//		extraData:GC.CreateSyncString(""),
//	}
//	return
//}
//func (se *SyncEntity) RegisterSyncVars(m *GC.ServerManager) {
//	m.RegisterSyncVarToAllClients(se.X, 		se.ACIDStart+0)
//	m.RegisterSyncVarToAllClients(se.Y, 		se.ACIDStart+1)
//	m.RegisterSyncVarToAllClients(se.fcID, 		se.ACIDStart+2)
//	m.RegisterSyncVarToAllClients(se.Dx, 		se.ACIDStart+3)
//	m.RegisterSyncVarToAllClients(se.Dy, 		se.ACIDStart+4)
//	m.RegisterSyncVarToAllClients(se.extraData, se.ACIDStart+5)
//}
//func (se *SyncEntity) GetRegisterdSyncVars(m *GC.ClientHandler) {
//	se.X = 			m.SyncvarsByACID[se.ACIDStart+0].(*GC.SyncUInt16)
//	se.Y = 			m.SyncvarsByACID[se.ACIDStart+1].(*GC.SyncUInt16)
//	se.fcID = 		m.SyncvarsByACID[se.ACIDStart+2].(*GC.SyncUInt16)
//	se.Dx = 		m.SyncvarsByACID[se.ACIDStart+3].(*GC.SyncByte)
//	se.Dy = 		m.SyncvarsByACID[se.ACIDStart+4].(*GC.SyncByte)
//	se.extraData =	m.SyncvarsByACID[se.ACIDStart+5].(*GC.SyncString)
//}
//func (se *SyncEntity) RegisterOnChange(m *GC.ClientHandler) {
//	m.RegisterOnChangeFunc(se.ACIDStart+0, se.OnXChange)
//	m.RegisterOnChangeFunc(se.ACIDStart+1, se.OnYChange)
//	m.RegisterOnChangeFunc(se.ACIDStart+2, se.OnfcIDChange)
//	m.RegisterOnChangeFunc(se.ACIDStart+3, se.OnxdChange)
//	m.RegisterOnChangeFunc(se.ACIDStart+4, se.OnydChange)
//	m.RegisterOnChangeFunc(se.ACIDStart+5, se.OnextraDataChange)
//}