package calculate

import "math"

// PointsValue calculates the value for the attack vs defence scenario using
// an ELO based calculation. k argument relates to the k factor used with elo
// calculation see https://en.wikipedia.org/wiki/Elo_rating_system#The_K-factor_used_by_the_USCF
func PointsValue(attack, defence int64, k int8, goals float64) float64 {
	ge := GoalExpectancy(attack, defence)
	val := (float64(k)*adjustGoals(goals)) * (goals - ge)
	return float64(int(val * 100)) / 100
}

// GoalExpectancy calculates the expected goal probability based on attack and defence rating values.
func GoalExpectancy(attack, defence int64) float64 {
	diff := attack - defence
	d := float64(-diff) / 400
	pow := 1 / (math.Pow(10, d) + 1)
	return float64(int(pow * 100)) / 100
}

func adjustGoals(g float64) float64 {
	if g == 0.0 {
		g = 1.0
	}

	return g
}
