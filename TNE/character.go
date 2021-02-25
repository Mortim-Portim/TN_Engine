package TNE

import (
	"fmt"
	"math"

	cmp "github.com/mortim-portim/GraphEng/compression"
)

/**
Class and Race represent DnD similar Classes and Races
Character represents attributes, class, race, name of Player
**/

const ERROR_WRONG_CHAR_VERSION = "Wrong character version: %v"
const NO_CHARACTER_DATA = "No character data supplied"

var Classes []*Class = []*Class{
	{"Fighter", 0, []string{"Berserker", "Defender"}},
	{"Knight", 1, []string{"Heavy Armor", "Combat Rager", "Guard", "Paladin"}},
	{"Wizard", 2, []string{"Test"}},
	{"Bard", 3, []string{"Test"}},
	{"Cleric", 4, []string{"Test"}},
	{"Druid", 5, []string{"Test"}},
	{"Ranger", 6, []string{"Test"}},
	{"Rogue", 7, []string{"Test"}},
	{"Sorcerer", 8, []string{"Test"}},
}
var Races []*Race = []*Race{
	{"Elv", 0, []int8{0, 2, 2, 0}, []int8{SCORE_NATURE, SCORE_ACROBATICS, SCORE_PERCEPTION, SCORE_INSIGHT, SCORE_INTELLIGENCE}, 1, []string{"Wood Elv", "High Elv"}},
	{"Human", 1, []int8{1, 1, 1, 1}, []int8{}, 4, []string{"Urban Human", "Country-Side Human", "Mountain Tribe"}},
	{"Half-Elv", 2, []int8{0, 2, 0, 2}, []int8{SCORE_ACROBATICS, SCORE_PERCEPTION, SCORE_INSIGHT, SCORE_DUNGEONEERING}, 1, []string{"Dark Elv", "City Elv"}},
	{"Ork", 3, []int8{3, 0, 0, 0}, []int8{SCORE_STRENGTH, SCORE_ENDURANCE}, 1, []string{"Mountain Ork", "Cave Ork"}},
	{"Goblin", 4, []int8{-1, 2, 0, 2}, []int8{SCORE_STEALTH, SCORE_THIEVERY, SCORE_ACROBATICS, SCORE_DECEPTION, SCORE_PERCEPTION}, 1, []string{"Ravin Goblin", "Sever Goblin"}},
	{"Dwarf", 5, []int8{2, 0, 2, 0}, []int8{SCORE_STRENGTH, SCORE_CRAFTSMANSHIP, SCORE_DUNGEONEERING}, 1, []string{"Hill Dwarf", "Mountain Dwarf"}},
	{"Halfling", 6, []int8{-1, 0, 1, 3}, []int8{SCORE_PERSUASION, SCORE_DECEPTION, SCORE_DUNGEONEERING}, 1, []string{"Rock Halfling", "Forest Halfling"}},
}

const (
	ABIL_STRENGTH = iota
	ABIL_DEXTERITY
	ABIL_INTELLIGENCE
	ABIL_CHARISMA

	ABIL_COUNT
)

const (
	SCORE_STRENGTH = iota
	SCORE_DEXTERITY
	SCORE_INTELLIGENCE
	SCORE_CHARISMA
	SCORE_ENDURANCE
	SCORE_PERSUASION
	SCORE_DECEPTION
	SCORE_PERFORMANCE
	SCORE_INSIGHT
	SCORE_THIEVERY
	SCORE_STEALTH
	SCORE_ACROBATICS
	SCORE_NATURE
	SCORE_ARCANA
	SCORE_PERCEPTION
	SCORE_CRAFTSMANSHIP
	SCORE_DUNGEONEERING

	SCORE_COUNT
)

type Class struct {
	Name     string
	id       int
	Subclass []string //will change later, placeholder
}

//const MAX_RACE_NAME_LENGTH = 20
type Race struct {
	Name          string
	id            int
	Attributes    []int8
	Proficiencies []int8
	Extraprof     int
	Subraces      []string //will change later, placeholder
}

const MAX_CHARACTER_NAME_LENGTH = 20

type Character struct {
	Name          string
	Class         *Class
	Race          *Race
	Attributes    []int8
	Proficiencies []int8
	Attacks       []byte
}

func (char *Character) GetAttack(idx int) int {
	if idx >= 0 && idx < len(char.Attacks) {
		return int(char.Attacks[idx])
	}
	return -1
}

func baseCompFunc(y float64, percents ...float64) func(vals ...float64) float64 {
	return func(vals ...float64) (f float64) {
		for i, val := range vals {
			f += math.Pow(percents[i], val)
		}
		f *= y
		return
	}
}

