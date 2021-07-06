package calculate_test

import (
	statistico "github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-ratings/internal/app/calculate"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPointsValue(t *testing.T) {
	t.Run("returns float value for elo calculation", func(t *testing.T) {
		t.Helper()

		s := []struct {
			Attack  float64
			Defence float64
			K       int8
			G       float64
			Value   float64
		}{
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
			e := calculate.PointsValue(st.Attack, st.Defence, st.K, st.G)

			assert.Equal(t, st.Value, e)
		}
	})
}
func TestGoalExpectancy(t *testing.T) {
	t.Run("returns calculated goal expectancy", func(t *testing.T) {
		t.Helper()

		s := []struct {
			Attack         float64
			Defence        float64
			GoalExpectancy float64
		}{
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

func TestAdjustedGoals(t *testing.T) {
	t.Run("returns values for home and away adjusted goals", func(t *testing.T) {
		t.Helper()

		s := []struct {
			HomeID    uint64
			AwayID    uint64
			Goals     []*statistico.GoalEvent
			Cards     []*statistico.CardEvent
			HomeGoals float64
			AwayGoals float64
		}{
			{
				1,
				2,
				[]*statistico.GoalEvent{},
				[]*statistico.CardEvent{},
				0.0,
				0.0,
			},
			{
				1,
				2,
				[]*statistico.GoalEvent{
					{
						TeamId: 1,
						Minute: 25,
					},
					{
						TeamId: 1,
						Minute: 42,
					},
				},
				[]*statistico.CardEvent{},
				2.0,
				0.0,
			},
			{
				1,
				2,
				[]*statistico.GoalEvent{
					{
						TeamId: 1,
						Minute: 25,
					},
					{
						TeamId: 1,
						Minute: 42,
					},
				},
				[]*statistico.CardEvent{
					{
						TeamId:   2,
						Type:     "redcard",
						PlayerId: 0,
						Minute:   2,
					},
				},
				1.5,
				0.0,
			},
			{
				1,
				2,
				[]*statistico.GoalEvent{
					{
						TeamId: 1,
						Minute: 25,
					},
					{
						TeamId: 1,
						Minute: 42,
					},
					{
						TeamId: 2,
						Minute: 85,
					},
					{
						TeamId: 2,
						Minute: 89,
					},
				},
				[]*statistico.CardEvent{
					{
						TeamId:   2,
						Type:     "redcard",
						PlayerId: 0,
						Minute:   2,
					},
				},
				1.5,
				2.66,
			},
			{
				1,
				2,
				[]*statistico.GoalEvent{
					{
						TeamId: 1,
						Minute: 25,
					},
					{
						TeamId: 1,
						Minute: 42,
					},
					{
						TeamId: 2,
						Minute: 85,
					},
					{
						TeamId: 2,
						Minute: 89,
					},
					{
						TeamId: 1,
						Minute: 95,
					},
				},
				[]*statistico.CardEvent{
					{
						TeamId:   2,
						Type:     "yellowcard",
						PlayerId: 0,
						Minute:   2,
					},
				},
				3.0,
				2.0,
			},
			{
				1,
				2,
				[]*statistico.GoalEvent{
					{
						TeamId: 1,
						Minute: 25,
					},
					{
						TeamId: 1,
						Minute: 42,
					},
					{
						TeamId: 1,
						Minute: 85,
					},
					{
						TeamId: 1,
						Minute: 89,
					},
					{
						TeamId: 1,
						Minute: 95,
					},
				},
				[]*statistico.CardEvent{
					{
						TeamId:   2,
						Type:     "redcard",
						PlayerId: 0,
						Minute:   2,
					},
				},
				2.08,
				0.0,
			},
		}

		for _, st := range s {
			home, away := calculate.AdjustedGoals(st.HomeID, st.AwayID, st.Goals, st.Cards)

			assert.Equal(t, st.HomeGoals, home)
			assert.Equal(t, st.AwayGoals, away)
		}
	})
}
