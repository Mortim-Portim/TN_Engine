package TNE

import (
	"fmt"
	"time"

	ws "github.com/gorilla/websocket"
	"github.com/hajimehoshi/ebiten"
	"github.com/mortim-portim/GameConn/GC"
	"github.com/mortim-portim/GraphEng/GE"
	cmp "github.com/mortim-portim/GraphEng/compression"
)

const (
	SMALLWORLD_WORLDSTRUCTURE_MSG = iota
	SMALLWORLD_TIMEPERFRAME_MSG
	SMALLWORLD_FRAMEANDTIME_MSG
	SMALLWORLD_SET_ACTIVEPLAYER_ID

	SMALLWORLD_CHAN_TO_CLIENT_PIPES
)

const (
	SMALLWORLD_PLACEHOLDER_TOSERV = iota

	SMALLWORLD_CHAN_TO_SERVER_PIPES
)

func GetSmallWorld(X, Y, W, H float64, tile_path, struct_path, entity_path string) (sm *SmallWorld, err error) {
	fc := 0
	ef, err := GetEntityFactory(entity_path, &fc, 3)
	sm = &SmallWorld{Ents: make([]*SyncEntity, SYNCENTITIES_PREP),
		Plys:         make([]*SyncPlayer, SYNCPLAYER_PREP),
		ChanToClient: GC.GetNewChannel(SMALLWORLD_CHAN_TO_CLIENT_PIPES),
		ChanToServer: GC.GetNewChannel(SMALLWORLD_CHAN_TO_SERVER_PIPES),
		Ef:           ef, X: X, Y: Y, W: W, H: H, tile_path: tile_path, struct_path: struct_path,
		FrameCounter: &fc,
		ActivePlayer: GetNewSyncPlayer(GetSVACID_Start_OwnPlayer(), ef, sm),
	}
	for i := range sm.Ents {
		sm.Ents[i] = GetNewSyncEntity(GetSVACID_Start_Entities(i), ef, sm)
	}
	for i := range sm.Plys {
		sm.Plys[i] = GetNewSyncPlayer(GetSVACID_Start_OtherPlayer(i), ef, sm)
	}
	return
}
func (sm *SmallWorld) New() (sm2 *SmallWorld) {
	sm2 = &SmallWorld{Ents: make([]*SyncEntity, SYNCENTITIES_PREP),
		Plys:         make([]*SyncPlayer, SYNCPLAYER_PREP),
		ChanToClient: GC.GetNewChannel(SMALLWORLD_CHAN_TO_CLIENT_PIPES),
		ChanToServer: GC.GetNewChannel(SMALLWORLD_CHAN_TO_SERVER_PIPES),
		Ef:           sm.Ef, X: sm.X, Y: sm.Y, W: sm.W, H: sm.H, tile_path: sm.tile_path, struct_path: sm.struct_path,
		FrameCounter:        sm.FrameCounter,
		ActivePlayer:        GetNewSyncPlayer(GetSVACID_Start_OwnPlayer(), sm.Ef, sm2),
		Struct:              sm.Struct,
		FrameChanSendPeriod: sm.FrameChanSendPeriod,
		TimePerFrame:        sm.TimePerFrame,
		World:               sm.World,
	}
	for i := range sm2.Ents {
		sm2.Ents[i] = GetNewSyncEntity(GetSVACID_Start_Entities(i), sm.Ef, sm2)
	}
	for i := range sm2.Plys {
		sm2.Plys[i] = GetNewSyncPlayer(GetSVACID_Start_OtherPlayer(i), sm.Ef, sm2)
	}
	sm2.ActivePlayer.Se.sm = sm2
	return
}
func (sm *SmallWorld) Clear() {
	sm.Ents = make([]*SyncEntity, SYNCENTITIES_PREP)
	sm.Plys = make([]*SyncPlayer, SYNCPLAYER_PREP)
	sm.ActivePlayer = GetNewSyncPlayer(GetSVACID_Start_OwnPlayer(), sm.Ef, sm)
}

type SmallWorld struct {
	X, Y, W, H             float64
	tile_path, struct_path string
	Ef                     *EntityFactory
	Ents                   []*SyncEntity
	Plys                   []*SyncPlayer

	ActivePlayer *SyncPlayer
	newPlayer    bool

	Struct *GE.WorldStructure

	ChanToClient, ChanToServer *GC.Channel
	World                      *World

	TimePerFrame        int64
	FrameChanSendPeriod int
	FrameCounter        *int
}

