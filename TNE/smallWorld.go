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

func GetSmallWorld(X, Y, W, H float64, tile_path, struct_path, entity_path string) (sm *SmallWorld, err error) {
	fc := 0
	ef, err := GetEntityFactory(entity_path, &fc, 3)
	sm = &SmallWorld{Ents: make([]*SyncEntity, SYNCENTITIES_PREP),
		Plys:          make([]*SyncPlayer, SYNCPLAYER_PREP),
		SyncFrameChan: GC.CreateSyncString(""),
		WorldChan:     GC.CreateSyncString(""),
		Ef:            ef, X: X, Y: Y, W: W, H: H, tile_path: tile_path, struct_path: struct_path,
		FrameCounter: &fc,
		ActivePlayer: GetNewSyncPlayer(GetSVACID_Start_OwnPlayer(), ef),
	}
	for i := range sm.Ents {
		sm.Ents[i] = GetNewSyncEntity(GetSVACID_Start_Entities(i), ef)
	}
	for i := range sm.Plys {
		sm.Plys[i] = GetNewSyncPlayer(GetSVACID_Start_OtherPlayer(i), ef)
	}
	return
}
func (sm *SmallWorld) New() (sm2 *SmallWorld) {
	sm2 = &SmallWorld{Ents: make([]*SyncEntity, SYNCENTITIES_PREP),
		Plys:          make([]*SyncPlayer, SYNCPLAYER_PREP),
		SyncFrameChan: GC.CreateSyncString(""),
		WorldChan:     GC.CreateSyncString(""),
		Ef:            sm.Ef, X: sm.X, Y: sm.Y, W: sm.W, H: sm.H, tile_path: sm.tile_path, struct_path: sm.struct_path,
		FrameCounter:        sm.FrameCounter,
		ActivePlayer:        GetNewSyncPlayer(GetSVACID_Start_OwnPlayer(), sm.Ef),
		Struct:              sm.Struct,
		FrameChanSendPeriod: sm.FrameChanSendPeriod,
		TimePerFrame:        sm.TimePerFrame,
	}
	for i := range sm2.Ents {
		sm2.Ents[i] = GetNewSyncEntity(GetSVACID_Start_Entities(i), sm.Ef)
	}
	for i := range sm2.Plys {
		sm2.Plys[i] = GetNewSyncPlayer(GetSVACID_Start_OtherPlayer(i), sm.Ef)
	}
	return
}
func (sm *SmallWorld) Clear() {
	sm.Ents = make([]*SyncEntity, SYNCENTITIES_PREP)
	sm.Plys = make([]*SyncPlayer, SYNCPLAYER_PREP)
	sm.ActivePlayer = GetNewSyncPlayer(GetSVACID_Start_OwnPlayer(), sm.Ef)
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

	SyncFrameChan *GC.SyncString
	WorldChan     *GC.SyncString

	TimePerFrame        int64
	FrameChanSendPeriod int
	FrameCounter        *int
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
		if !containsI(setted, i) {
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
		sm.ActivePlayer.UpdateAll(nil, server, sm.Struct.Collides)
	}
	for _, pl := range sm.Plys {
		if pl.HasPlayer() {
			pl.Player.UpdateAll(nil, server, sm.Struct.Collides)
		}
	}
	for _, ent := range sm.Ents {
		if ent.HasEntity() {
			ent.Entity.UpdateAll(nil, server, sm.Struct.Collides)
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
		out += fmt.Sprintf("(AP)(%p)|%0.2f, %0.2f, %s|", sm.ActivePlayer.Player, x, y, sm.ActivePlayer.Player.Entity.Actions().Print())
		c++
	}
	for _, pl := range sm.Plys {
		if pl.HasPlayer() {
			x, y, _ := pl.Player.GetPos()
			out += fmt.Sprintf("(P)(%p)|%0.2f, %0.2f, %s|", pl.Player, x, y, pl.Player.Entity.Actions().Print())
			c++
		}
	}
	if ents {
		for _, ent := range sm.Ents {
			if ent.HasEntity() {
				x, y, _ := ent.Entity.GetPos()
				out += fmt.Sprintf("(E)(%p)|%0.2f, %0.2f, %s|", ent.Entity, x, y, ent.Entity.Actions().Print())
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
func (sm *SmallWorld) SetWorldStruct(wS *GE.WorldStructure) error {
	if wS != nil {
		sm.Struct = wS
		bs := wS.ToBytes()
		sm.WorldChan.SetBs(bs)
	}
	return nil
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
		sm.SetFrameChangeVar()
	}
}
func (sm *SmallWorld) SetTimePerFrame(tpf int64) {
	sm.TimePerFrame = tpf
	bs := make([]byte, 9)
	bs[0] = 0
	copy(bs[1:9], cmp.Int64ToBytes(sm.TimePerFrame))
	sm.SyncFrameChan.SetBs(bs)
}
func (sm *SmallWorld) SetFrameChangeVar() {
	bs := make([]byte, 24)
	bs[0] = 1
	copy(bs[1:9], cmp.Int64ToBytes(int64(*sm.FrameCounter)))
	timBs, err := sm.Struct.CurrentTime.MarshalBinary()
	GE.ShitImDying(err)
	copy(bs[9:24], timBs)
	sm.SyncFrameChan.SetBs(bs)
}
func (sm *SmallWorld) GetFrameChangeVar(bs []byte) {
	if len(bs) < 9 {
		return
	}
	idx := bs[0]
	bs = bs[1:]
	if idx == 0 {
		sm.TimePerFrame = cmp.BytesToInt64(bs[0:8])
	} else if idx == 1 {
		*sm.FrameCounter = int(cmp.BytesToInt64(bs[0:8]))
		GE.ShitImDying(sm.Struct.CurrentTime.UnmarshalBinary(bs[8:23]))
		sm.Struct.UpdateTime(time.Duration(sm.TimePerFrame))
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
func (sm *SmallWorld) OnWorldChanChange(sv GC.SyncVar, id int) {
	data := sm.WorldChan.GetBs()
	if len(data) > 0 {
		wS, err := GE.LoadWorldStructureFromBytes(sm.X, sm.Y, sm.W, sm.H, data, sm.tile_path, sm.struct_path)
		if err != nil {
			panic(err)
		}
		sm.Struct = wS
	}
}
func (sm *SmallWorld) OnSyncFrameChanChange(sv GC.SyncVar, id int) {
	sm.GetFrameChangeVar(sm.SyncFrameChan.GetBs())
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
	AllSVs[WorldFrameChan_ACID] = sm.SyncFrameChan
	AllSVs[WorldStructChan_ACID] = sm.WorldChan
	//m.RegisterOnChangeFunc(WorldChannel_ACID, []func(GC.SyncVar, int){sm.OnChannelChange}, clients...)
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
	sm.SyncFrameChan = m.SyncvarsByACID[WorldFrameChan_ACID].(*GC.SyncString)
	m.RegisterOnChangeFunc(WorldFrameChan_ACID, sm.OnSyncFrameChanChange)

	sm.WorldChan = m.SyncvarsByACID[WorldStructChan_ACID].(*GC.SyncString)
	m.RegisterOnChangeFunc(WorldStructChan_ACID, sm.OnWorldChanChange)
}
