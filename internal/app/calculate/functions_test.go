package calculate_test

import (
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-ratings/internal/app/calculate"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPointsValue(t *testing.T) {
	t.Run("returns float values for rating calculation", func(t *testing.T) {
		t.Helper()

		s := []struct {
			Attack  float64
			Defence float64
			G       float64
			KFactor float64
			Points   float64
		}{
			{
				80,
				30,
				2,
				5,
				26.66,
			},
			{
				120,
				150,
				4,
				5,
				16,
			},
			{
				220,
				200,
				2,
				5,
				11,
			},
			{
				5,
				40,
				5,
				35,
				21.87,
			},
			{
				220,
				200,
				0,
				5,
				-8.25,
			},
			{
				220,
				50,
				0,
				5,
				-33,
			},
			{
				100,
				200,
				0,
				5,
				-3.75,
			},
		}

		for _, st := range s {
			a := calculate.PointsValue(st.Attack, st.Defence, st.KFactor, st.G)

			assert.Equal(t, st.Points, a)
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
