package TNE

import (
	"github.com/mortim-portim/GameConn/GC"
)

const NumberOfSVACIDs_Msg = "SVACIDs"

const (
	SYNCENT_CHAN_ENTITY_CREATION = iota
	SYNCENT_CHAN_ACTIONS

	SYNCENT_CHAN_PLAYER_CREATION

	SYNCENT_CHAN_NUM
)

const (
	//amount of syncVars needed by one Entity
	SYNCVARS_PER_ENTITY = 1
	//count of SyncEntities to be prepared
	SYNCENTITIES_PREP = 100
	//amount of syncVars needed by one Player
	SYNCVARS_PER_PLAYER = SYNCVARS_PER_ENTITY
	//count of SyncPlayer to be prepared besides the own player
	SYNCPLAYER_PREP = 10
)
const (
	//SyncVars that are registered by the world
	WorldStructChan_ACID = iota
	WorldFrameChan_ACID
	WORLD_SYNCVARS
)

//Returns the startACID for the own player
func GetSVACID_Start_OwnPlayer() int {
	return WORLD_SYNCVARS
}

//Returns the startACID for the other player
func GetSVACID_Start_OtherPlayer(idx int) int {
	return GetSVACID_Start_OwnPlayer() + SYNCVARS_PER_PLAYER + idx*SYNCVARS_PER_PLAYER
}

//Returns the startACID for the entity
func GetSVACID_Start_Entities(idx int) int {
	return GetSVACID_Start_OtherPlayer(SYNCPLAYER_PREP) + idx*SYNCVARS_PER_ENTITY
}
func GetSVACID_Count() int {
	return GetSVACID_Start_Entities(SYNCENTITIES_PREP)
}

//Initializes the SyncClient
func InitialSyncClient() {
	GC.InitSyncVarStandardTypes()
}

const (
	ERR_ENTITY_IS_NIL         = "Entity: %v is nil"
	ERR_NO_FACTORY_FOR_ENTITY = "No factory for Entity: %v, with fcID: %v"
)
