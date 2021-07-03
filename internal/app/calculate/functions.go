package calculate

import "math"

func EloValue(attack, defence int64, k int8, g float64) float64 {
	ge := GoalExpectancy(attack, defence)
	val := (float64(k)*adjustGoals(g)) * (g - ge)
	return float64(int(val * 100)) / 100
}

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