func (sm *SmallWorld) IsOnServer() bool {
	return sm.World != nil
}
func (sm *SmallWorld) SendToClient(idx int, msg []byte, force bool) {
	sm.ChanToClient.SendToPipe(idx, msg, force)
}
func (sm *SmallWorld) SendToServer(idx int, msg []byte, force bool) {
	sm.ChanToServer.SendToPipe(idx, msg, force)
}
func (sm *SmallWorld) ReceiveFromClient(sv GC.SyncVar, id int) {
	defer sm.ChanToClient.ResetJustChanged(SMALLWORLD_PLACEHOLDER_TOSERV, SMALLWORLD_PLACEHOLDER_TOSERV)
}
func (sm *SmallWorld) ReceiveFromServer(sv GC.SyncVar, id int) {
	defer sm.ChanToClient.ResetJustChanged(SMALLWORLD_WORLDSTRUCTURE_MSG, SMALLWORLD_SET_ACTIVEPLAYER_ID)
	if sm.ChanToClient.JustChanged(SMALLWORLD_WORLDSTRUCTURE_MSG) {
		data := sm.ChanToClient.Pipes[SMALLWORLD_WORLDSTRUCTURE_MSG]
		fmt.Println("WorldStructure received: ", data)
		sm.ChangeWorldStruct(data)
	}
	if sm.ChanToClient.JustChanged(SMALLWORLD_TIMEPERFRAME_MSG) {
		data := sm.ChanToClient.Pipes[SMALLWORLD_TIMEPERFRAME_MSG]
		fmt.Println("TimePerFrame received: ", data)
		sm.TimePerFrame = cmp.BytesToInt64(data[0:8])
	}
	if sm.ChanToClient.JustChanged(SMALLWORLD_FRAMEANDTIME_MSG) {
		data := sm.ChanToClient.Pipes[SMALLWORLD_FRAMEANDTIME_MSG]
		fmt.Println("FrameAndTime received: ", data)
		*sm.FrameCounter = int(cmp.BytesToInt64(data[0:8]))
		GE.ShitImDying(sm.Struct.CurrentTime.UnmarshalBinary(data[8:23]))
		sm.Struct.UpdateTime(time.Duration(sm.TimePerFrame))
	}
	if sm.ChanToClient.JustChanged(SMALLWORLD_SET_ACTIVEPLAYER_ID) {
		data := sm.ChanToClient.Pipes[SMALLWORLD_SET_ACTIVEPLAYER_ID]
		if sm.ActivePlayer.HasPlayer() {
			fmt.Println("New player ID received")
			sm.ActivePlayer.ID = cmp.BytesToInt16(data[0:2])
		} else {
			fmt.Println("Error: Cannot set ID of nil player")
		}
	}
}
func (sm *SmallWorld) SetWorldStruct(wS *GE.WorldStructure) error {
	if wS != nil {
		sm.Struct = wS
		bs := wS.ToBytes()
		sm.SendToClient(SMALLWORLD_WORLDSTRUCTURE_MSG, bs, true)
	}
	return nil
}
func (sm *SmallWorld) SetTimePerFrame(tpf int64) {
	sm.TimePerFrame = tpf
	sm.SendToClient(SMALLWORLD_TIMEPERFRAME_MSG, cmp.Int64ToBytes(sm.TimePerFrame), true)
}
func (sm *SmallWorld) SyncFrameAndTime() {
	bs := make([]byte, 23)
	copy(bs[0:8], cmp.Int64ToBytes(int64(*sm.FrameCounter)))
	timBs, err := sm.Struct.CurrentTime.MarshalBinary()
	GE.ShitImDying(err)
	copy(bs[8:23], timBs)
	sm.SendToClient(SMALLWORLD_FRAMEANDTIME_MSG, bs, true)
}
func (sm *SmallWorld) SetActivePlayerID(id int16) {
	sm.ActivePlayer.ID = id
	sm.SendToClient(SMALLWORLD_SET_ACTIVEPLAYER_ID, cmp.Int16ToBytes(id), true)
}
func (sm *SmallWorld) ChangeWorldStruct(data []byte) {
	if len(data) > 0 {
		wS, err := GE.LoadWorldStructureFromBytes(sm.X, sm.Y, sm.W, sm.H, data, sm.tile_path, sm.struct_path)
		if err != nil {
			panic(err)
		}
		sm.Struct = wS
	}
}

