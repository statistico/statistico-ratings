package team_test

import (
	"context"
	"errors"
	"github.com/jonboulle/clockwork"
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-ratings/internal/app/team"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestRatingCalculator_ForFixture(t *testing.T) {
	fixture := statistico.Fixture{
		Id:                   26,
		Season:               &statistico.Season{
			Id:                   17462,
		},
		HomeTeam:             &statistico.Team{
			Id:                   10,
		},
		AwayTeam:             &statistico.Team{
			Id:                   12,
		},
	}

	home := team.Rating{
		TeamID:    10,
		FixtureID: 25,
		SeasonID:  17462,
		Attack:    team.Points{
			Total:      1670,
			Difference: 5,
		},
		Defence:   team.Points{
			Total:      1542.45,
			Difference: -4,
		},
		Timestamp: time.Now(),
	}

	away := team.Rating{
		TeamID:    12,
		FixtureID: 25,
		SeasonID:  17462,
		Attack:    team.Points{
			Total:      1870,
			Difference: 15,
		},
		Defence:   team.Points{
			Total:      1642.45,
			Difference: -14,
		},
		Timestamp: time.Now(),
	}

	t.Run("calculates new team ratings for a fixture", func(t *testing.T) {
		t.Helper()

		events := new(MockEventClient)
		clock := clockwork.NewFakeClock()

		calculator := team.NewRatingCalculator(events, clock)

		ctx := context.Background()

		res := statistico.FixtureEventsResponse{
			Goals:                []*statistico.GoalEvent{
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

		a := assert.New(t)

		a.Equal(uint64(10), newHome.TeamID)
		a.Equal(uint64(26), newHome.FixtureID)
		a.Equal(uint64(17462), newHome.SeasonID)
		a.Equal(1679.39, newHome.Attack.Total)
		a.Equal(9.39, newHome.Attack.Difference)
		a.Equal(1559.65, newHome.Defence.Total)
		a.Equal(17.20, newHome.Defence.Difference)
		a.Equal(time.Date(1984, time.April, 4, 0, 0, 0, 0, time.UTC), newHome.Timestamp)

		a.Equal(uint64(12), newAway.TeamID)
		a.Equal(uint64(26), newAway.FixtureID)
		a.Equal(uint64(17462), newAway.SeasonID)
		a.Equal(1852.8, newAway.Attack.Total)
		a.Equal(-17.20, newAway.Attack.Difference)
		a.Equal(1633.06, newAway.Defence.Total)
		a.Equal(-9.39, newAway.Defence.Difference)
		a.Equal(time.Date(1984, time.April, 4, 0, 0, 0, 0, time.UTC), newHome.Timestamp)
	})

	t.Run("returns an error if returned by event client", func(t *testing.T) {
		t.Helper()

		events := new(MockEventClient)
		clock := clockwork.NewFakeClock()

		calculator := team.NewRatingCalculator(events, clock)

		ctx := context.Background()

		e := errors.New("error in event client")

		events.On("FixtureEvents", ctx, uint64(26)).Return(&statistico.FixtureEventsResponse{}, e)

		newHome, newAway, err := calculator.ForFixture(ctx, &fixture, &home, &away)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

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
