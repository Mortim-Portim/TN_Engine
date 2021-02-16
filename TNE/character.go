package TNE

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

type Class struct {
	Name     string
	id       int
	Subclass []string //will change later, placeholder
}

const ABIL_STRENGTH = 0
const ABIL_DEXTERITY = 1
const ABIL_INTELLIGENCE = 2
const ABIL_CHARISMA = 3

const SCORE_STRENGTH = 0
const SCORE_DEXTERITY = 1
const SCORE_INTELLIGENCE = 2
const SCORE_CHARISMA = 3
const SCORE_ENDURANCE = 4
const SCORE_PERSUASION = 5
const SCORE_DECEPTION = 6
const SCORE_PERFORMANCE = 7
const SCORE_INSIGHT = 8
const SCORE_THIEVERY = 9
const SCORE_STEALTH = 10
const SCORE_ACROBATICS = 11
const SCORE_NATURE = 12
const SCORE_ARCANA = 13
const SCORE_PERCEPTION = 14
const SCORE_CRAFTSMANSHIP = 15
const SCORE_DUNGEONEERING = 16

var Races []*Race = []*Race{
	{"Elv", 0, []int8{0, 2, 2, 0}, []int{SCORE_NATURE, SCORE_ACROBATICS, SCORE_PERCEPTION, SCORE_INSIGHT, SCORE_INTELLIGENCE}, 1, []string{"Wood Elv", "High Elv"}},
	{"Human", 1, []int8{1, 1, 1, 1}, []int{}, 4, []string{"Urban Human", "Country-Side Human", "Mountain Tribe"}},
	{"Half-Elv", 2, []int8{0, 2, 0, 2}, []int{SCORE_ACROBATICS, SCORE_PERCEPTION, SCORE_INSIGHT, SCORE_DUNGEONEERING}, 1, []string{"Dark Elv", "City Elv"}},
	{"Ork", 3, []int8{3, 0, 0, 0}, []int{SCORE_STRENGTH, SCORE_ENDURANCE}, 1, []string{"Mountain Ork", "Cave Ork"}},
	{"Goblin", 4, []int8{-1, 2, 0, 2}, []int{SCORE_STEALTH, SCORE_THIEVERY, SCORE_ACROBATICS, SCORE_DECEPTION, SCORE_PERCEPTION}, 1, []string{"Ravin Goblin", "Sever Goblin"}},
	{"Dwarf", 5, []int8{+2, 0, +2, 0}, []int{SCORE_STRENGTH, SCORE_CRAFTSMANSHIP, SCORE_DUNGEONEERING}, 1, []string{"Hill Dwarf", "Mountain Dwarf"}},
	{"Halfling", 6, []int8{-1, 0, +1, +3}, []int{SCORE_PERSUASION, SCORE_DECEPTION, SCORE_DUNGEONEERING}, 1, []string{"Rock Halfling", "Forest Halfling"}},
}

type Race struct {
	Name       string
	id         int
	Attributes []int8
	Profencies []int
	Extraprof  int
	Subraces   []string //will change later, placeholder
}

type Character struct {
	Name          string
	Class         *Class
	Race          *Race
	Attributes    []int8
	Proficiencies []int8
}
func (char *Character) Copy() *Character {
	if char == nil {return nil}
	return LoadChar(char.ToByte())
}
const CHARACTER_BYTES_LENGTH = 23
func (char *Character) ToByte() []byte {
	byteattrib := make([]byte, len(char.Attributes))
	for i, attrib := range char.Attributes {
		byteattrib[i] = byte(attrib)
	}

	byteprof := make([]byte, len(char.Proficiencies))
	for i, prof := range char.Proficiencies {
		byteprof[i] = byte(prof)
	}

	bytearray := make([]byte, 0)
	bytearray = append(bytearray, byte(char.Race.id), byte(char.Class.id))
	bytearray = append(bytearray, byteattrib...)
	bytearray = append(bytearray, byteprof...)
	return bytearray
}

func LoadChar(bytes []byte) *Character {
	race := Races[bytes[0]]
	class := Classes[bytes[1]]

	attrib := make([]int8, 4)
	for i := 0; i < len(attrib); i++ {
		attrib[i] = int8(bytes[i+2])
	}

	proficiencies := make([]int8, 17)
	for i := 0; i < len(proficiencies); i++ {
		proficiencies[i] = int8(bytes[i+6])
	}

	char := &Character{Race: race, Class: class, Attributes: attrib, Proficiencies: proficiencies}

	return char
}
