package GE

import (
	"math/rand"
)
type RandomPool []int
func (p *RandomPool) Add(i int) {
	*p = (RandomPool(append(([]int)(*p), i)))
}
func (p *RandomPool) AddL(iL ...int) {
	for _,i := range(iL) {
		p.Add(i)
	}
}
func (p *RandomPool) Get() (i int) {
	i = ([]int)(*p)[0]
	*p = RandomPool(([]int)(*p)[1:])
	return
}
func GenerateDs(sides, count int) RandomPool {
	random_d20s := make([]int, count)
	for i,_ := range(random_d20s) {
		random_d20s[i] = RollDice(1, sides)
	}
	return RandomPool(random_d20s)
}

//Rolls a n dices of m sides and returns the result
func RollDice(dices, sides int) int {
	res := 0
	for i := 0; i < dices; i++ {
		res += rand.Intn(sides)+1
	}
	return res
}