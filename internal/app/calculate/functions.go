package calculate

import (
	"github.com/statistico/statistico-proto/go"
	"math"
)

// PointsValue calculates the value for the attack vs defence scenario using
// an ELO based calculation. k argument relates to the k factor used with elo
// calculation see https://en.wikipedia.org/wiki/Elo_rating_system#The_K-factor_used_by_the_USCF
func PointsValue(attack, defence float64, k uint8, goals float64) (float64, float64) {
	// Adjust goals for true reflection
	if goals == 0 {
		d := float64(k) / math.Abs(attack - defence)

		return -float64(int(d*100)) / 100, -float64(int(d*100)) / 100
	}

	kg := float64(k) * goals

	a := kg / defence
	d := kg / attack

	return float64(int(a*100)) / 100, float64(int(d*100)) / 100
}

// GoalExpectancy calculates the expected goal probability based on attack and defence rating values.
func GoalExpectancy(attack, defence float64) float64 {
	diff := -defence + attack
	//d := float64(diff) / 400
	//pow := 1 / (math.Pow(10, d) + 1)
	//val := float64(int(pow*100)) / 100
	//
	//if defence > attack {
	//	return 1 - val
	//}

	return 1 - (1 / math.Sqrt(diff))
}

// AdjustedGoals calculates the value of the goals scored for each team. A goal value can be increased or decreased
// based on factors such as red cards, minute of goal and current score difference. The two float64 values returned
// are the home goals as the first value and away goals as the second value.
func AdjustedGoals(homeID, awayID uint64, goals []*statistico.GoalEvent, cards []*statistico.CardEvent) (float64, float64) {
	var home int8
	var away int8
	var homeAdj float64
	var awayAdj float64

	for _, goal := range goals {
		homeRed := hasBeenRedCard(cards, homeID, goal.Minute)
		awayRed := hasBeenRedCard(cards, awayID, goal.Minute)

		if goal.TeamId == homeID {
			home++
			diff := float64(home - away)
			homeAdj += calculateGoalValue(diff, goal.Minute, 70, homeRed, awayRed)
		}

		if goal.TeamId == awayID {
			away++
			diff := float64(away - home)
			awayAdj += calculateGoalValue(diff, goal.Minute, 70, awayRed, homeRed)
		}
	}

	return float64(int(homeAdj*100)) / 100, float64(int(awayAdj*100)) / 100
}

func calculateGoalValue(diff float64, min, clock uint32, teamRed, oppRed bool) float64 {
	g := 1.0

	if min > clock {
		if diff >= 2 {
			g = g / diff
		}
	}

	if oppRed {
		g = g * 0.75
	}

	if teamRed {
		g = g / 0.75
	}

	return g
}

func hasBeenRedCard(cards []*statistico.CardEvent, teamID uint64, min uint32) bool {
	for _, card := range cards {
		if card.Type == "redcard" && card.TeamId == teamID && card.Minute < min {
			return true
		}
	}

	return false
}
