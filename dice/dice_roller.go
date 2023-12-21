package dice

import (
	"math/rand"
	"time"
)

func DiceRoll(numberOfDice int, numberOfSides int) int {
	rand.Seed(time.Now().UnixNano())
	total := 0
	for i := 0; i < numberOfDice; i++ {
		total += rand.Intn(numberOfSides) + 1
	}
	return total
}
