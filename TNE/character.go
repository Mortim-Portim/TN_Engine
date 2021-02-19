package TNE

import (
	"fmt"
	"strings"
)

const ERROR_WRONG_CHAR_VERSION = "Wrong character version: %v"
var Classes []*Class = []*Class{
	{"Fighter", 	0, []string{"Berserker", "Defender"}},
	{"Knight", 		1, []string{"Heavy Armor", "Combat Rager", "Guard", "Paladin"}},
	{"Wizard", 		2, []string{"Test"}},
	{"Bard", 		3, []string{"Test"}},
	{"Cleric", 		4, []string{"Test"}},
	{"Druid", 		5, []string{"Test"}},
	{"Ranger", 		6, []string{"Test"}},
	{"Rogue", 		7, []string{"Test"}},
	{"Sorcerer", 	8, []string{"Test"}},
}
var Races []*Race = []*Race{
	{"Elv", 		0, []int8{0, 2, 2, 0}, 	[]int8{SCORE_NATURE, SCORE_ACROBATICS, SCORE_PERCEPTION, SCORE_INSIGHT, SCORE_INTELLIGENCE},1, []string{"Wood Elv", "High Elv"}},
	{"Human", 		1, []int8{1, 1, 1, 1}, 	[]int8{}, 																					4, []string{"Urban Human", "Country-Side Human", "Mountain Tribe"}},
	{"Half-Elv", 	2, []int8{0, 2, 0, 2}, 	[]int8{SCORE_ACROBATICS, SCORE_PERCEPTION, SCORE_INSIGHT, SCORE_DUNGEONEERING}, 			1, []string{"Dark Elv", "City Elv"}},
	{"Ork", 		3, []int8{3, 0, 0, 0}, 	[]int8{SCORE_STRENGTH, SCORE_ENDURANCE}, 													1, []string{"Mountain Ork", "Cave Ork"}},
	{"Goblin", 		4, []int8{-1, 2, 0, 2}, []int8{SCORE_STEALTH, SCORE_THIEVERY, SCORE_ACROBATICS, SCORE_DECEPTION, SCORE_PERCEPTION}, 1, []string{"Ravin Goblin", "Sever Goblin"}},
	{"Dwarf", 		5, []int8{2, 0, 2, 0}, 	[]int8{SCORE_STRENGTH, SCORE_CRAFTSMANSHIP, SCORE_DUNGEONEERING}, 							1, []string{"Hill Dwarf", "Mountain Dwarf"}},
	{"Halfling", 	6, []int8{-1, 0, 1, 3}, []int8{SCORE_PERSUASION, SCORE_DECEPTION, SCORE_DUNGEONEERING}, 							1, []string{"Rock Halfling", "Forest Halfling"}},
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
	Name			string
	id				int
	Attributes 		[]int8
	Proficiencies 	[]int8
	Extraprof 		int
	Subraces   		[]string //will change later, placeholder
}

const MAX_CHARACTER_NAME_LENGTH = 20
type Character struct {
	Name          string
	Class         *Class
	Race          *Race
	Attributes    []int8
	Proficiencies []int8
}
func GetCharacter(name string, raceId, classId int) (char *Character) {
	char = &Character{
		Name:			name,
		Race:			Races[raceId],
		Class:			Classes[classId],
		Attributes:		make([]int8, len(Races[raceId].Attributes)),
		Proficiencies:	make([]int8, len(Races[raceId].Proficiencies)),
	}
	copy(char.Attributes, char.Race.Attributes)
	copy(char.Proficiencies, char.Race.Proficiencies)
	return
}
func (char *Character) Copy() *Character {
	if char == nil {return nil}
	c, _ := LoadChar(char.ToByte())
	return c
}
const CHARACTER_BYTES_LENGTH = ABIL_COUNT+SCORE_COUNT+MAX_CHARACTER_NAME_LENGTH+2+1
func (char *Character) ToByte() (bs []byte) {
	bs = make([]byte, CHARACTER_BYTES_LENGTH)
	for i, attrib := range char.Attributes {
		bs[i] = byte(attrib+1)
	}
	for i, prof := range char.Proficiencies {
		bs[i+ABIL_COUNT] = byte(prof+1)
	}
	idx := ABIL_COUNT+SCORE_COUNT
	name := char.Name
	for len(name) < MAX_CHARACTER_NAME_LENGTH {
		name += " "
	}
	copy(bs[idx:idx+MAX_CHARACTER_NAME_LENGTH], []byte(name))
	idx += MAX_CHARACTER_NAME_LENGTH
	bs[idx] = 	byte(char.Race.id)
	bs[idx+1] = byte(char.Class.id)
	bs[idx+2] = 0
	return
}

var CharLoader = map[byte]func([]byte)*Character{
	0:func(bs []byte)(char *Character){
		char = &Character{
			Race:  Races[int(bs[len(bs)-2])],
			Class: Classes[int(bs[len(bs)-1])],
			Attributes:		make([]int8, 0),
			Proficiencies:	make([]int8, 0),
		}
		idx := 0
		for idx < ABIL_COUNT {
			v := int8(bs[idx])-1
			if v >= 0 {
				char.Attributes = append(char.Attributes, v)
			}
			idx ++
		}
		for idx < ABIL_COUNT+SCORE_COUNT {
			v := int8(bs[idx])-1
			if v >= 0 {
				char.Proficiencies = append(char.Proficiencies, v)
			}
			idx ++
		}
		char.Name = strings.ReplaceAll(string(bs[idx:idx+MAX_CHARACTER_NAME_LENGTH]), " ", "")
		return
	},
}

func LoadChar(bs []byte) (*Character, error) {
	v := bs[len(bs)-1]; bs = bs[:len(bs)-1]
	fnc, ok := CharLoader[v]
	if ok {
		return fnc(bs), nil
	}
	return nil, fmt.Errorf(ERROR_WRONG_CHAR_VERSION, v)
}
