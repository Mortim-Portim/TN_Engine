package TNE

import "math/rand"

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

//SHOULD contain information about the races stats
type Race struct {
	*Entity
	Ability       []int8
	Proficiencies []int8
}

func (r *Race) Copy() (r2 *Race) {
	r2 = &Race{Entity: r.Entity.Copy()}
	return
}

func (race *Race) GetScore(index int) int {
	return rand.Intn(20) + 1 + int(race.Proficiencies[index])
}

func LoadFromByte(byt []byte) *Race {
	ability := make([]int8, 4)
	for i := 2; i < 6; i++ {
		ability[i] = int8(byt[i])
	}

	proficiencies := make([]int8, 17)
	for i := 6; i < 23; i++ {
		proficiencies[i] = int8(proficiencies[i])
	}

	return nil
}