const (
	MODIFIER_SPEED_BASE = 5
	MODIFIER_SPEED_STR  = 1.03
	MODIFIER_SPEED_DEX  = 1.05

	MODIFIER_HEALTH_BASE = 50
	MODIFIER_HEALTH_STR  = 1.1
	MODIFIER_HEALTH_DEX  = 1.05

	MODIFIER_STAMINA_BASE = 50
	MODIFIER_STAMINA_STR  = 1.05
	MODIFIER_STAMINA_DEX  = 1.1

	MODIFIER_MANA_BASE = 50
	MODIFIER_MANA_INT  = 1.1
	MODIFIER_MANA_CHA  = 1.05
)

func (char *Character) SetEntityValues(e *Entity) {
	e.Speed = char.CompSpeed(MODIFIER_SPEED_BASE, MODIFIER_SPEED_STR, MODIFIER_SPEED_DEX)
	e.SetMaxHealth(float32(char.CompHealth(MODIFIER_HEALTH_BASE, MODIFIER_HEALTH_STR, MODIFIER_HEALTH_DEX)))
	e.SetMaxStamina(float32(char.CompStamina(MODIFIER_STAMINA_BASE, MODIFIER_STAMINA_STR, MODIFIER_STAMINA_DEX)))
	e.SetMaxMana(float32(char.CompMana(MODIFIER_MANA_BASE, MODIFIER_MANA_INT, MODIFIER_MANA_CHA)))
	e.ResetHSM()
}

func (char *Character) CompSpeed(base, strWeight, dexWeight float64) float64 {
	return baseCompFunc(base, strWeight, dexWeight)(float64(char.Proficiencies[SCORE_STRENGTH]), float64(char.Proficiencies[SCORE_DEXTERITY]))
}
func (char *Character) CompHealth(base, strWeight, dexWeight float64) float64 {
	return baseCompFunc(base, strWeight, dexWeight)(float64(char.Proficiencies[SCORE_STRENGTH]), float64(char.Proficiencies[SCORE_DEXTERITY]))
}
func (char *Character) CompStamina(base, strWeight, dexWeight float64) float64 {
	return baseCompFunc(base, strWeight, dexWeight)(float64(char.Proficiencies[SCORE_STRENGTH]), float64(char.Proficiencies[SCORE_DEXTERITY]))
}
func (char *Character) CompMana(base, intWeight, chaWeight float64) float64 {
	return baseCompFunc(base, intWeight, chaWeight)(float64(char.Proficiencies[SCORE_INTELLIGENCE]), float64(char.Proficiencies[SCORE_CHARISMA]))
}
func GetCharacter(name string, raceId, classId int) (char *Character) {
	char = &Character{
		Name:          name,
		Race:          Races[raceId],
		Class:         Classes[classId],
		Attributes:    make([]int8, len(Races[raceId].Attributes)),
		Proficiencies: make([]int8, len(Races[raceId].Proficiencies)),
	}
	copy(char.Attributes, char.Race.Attributes)
	copy(char.Proficiencies, char.Race.Proficiencies)
	return
}
func (char *Character) Copy() *Character {
	if char == nil {
		return nil
	}
	c, _ := LoadChar(char.ToByte())
	return c
}

func (char *Character) ToByte() (bs []byte) {
	attrBs := make([]byte, len(char.Attributes))
	for i, attrib := range char.Attributes {
		attrBs[i] = byte(attrib + 1)
	}
	profBs := make([]byte, len(char.Proficiencies))
	for i, prof := range char.Proficiencies {
		profBs[i] = byte(prof + 1)
	}
	bs = cmp.Merge([][]byte{[]byte(char.Name), char.Attacks, profBs}, attrBs, []byte{byte(char.Race.id), byte(char.Class.id)})
	bs = append(bs, byte(0))
	return
}

var CharLoader = map[byte]func([]byte) *Character{
	0: func(bs []byte) (char *Character) {
		data := cmp.Demerge(bs, []int{ABIL_COUNT, 2})
		char = &Character{
			Race:          Races[int(data[1][0])],
			Class:         Classes[int(data[1][1])],
			Attributes:    make([]int8, 0),
			Proficiencies: make([]int8, 0),
		}
		for _, attrB := range data[0] {
			char.Attributes = append(char.Attributes, int8(attrB)-1)
		}
		for _, profB := range data[4] {
			char.Proficiencies = append(char.Proficiencies, int8(profB)-1)
		}
		char.Name = string(data[2])
		char.Attacks = data[3]
		return
	},
}

func LoadChar(bs []byte) (*Character, error) {
	if len(bs) <= 1 {
		return nil, fmt.Errorf(NO_CHARACTER_DATA)
	}
	v := bs[len(bs)-1]
	bs = bs[:len(bs)-1]
	fnc, ok := CharLoader[v]
	if ok {
		return fnc(bs), nil
	}
	return nil, fmt.Errorf(ERROR_WRONG_CHAR_VERSION, v)
}
