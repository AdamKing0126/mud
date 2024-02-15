package dice

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func DiceRoll(dice string) int32 {
	rand.Seed(time.Now().UnixNano())

	// Split the dice string into its components
	parts := strings.Split(dice, "d")
	numberOfDice, _ := strconv.Atoi(parts[0])
	bonus := 0

	// Check if the bonus is present
	var numberOfSides int
	if strings.Contains(parts[1], "+") {
		bonusParts := strings.Split(parts[1], "+")
		numberOfSides, _ = strconv.Atoi(bonusParts[0])
		bonus, _ = strconv.Atoi(bonusParts[1])
	} else {
		numberOfSides, _ = strconv.Atoi(parts[1])
	}

	total := 0
	for i := 0; i < numberOfDice; i++ {
		total += rand.Intn(numberOfSides) + 1
	}
	return int32(total + bonus)
}
