package calculate_test

import (
	"github.com/statistico/statistico-ratings/internal/app/calculate"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEloValue(t *testing.T) {
	t.Run("returns float value for elo calculation", func(t *testing.T) {
		t.Helper()

		s := []struct{
			Attack  int64
			Defence int64
			K       int8
			G       float64
			Value   float64
		} {
			{
				1745,
				1300,
				30,
				1.50,
				26.10,
			},
			{
				1745,
				1300,
				30,
				0.0,
				-27.60,
			},
			{
				1645,
				1700,
				10,
				3.0,
				77.40,
			},
			{
				1645,
				1640,
				10,
				0,
				-5.0,
			},
		}

		for _, st := range s {
			e := calculate.EloValue(st.Attack, st.Defence, st.K, st.G)

			assert.Equal(t, st.Value, e)
		}
	})
}
func TestGoalExpectancy(t *testing.T) {
	t.Run("returns calculated goal expectancy", func(t *testing.T) {
		t.Helper()

		s := []struct{
			Attack  int64
			Defence int64
			GoalExpectancy float64
		} {
			{
				1745,
				1300,
				0.92,
			},
			{
				1610,
				1600,
				0.51,
			},
			{
				1610,
				1800,
				0.25,
			},
		}

		for _, st := range s {
			e := calculate.GoalExpectancy(st.Attack, st.Defence)

			assert.Equal(t, st.GoalExpectancy, e)
		}
	})
}
