package TNE

import (
	"math/rand"
	//"sync"
	"time"
)

//Rolls n dices of m sides and returns the sum of the results
func RollDice(dices, sides int) int {
	res := 0
	for i := 0; i < dices; i++ {
		res += rand.Intn(sides)+1
	}
	return res
}

//Generates a random number including l and u
func RandomInt(l, u int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(u-l+1) + l
}

func RandomFloat(l, u float64) float64 {
	return rand.Float64()*(u-l+1)+l
}