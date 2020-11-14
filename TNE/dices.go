package TNE

import (
	"math/rand"
	//"sync"
	//"time"
)

//Rolls n dices of m sides and returns the result
func RollDice(dices, sides int) int {
	res := 0
	for i := 0; i < dices; i++ {
		res += rand.Intn(sides)+1
	}
	return res
}