func (sm *SmallWorld) SetEntitiesFromChunks(chL []*Chunk, idxs ...int) {
	setted := make([]int, 0)
	for _, idx := range idxs {
		ents := chL[idx].GetEntities()
		for _, e := range ents {
			eidx := sm.HasEntity(e)
			if eidx < 0 {
				eidx = sm.GetIdxOfNilEnt()
				sm.Ents[eidx].SetEntity(e)
			}
			setted = append(setted, eidx)
		}
	}
	for i := range sm.Ents {
		if containsI(setted, i) == -1 {
			sm.Ents[i].SetEntity(nil)
		}
	}
}
func (sm *SmallWorld) GetIdxOfNilEnt() int {
	for i, se := range sm.Ents {
		if !se.HasEntity() {
			return i
		}
	}
	return 0
}
func (sm *SmallWorld) HasEntityWithID(id int16) *SyncEntity {
	for _, se := range sm.Ents {
		if se.HasEntity() && se.Entity.ID == id {
			return se
		}
	}
	return nil
}
func (sm *SmallWorld) HasEntity(e *Entity) int {
	for i, se := range sm.Ents {
		if se.HasEntity() && se.Entity == e {
			return i
		}
	}
	return -1
}
func (sm *SmallWorld) UpdateAll(server bool) {
	if !server {
		*sm.FrameCounter++
		sm.Struct.UpdateTime(time.Duration(sm.TimePerFrame))
	}
	if sm.ActivePlayer.HasPlayer() {
		sm.ActivePlayer.UpdateAll(sm, server, sm.Struct.Collides)
	}
	for _, pl := range sm.Plys {
		if pl.HasPlayer() {
			pl.Player.UpdateAll(sm, server, sm.Struct.Collides)
		}
	}
	for _, ent := range sm.Ents {
		if ent.HasEntity() {
			ent.Entity.UpdateAll(sm, server, sm.Struct.Collides)
		}
	}
}
func (sm *SmallWorld) ResetActions() {
	if sm.ActivePlayer.HasPlayer() {
		sm.ActivePlayer.Player.Actions().Reset()
	}
	for _, pl := range sm.Plys {
		if pl.HasPlayer() {
			pl.Player.Actions().Reset()
		}
	}
	for _, ent := range sm.Ents {
		if ent.HasEntity() {
			ent.Entity.Actions().Reset()
		}
	}
}
func (sm *SmallWorld) Print(ents bool) (out string, c int) {
	out = fmt.Sprintf("%v: ", *sm.FrameCounter)
	if sm.ActivePlayer.HasPlayer() {
		x, y, _ := sm.ActivePlayer.Player.GetPos()
		out += fmt.Sprintf("(AP)(%p)(%v)|%0.2f, %0.2f, %s|", sm.ActivePlayer.Player, sm.ActivePlayer.ID, x, y, sm.ActivePlayer.Player.Entity.Actions().Print())
		c++
	}
	for _, pl := range sm.Plys {
		if pl.HasPlayer() {
			x, y, _ := pl.Player.GetPos()
			out += fmt.Sprintf("(P)(%p)(%v)|%0.2f, %0.2f, %s|", pl.Player, pl.ID, x, y, pl.Player.Entity.Actions().Print())
			c++
		}
	}
	if ents {
		for _, ent := range sm.Ents {
			if ent.HasEntity() {
				x, y, _ := ent.Entity.GetPos()
				out += fmt.Sprintf("(E)(%p)(%v)|%0.2f, %0.2f, %s|", ent.Entity, ent.Entity.ID, x, y, ent.Entity.Actions().Print())
				c++
			}
		}
	}
	return
}
func (sm *SmallWorld) HasNewActivePlayer() (bool, *Player) {
	if sm.newPlayer {
		sm.newPlayer = false
		if sm.ActivePlayer.HasPlayer() {
			return true, sm.ActivePlayer.Player
		}
	}
	return false, nil
}
func (sm *SmallWorld) GetSyncPlayersFromWorld(w *World) {
	idx := 0
	for _, pl := range w.Players {
		if pl != sm.ActivePlayer.Player {
			if pl != sm.Plys[idx].Player {
				sm.Plys[idx].SetPlayer(pl)
			}
			idx++
		}
	}
	for i := idx; i < len(sm.Plys); i++ {
		sm.Plys[idx].SetPlayer(nil)
	}
}
func (sm *SmallWorld) HasWorldStruct() bool {
	return sm.Struct != nil
}
func (sm *SmallWorld) ReassignAllEntities() {
	if sm.HasWorldStruct() {
		sm.Struct.Add_Drawables.Clear()
		for _, e := range sm.Ents {
			if e.HasEntity() {
				sm.Struct.Add_Drawables.Add(e.Entity)
			}
		}
		for _, p := range sm.Plys {
			if p.HasPlayer() {
				sm.Struct.Add_Drawables.Add(p.Player)
			}
		}
		if sm.ActivePlayer.HasPlayer() {
			sm.Struct.Add_Drawables.Add(sm.ActivePlayer.Player)
		}
	}
}
func (sm *SmallWorld) StandardOnEntityChange(se interface{}, oldE, newE GE.Drawable) {
	if sm.HasWorldStruct() {
		if !IsInterfaceNil(oldE) {
			err, dp := sm.Struct.Add_Drawables.Remove(oldE)
			if err == nil {
				sm.Struct.Add_Drawables = dp
			} else {
				fmt.Printf("Error removing %p: %v\n", oldE, err)
			}
		}
		if !IsInterfaceNil(newE) {
			sm.Struct.Add_Drawables = sm.Struct.Add_Drawables.Add(newE)
		}
	}
}
func IsInterfaceNil(i interface{}) bool {
	pnt := fmt.Sprintf("%p", i)
	return pnt == "0x0" || i == nil
}
func (sm *SmallWorld) OnActivePlayerChange(se interface{}, oldE, newE GE.Drawable) {
	sm.newPlayer = true
	sm.StandardOnEntityChange(se, oldE, newE)
}
func (sm *SmallWorld) RegisterOnEntityChangeListeners() {
	for _, e := range sm.Ents {
		e.OnNewEntity = sm.StandardOnEntityChange
	}
	for _, p := range sm.Plys {
		p.OnNewPlayer = sm.StandardOnEntityChange
	}
	sm.ActivePlayer.OnNewPlayer = sm.OnActivePlayerChange
}
func (sm *SmallWorld) Register(m *GC.ServerManager, client *ws.Conn) {
	AllSVs := make(map[int]GC.SyncVar)

	sm.ActivePlayer.GetSyncVars(AllSVs)
	for _, e := range sm.Ents {
		e.GetSyncVars(AllSVs)
		//e.RegisterOnChange(m)
	}
	for _, p := range sm.Plys {
		p.GetSyncVars(AllSVs)
		//p.RegisterOnChange(m)
	}
	AllSVs[SM_TO_CLIENT] = sm.ChanToClient
	AllSVs[SM_TO_SERVER] = sm.ChanToServer
	m.RegisterOnChangeFunc(SM_TO_SERVER, []func(GC.SyncVar, int){sm.ReceiveFromClient}, client)

	m.RegisterSyncVars(true, AllSVs, client)
	m.Server.WaitForConfirmation(client)
	sm.ActivePlayer.RegisterOnChange(m.Handler[client])
	sm.ActivePlayer.OnNewPlayer = sm.OnActivePlayerChange
}
func (sm *SmallWorld) GetRegistered(m *GC.ClientManager) {
	sm.ActivePlayer.GetRegisterdSyncVars(m)
	for _, e := range sm.Ents {
		e.GetRegisterdSyncVars(m)
		e.RegisterOnChange(m)
	}
	for _, p := range sm.Plys {
		p.GetRegisterdSyncVars(m)
		p.RegisterOnChange(m)
	}
	sm.ChanToClient = m.SyncvarsByACID[SM_TO_CLIENT].(*GC.Channel)
	sm.ChanToServer = m.SyncvarsByACID[SM_TO_SERVER].(*GC.Channel)
	m.RegisterOnChangeFunc(SM_TO_CLIENT, sm.ReceiveFromServer)
}
func (sm *SmallWorld) Draw(screen *ebiten.Image) {
	sm.ActivePlayer.MoveWorld(sm.Struct)
	sm.Struct.UpdateObjDrawables()
	sm.Struct.Draw(screen)
}
func (sm *SmallWorld) UpdateVars() {
	for _, e := range sm.Ents {
		e.UpdateChanFromEnt()
	}
	for _, p := range sm.Plys {
		p.UpdateChanFromPlayer()
	}
	if sm.HasWorldStruct() && *sm.FrameCounter%sm.FrameChanSendPeriod == 0 {
		sm.SyncFrameAndTime()
	}
}
