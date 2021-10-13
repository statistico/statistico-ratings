package calculate

import (
	"github.com/statistico/statistico-proto/go"
	"math"
)

func PointsValue(attack, defence float64, k, goals float64) float64 {
	if goals == 0 {
		val := math.Abs(defence - attack)

		if val == 0 {
			return -k
		}

		kg := (defence / attack) * (k * 1.5)

		return -float64(int(kg*100)) / 100
	}

	kg := (defence / attack) * k * goals

	return float64(int(kg*100)) / 100
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
