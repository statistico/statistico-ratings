package team

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-ratings/internal/app/fixture"
	"time"
)

type RatingHandler struct {
	fetcher        fixture.Fetcher
	processor      RatingProcessor
	logger         *logrus.Logger
}

func (r *RatingHandler) ByCompetition(ctx context.Context, competitionID uint64, numSeasons int8) error {
	fixtures, err := r.fetcher.ByCompetition(ctx, competitionID, numSeasons)

	if err != nil {
		return err
	}

	return r.handleFixtures(ctx, fixtures)
}

func (r *RatingHandler) ByDate(ctx context.Context, time time.Time) error {
	fixtures, err := r.fetcher.ByDate(ctx, time)

	if err != nil {
		return err
	}

	return r.handleFixtures(ctx, fixtures)
}

func (r *RatingHandler) handleFixtures(ctx context.Context, f []*statistico.Fixture) error {
	for _, fix := range f {
		err := r.processor.ByFixture(ctx, fix)

		if err != nil {
			return err
		}
	}

	return nil
}
