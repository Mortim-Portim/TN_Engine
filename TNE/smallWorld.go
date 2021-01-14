package TNE

import (
	ws "github.com/gorilla/websocket"
	"github.com/mortim-portim/GraphEng/GE"
	"github.com/mortim-portim/GameConn/GC"
	"github.com/hajimehoshi/ebiten"
)

/**
Channel communication:
(1)[mt]+(n)[msg]
**/
/**
Dont forget do set FrameCounter on Server and assign WorldStructure
**/

func GetSmallWorld(X, Y, W, H float64, tile_path, struct_path, entity_path string) (sm *SmallWorld, err error) {
	fc := 0
	ef, err := GetEntityFactory(entity_path, &fc, 3)
	sm = &SmallWorld{Ents:make([]*SyncEntity, SYNCENTITIES_PREP),
					 Plys:make([]*SyncPlayer, SYNCPLAYER_PREP),
					 SyncFrame:GC.CreateSyncInt64(0),
					 SyncLightLevel:GC.CreateSyncInt16(0),
					 WorldChan:GC.CreateSyncString(""),
					 Ef:ef,X:X,Y:Y,W:W,H:H,tile_path:tile_path,struct_path:struct_path,
					 FrameCounter:&fc,
					 ActivePlayer:GetNewSyncPlayer(GetSVACID_Start_OwnPlayer(), ef),
	}
	for i,_ := range(sm.Ents) {
		sm.Ents[i] = GetNewSyncEntity(GetSVACID_Start_Entities(i), ef)
	}
	for i,_ := range(sm.Plys) {
		sm.Plys[i] = GetNewSyncPlayer(GetSVACID_Start_OtherPlayer(i), ef)
	}
	return
}
type SmallWorld struct {
	X,Y,W,H float64
	tile_path, struct_path string
	Ef *EntityFactory
	Ents []*SyncEntity
	Plys []*SyncPlayer
	
	ActivePlayer *SyncPlayer
	
	Struct *GE.WorldStructure
	
	SyncFrame *GC.SyncInt64
	SyncLightLevel *GC.SyncInt16
	WorldChan *GC.SyncString
	
	FrameCounter *int
}
func (sm *SmallWorld) SetWorldStruct(wS *GE.WorldStructure) error {
	if wS != nil {
		sm.Struct = wS
		bs,err := wS.ToBytes()
		if err != nil {
			return err
		}
		sm.WorldChan.SetBs(bs)
	}
	return nil
}
func (sm *SmallWorld) Draw(screen *ebiten.Image) {
	sm.ActivePlayer.MoveWorld(sm.Struct)
	sm.Struct.Draw(screen)
}
func (sm *SmallWorld) UpdateVars() {
	for _,e := range(sm.Ents) {
		if e.HasEntity() {
			e.UpdateVarsFromEnt()
		}
	}
	for _,p := range(sm.Plys) {
		if p.HasPlayer() {
			p.UpdateVarsFromPlayer()
		}
	}
	if sm.HasWorldStruct() {
		sm.SyncFrame.SetInt(int64(*sm.FrameCounter))
		sm.SyncLightLevel.SetInt(sm.Struct.GetLightLevel())
	}
}
func (sm *SmallWorld) HasWorldStruct() bool {
	return sm.Struct != nil
}
func (sm *SmallWorld) ReassignAllEntities() {
	if sm.HasWorldStruct() {
		sm.Struct.Add_Drawables.Clear()
		for _,e := range(sm.Ents) {
			if e.HasEntity() {
				sm.Struct.Add_Drawables.Add(e.Entity)
			}
		}
		for _,p := range(sm.Plys) {
			if p.HasPlayer() {
				sm.Struct.Add_Drawables.Add(p.Player)
			}
		}
		if sm.ActivePlayer.HasPlayer() {
			sm.Struct.Add_Drawables.Add(sm.ActivePlayer.Player)
		}
	}
}
func (sm *SmallWorld) StandardOnEntityChange(se interface{}, oldE, newE GE.Drawable){
	if sm.HasWorldStruct() {
		if oldE != nil {
			sm.Struct.Add_Drawables.Remove(oldE)
		}
		if newE != nil {
			sm.Struct.Add_Drawables.Add(newE)
		}
	}
}
func (sm *SmallWorld) RegisterOnEntityChangeListeners() {
	for _,e := range(sm.Ents) {
		e.OnNewEntity = sm.StandardOnEntityChange
	}
	for _,p := range(sm.Plys) {
		p.OnNewPlayer = sm.StandardOnEntityChange
	}
	sm.ActivePlayer.OnNewPlayer = sm.StandardOnEntityChange
}
func (sm *SmallWorld) OnWorldChanChange(sv GC.SyncVar, id int) {
	data := sm.WorldChan.GetBs()
	wS,err := GE.GetWorldStructureFromBytes(sm.X, sm.Y, sm.W, sm.H, data, sm.tile_path, sm.struct_path)
	if err != nil {
		panic(err)
	}
	sm.Struct = wS
}
func (sm *SmallWorld) OnFrameChange(sv GC.SyncVar, id int) {
	*sm.FrameCounter = int(sm.SyncFrame.GetInt())
}
func (sm *SmallWorld) OnLightLevelChange(sv GC.SyncVar, id int) {
	sm.Struct.SetLightLevel(sm.SyncLightLevel.GetInt())
}
func (sm *SmallWorld) Register(m *GC.ServerManager, clients ...*ws.Conn) {
	sm.ActivePlayer.RegisterSyncVars(m, clients...)
	for _,e := range(sm.Ents) {
		e.RegisterSyncVars(m, clients...)
		//e.RegisterOnChange(m)
	}
	for _,p := range(sm.Plys) {
		p.RegisterSyncVars(m, clients...)
		//p.RegisterOnChange(m)
	}
	m.RegisterSyncVar(sm.SyncFrame, WorldFrameChan_ACID, clients...)
	m.RegisterSyncVar(sm.SyncLightLevel, WorldLightLevelChan_ACID, clients...)
	m.RegisterSyncVar(sm.WorldChan, WorldStructChan_ACID, clients...)
	//m.RegisterOnChangeFunc(WorldChannel_ACID, []func(GC.SyncVar, int){sm.OnChannelChange}, clients...)
}
func (sm *SmallWorld) GetRegistered(m *GC.ClientManager) {
	sm.ActivePlayer.GetRegisterdSyncVars(m)
	for _,e := range(sm.Ents) {
		e.GetRegisterdSyncVars(m)
		e.RegisterOnChange(m)
	}
	for _,p := range(sm.Plys) {
		p.GetRegisterdSyncVars(m)
		p.RegisterOnChange(m)
	}
	sm.SyncFrame = m.SyncvarsByACID[WorldFrameChan_ACID].(*GC.SyncInt64)
	m.RegisterOnChangeFunc(WorldFrameChan_ACID, sm.OnFrameChange)
	
	sm.SyncLightLevel = m.SyncvarsByACID[WorldLightLevelChan_ACID].(*GC.SyncInt16)
	m.RegisterOnChangeFunc(WorldLightLevelChan_ACID, sm.OnLightLevelChange)
	
	sm.WorldChan = m.SyncvarsByACID[WorldStructChan_ACID].(*GC.SyncString)
	m.RegisterOnChangeFunc(WorldStructChan_ACID, sm.OnWorldChanChange)
}