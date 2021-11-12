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
			Id: 1,
		},
		AwayTeam: &statistico.Team{
			Id: 8,
		},
		DateTime: &statistico.Date{Utc: 1630343736},
	}

	home := team.Rating{
		TeamID:    1,
		FixtureID: 25,
		SeasonID:  17462,
		Attack: team.Points{
			Total:      1225.67,
			Difference: 8.06,
		},
		Defence: team.Points{
			Total:      1228.47,
			Difference: -6.85,
		},
		FixtureDate: time.Unix(1630343736, 0),
		Timestamp:   time.Now(),
	}

	away := team.Rating{
		TeamID:    8,
		FixtureID: 25,
		SeasonID:  17462,
		Attack: team.Points{
			Total:      1518.33,
			Difference: -5.67,
		},
		Defence: team.Points{
			Total:      790.72,
			Difference: 17.07,
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
					TeamId: 1,
					Minute: 4,
				},
				{
					TeamId: 8,
					Minute: 42,
				},
				{
					TeamId: 1,
					Minute: 67,
				},
				{
					TeamId: 1,
					Minute: 74,
				},
				{
					TeamId: 8,
					Minute: 83,
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

		a.Equal(uint64(1), newHome.TeamID)
		a.Equal(uint64(26), newHome.FixtureID)
		a.Equal(uint64(17462), newHome.SeasonID)
		a.Equal(1245.04, newHome.Attack.Total)
		a.Equal(19.37, newHome.Attack.Difference)
		a.Equal(1240.82, newHome.Defence.Total)
		a.Equal(12.35, newHome.Defence.Difference)
		a.Equal(time.Unix(1630343736, 0), newHome.FixtureDate)
		a.Equal(time.Date(1984, time.April, 4, 0, 0, 0, 0, time.UTC), newHome.Timestamp)

		a.Equal(uint64(8), newAway.TeamID)
		a.Equal(uint64(26), newAway.FixtureID)
		a.Equal(uint64(17462), newAway.SeasonID)
		a.Equal(1530.6799999999998, newAway.Attack.Total)
		a.Equal(12.35, newAway.Attack.Difference)
		a.Equal(810.09, newAway.Defence.Total)
		a.Equal(19.37, newAway.Defence.Difference)
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
