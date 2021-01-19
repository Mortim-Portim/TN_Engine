package TNE

import (
	ws "github.com/gorilla/websocket"
	"github.com/mortim-portim/GraphEng/GE"
	"github.com/mortim-portim/GameConn/GC"
	"github.com/hajimehoshi/ebiten"
	
	"fmt"
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
func (sm *SmallWorld) New() (sm2 *SmallWorld) {
	sm2 = &SmallWorld{Ents:make([]*SyncEntity, SYNCENTITIES_PREP),
					 Plys:make([]*SyncPlayer, SYNCPLAYER_PREP),
					 SyncFrame:GC.CreateSyncInt64(0),
					 SyncLightLevel:GC.CreateSyncInt16(0),
					 WorldChan:GC.CreateSyncString(""),
					 Ef:sm.Ef,X:sm.X,Y:sm.Y,W:sm.W,H:sm.H,tile_path:sm.tile_path,struct_path:sm.struct_path,
					 FrameCounter:sm.FrameCounter,
					 ActivePlayer:GetNewSyncPlayer(GetSVACID_Start_OwnPlayer(), sm.Ef),
					 Struct:sm.Struct,
	}
	for i,_ := range(sm2.Ents) {
		sm2.Ents[i] = GetNewSyncEntity(GetSVACID_Start_Entities(i), sm.Ef)
	}
	for i,_ := range(sm2.Plys) {
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
	X,Y,W,H float64
	tile_path, struct_path string
	Ef *EntityFactory
	Ents []*SyncEntity
	Plys []*SyncPlayer
	
	ActivePlayer *SyncPlayer
	newPlayer bool
	
	Struct *GE.WorldStructure
	
	SyncFrame *GC.SyncInt64
	SyncLightLevel *GC.SyncInt16
	WorldChan *GC.SyncString
	
	FrameCounter *int
}
func (sm *SmallWorld) UpdateAll() {
	if sm.ActivePlayer.HasPlayer() {
		sm.ActivePlayer.UpdateAll(nil)
	}
	for _,pl := range(sm.Plys) {
		if pl.HasPlayer() {
			pl.Player.UpdateAll(nil)
		}
	}
	for _,ent := range(sm.Ents) {
		if ent.HasEntity() {
			ent.Entity.UpdateAll(nil)
		}
	}
}
func (sm *SmallWorld) Print() (out string) {
	out = fmt.Sprintf("SHOULD print information about the smallWorld")
	return
}
func (sm *SmallWorld) HasNewActivePlayer() (bool, *Player) {
	if sm.newPlayer {
		sm.newPlayer = false
		if sm.ActivePlayer.Player != nil {
			return true, sm.ActivePlayer.Player
		}
	}
	return false, nil
}
func (sm *SmallWorld) GetSyncPlayersFromWorld(w *World) {
	idx := 0
	for _,pl := range(w.Players) {
		if pl != sm.ActivePlayer.Player {
			if pl != sm.Plys[idx].Player {
				sm.Plys[idx].SetPlayer(pl)
			}
			idx ++
		}
	}
	for i := idx; i < len(sm.Plys); i++ {
		sm.Plys[idx].SetPlayer(nil)
	}
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
	sm.Struct.UpdateObjDrawables()
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
	fmt.Printf("Changing Entity: %p\n", se)
	if sm.HasWorldStruct() {
		if fmt.Sprintf("%p", oldE) != "0x0" {
			fmt.Printf("Removing %p\n", oldE)
			err, dp := sm.Struct.Add_Drawables.Remove(oldE)
			if err == nil {
				sm.Struct.Add_Drawables = dp
			}else{
				fmt.Println(err)
			}
			fmt.Println(([]GE.Drawable)(*sm.Struct.Add_Drawables))
		}
		if fmt.Sprintf("%p", newE) != "0x0" {
			fmt.Printf("Adding %p\n", newE)
			sm.Struct.Add_Drawables.Add(newE)
			fmt.Println(([]GE.Drawable)(*sm.Struct.Add_Drawables))
		}
	}
}
func (sm *SmallWorld) OnActivePlayerChange(se interface{}, oldE, newE GE.Drawable) {
	sm.newPlayer = true
	fmt.Printf("Changing Player: %p\n", se)
	sm.StandardOnEntityChange(se, oldE, newE)
}
func (sm *SmallWorld) RegisterOnEntityChangeListeners() {
	for _,e := range(sm.Ents) {
		e.OnNewEntity = sm.StandardOnEntityChange
	}
	for _,p := range(sm.Plys) {
		p.OnNewPlayer = sm.StandardOnEntityChange
	}
	sm.ActivePlayer.OnNewPlayer = sm.OnActivePlayerChange
}
func (sm *SmallWorld) OnWorldChanChange(sv GC.SyncVar, id int) {
	data := sm.WorldChan.GetBs()
	if len(data) > 0 {
		wS,err := GE.GetWorldStructureFromBytes(sm.X, sm.Y, sm.W, sm.H, data, sm.tile_path, sm.struct_path)
		if err != nil {
			panic(err)
		}
		sm.Struct = wS
	}
}
func (sm *SmallWorld) OnFrameChange(sv GC.SyncVar, id int) {
	*sm.FrameCounter = int(sm.SyncFrame.GetInt())
	//fmt.Println("FrameCounter Change: ", *sm.FrameCounter)
}
func (sm *SmallWorld) OnLightLevelChange(sv GC.SyncVar, id int) {
	if sm.HasWorldStruct() {
		sm.Struct.SetLightLevel(sm.SyncLightLevel.GetInt())
		//fmt.Println("LightLevel Change: ", sm.SyncLightLevel.GetInt())
	}
}
func (sm *SmallWorld) Register(m *GC.ServerManager, client *ws.Conn) {
	AllSVs := make(map[int]GC.SyncVar)
	
	sm.ActivePlayer.GetSyncVars(AllSVs)
	for _,e := range(sm.Ents) {
		e.GetSyncVars(AllSVs)
		//e.RegisterOnChange(m)
	}
	for _,p := range(sm.Plys) {
		p.GetSyncVars(AllSVs)
		//p.RegisterOnChange(m)
	}
	AllSVs[WorldFrameChan_ACID] = sm.SyncFrame
	AllSVs[WorldLightLevelChan_ACID] = sm.SyncLightLevel
	AllSVs[WorldStructChan_ACID] = sm.WorldChan
	//m.RegisterOnChangeFunc(WorldChannel_ACID, []func(GC.SyncVar, int){sm.OnChannelChange}, clients...)
	m.RegisterSyncVars(AllSVs, client)
	sm.ActivePlayer.RegisterOnChange(m.Handler[client])
	sm.ActivePlayer.OnNewPlayer = sm.OnActivePlayerChange
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