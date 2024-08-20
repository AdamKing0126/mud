package utilities

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

func CalculateAbilityModifier(val int32) int32 {
	switch {
	case val < 1:
		return -6
	case val == 1:
		return -5
	case val <= 3:
		return -4
	case val <= 5:
		return -3
	case val <= 7:
		return -2
	case val <= 9:
		return -1
	case val <= 11:
		return 0
	case val <= 13:
		return 1
	case val <= 15:
		return 2
	case val <= 17:
		return 3
	case val <= 19:
		return 4
	case val <= 21:
		return 5
	case val <= 23:
		return 6
	case val <= 25:
		return 7
	case val <= 27:
		return 8
	case val <= 29:
		return 9
	case val <= 31:
		return 10
	case val <= 33:
		return 11
	case val <= 35:
		return 12
	case val <= 37:
		return 13
	case val <= 39:
		return 14
	case val <= 41:
		return 15
	case val <= 43:
		return 16
	default:
		return 17
	}
}
