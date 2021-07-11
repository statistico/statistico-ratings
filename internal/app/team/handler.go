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

func (r *RatingHandler) ByCompetition(ctx context.Context, competitionID uint64, numSeasons int8) {
	fixtures, err := r.fetcher.ByCompetition(ctx, competitionID, numSeasons)

	if err != nil {
		r.logger.Error("error fetching fixtures in team rating handler")
		return
	}

	for _, fix := range fixtures {
		err := r.processor.ByFixture(ctx, fix)

		if err != nil {
			r.logger.Error("error processing fixtures in team rating handler")
			return
		}
	}

	return
}

func (r *RatingHandler) ByDate(ctx context.Context, time time.Time) {
	fixtures, err := r.fetcher.ByDate(ctx, time)

	if err != nil {
		r.logger.Error("error fetching fixtures in team rating handler")
		return
	}

	for _, fix := range fixtures {
		err := r.processor.ByFixture(ctx, fix)

		if err != nil {
			r.logger.Error("error processing fixtures in team rating handler")
			continue
		}
	}
}

func (r *RatingHandler) handleFixtures(ctx context.Context, f []*statistico.Fixture) {
	for _, fix := range f {
		err := r.processor.ByFixture(ctx, fix)

		if err != nil {
			r.logger.Error("error processing fixtures in team rating handler")
			return
		}
	}

	return
}

func NewHandler(f fixture.Fetcher, p RatingProcessor, l *logrus.Logger) RatingHandler {
	return RatingHandler{
		fetcher:   f,
		processor: p,
		logger:    l,
	}
}
