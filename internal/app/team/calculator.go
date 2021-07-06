package team

import (
	"context"
	"github.com/jonboulle/clockwork"
	"github.com/statistico/statistico-football-data-go-grpc-client"
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-ratings/internal/app/calculate"
)

type RatingCalculator interface {
	// ForFixture receives a statistico.Fixture struct and home and away Rating struct, calculates points for
	// each team for the fixture and returns the newly calculated Rating struct for each team.
	ForFixture(ctx context.Context, f *statistico.Fixture, home, away *Rating) (*Rating, *Rating, error)
}

type ratingCalculator struct {
	event statisticodata.EventClient
	clock clockwork.Clock
}

func (r *ratingCalculator) ForFixture(ctx context.Context, f *statistico.Fixture, home, away *Rating) (*Rating, *Rating, error) {
	events, err := r.event.FixtureEvents(ctx, uint64(f.Id))

	if err != nil {
		return nil, nil, err
	}

	hg, ag := calculate.AdjustedGoals(f.HomeTeam.Id, f.AwayTeam.Id, events.Goals, events.Cards)

	hv := calculate.PointsValue(home.Attack.Total, away.Defence.Total, 20, hg)
	av := calculate.PointsValue(away.Attack.Total, home.Defence.Total, 20, ag)

	newHome := r.applyRating(home, uint64(f.Id), f.Season.Id, hv, av)
	newAway := r.applyRating(away, uint64(f.Id), f.Season.Id, av, hv)

	return newHome, newAway, nil
}

func (r *ratingCalculator) applyRating(rt *Rating, fixtureID, seasonID uint64, attack, defence float64) *Rating {
	return &Rating{
		TeamID:    rt.TeamID,
		FixtureID: fixtureID,
		SeasonID:  seasonID,
		Attack:    Points{
			Total:      rt.Attack.Total + attack,
			Difference: attack,
		},
		Defence:   Points{
			Total:      rt.Defence.Total - defence,
			Difference: -defence,
		},
		Timestamp: r.clock.Now(),
	}
}

func NewRatingCalculator(e statisticodata.EventClient, c clockwork.Clock) RatingCalculator {
	return &ratingCalculator{
		event: e,
		clock: c,
	}
}
