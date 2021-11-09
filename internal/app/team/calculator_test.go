package team_test

import (
	"context"
	"errors"
	"github.com/jonboulle/clockwork"
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-ratings/internal/app/bootstrap"
	"github.com/statistico/statistico-ratings/internal/app/team"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestRatingCalculator_ForFixture(t *testing.T) {
	fixture := statistico.Fixture{
		Id: 26,
		Competition: &statistico.Competition{
			Id: 8,
		},
		Season: &statistico.Season{
			Id: 17462,
		},
		HomeTeam: &statistico.Team{
			Id: 10,
		},
		AwayTeam: &statistico.Team{
			Id: 12,
		},
		DateTime: &statistico.Date{Utc: 1630343736},
	}

	home := team.Rating{
		TeamID:    10,
		FixtureID: 25,
		SeasonID:  17462,
		Attack: team.Points{
			Total:      167.10,
			Difference: 5,
		},
		Defence: team.Points{
			Total:      152.45,
			Difference: -4,
		},
		FixtureDate: time.Unix(1630343736, 0),
		Timestamp:   time.Now(),
	}

	away := team.Rating{
		TeamID:    12,
		FixtureID: 25,
		SeasonID:  17462,
		Attack: team.Points{
			Total:      170,
			Difference: 15,
		},
		Defence: team.Points{
			Total:      42.45,
			Difference: -14,
		},
		FixtureDate: time.Unix(1630343736, 0),
		Timestamp:   time.Now(),
	}

	t.Run("calculates new team ratings for a fixture", func(t *testing.T) {
		t.Helper()

		events := new(MockEventClient)
		config := bootstrap.BuildConfig()
		clock := clockwork.NewFakeClock()

		calculator := team.NewRatingCalculator(events, config.KFactorMapping, clock)

		ctx := context.Background()

		res := statistico.FixtureEventsResponse{
			Goals: []*statistico.GoalEvent{
				{
					TeamId: 10,
					Minute: 47,
				},
			},
		}

		events.On("FixtureEvents", ctx, uint64(26)).Return(&res, nil)

		newHome, newAway, err := calculator.ForFixture(ctx, &fixture, &home, &away)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		events.AssertExpectations(t)

		a := assert.New(t)

		a.Equal(uint64(10), newHome.TeamID)
		a.Equal(uint64(26), newHome.FixtureID)
		a.Equal(uint64(17462), newHome.SeasonID)
		a.Equal(168.37, newHome.Attack.Total)
		a.Equal(1.27, newHome.Attack.Difference)
		a.Equal(145.73, newHome.Defence.Total)
		a.Equal(-6.72, newHome.Defence.Difference)
		a.Equal(time.Unix(1630343736, 0), newHome.FixtureDate)
		a.Equal(time.Date(1984, time.April, 4, 0, 0, 0, 0, time.UTC), newHome.Timestamp)

		a.Equal(uint64(12), newAway.TeamID)
		a.Equal(uint64(26), newAway.FixtureID)
		a.Equal(uint64(17462), newAway.SeasonID)
		a.Equal(163.28, newAway.Attack.Total)
		a.Equal(-6.72, newAway.Attack.Difference)
		a.Equal(43.720000000000006, newAway.Defence.Total)
		a.Equal(1.27, newAway.Defence.Difference)
		a.Equal(time.Unix(1630343736, 0), newAway.FixtureDate)
		a.Equal(time.Date(1984, time.April, 4, 0, 0, 0, 0, time.UTC), newHome.Timestamp)
	})

	t.Run("returns an error if returned by event client", func(t *testing.T) {
		t.Helper()

		events := new(MockEventClient)
		config := bootstrap.BuildConfig()
		clock := clockwork.NewFakeClock()

		calculator := team.NewRatingCalculator(events, config.KFactorMapping, clock)

		ctx := context.Background()

		e := errors.New("error in event client")

		events.On("FixtureEvents", ctx, uint64(26)).Return(&statistico.FixtureEventsResponse{}, e)

		newHome, newAway, err := calculator.ForFixture(ctx, &fixture, &home, &away)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		events.AssertExpectations(t)

		assert.Nil(t, newHome)
		assert.Nil(t, newAway)
	})
}

type MockEventClient struct {
	mock.Mock
}

func (m *MockEventClient) FixtureEvents(ctx context.Context, fixtureID uint64) (*statistico.FixtureEventsResponse, error) {
	args := m.Called(ctx, fixtureID)
	return args.Get(0).(*statistico.FixtureEventsResponse), args.Error(1)
}
