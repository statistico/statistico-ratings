package calculate_test

import (
	statistico "github.com/statistico/statistico-proto/go"
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
			K       uint8
			G       float64
			AttackValue   float64
			DefenceValue   float64
		}{
			{
				80,
				30,
				50,
				2,
				3.33,
				1.25,
			},
			{
				120,
				150,
				50,
				4,
				1.33,
				1.66,
			},
			{
				220,
				200,
				50,
				2,
				0.50,
				0.45,
			},
			{
				50,
				40,
				50,
				5,
				6.25,
				5,
			},
			{
				220,
				200,
				50,
				0,
				-2.50,
				-2.50,
			},
			{
				220,
				50,
				50,
				0,
				-0.29,
				-0.29,
			},
			{
				100,
				200,
				50,
				0,
				-0.50,
				-0.50,
			},
			{
				170,
				100,
				50,
				1,
				0.50,
				0.29,
			},
		}

		for _, st := range s {
			a, d := calculate.PointsValue(st.Attack, st.Defence, st.K, st.G)

			assert.Equal(t, st.AttackValue, a)
			assert.Equal(t, st.DefenceValue, d)
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
