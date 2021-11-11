package team

import (
	"context"
	"errors"
	"github.com/jonboulle/clockwork"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-ratings/internal/app/fixture"
	"time"
)

type RatingHandler struct {
	fetcher   fixture.Fetcher
	processor RatingProcessor
	clock     clockwork.Clock
	logger    *logrus.Logger
}

func (r *RatingHandler) ByCompetition(ctx context.Context, competitionID, seasonID uint64) {
	fixtures, err := r.fetcher.ByCompetition(ctx, competitionID, seasonID)

	if err != nil {
		r.logger.Errorf("error fetching fixtures in team rating handler: %s", err.Error())
		return
	}

	for _, fix := range fixtures {
		err := r.processor.ByFixture(ctx, fix)

		if err != nil {
			r.logger.Errorf("error processing fixtures in team rating handler: %s", err.Error())
			return
		}
	}

	return
}

// Today processes team ratings for the days fixtures. The hour argument is used to determine what hour of the
// day the ratings are to be processed i.e. hour 20 means process all fixture ratings for fixtures before 8pm
func (r *RatingHandler) Today(ctx context.Context, hour int) error {
	now := r.clock.Now()

	if hour > now.Hour() {
		return errors.New("hour provided is greater than the current time hour")
	}

	year, month, day := now.Date()

	start := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	end := time.Date(year, month, day, hour, 0, 0, 0, time.UTC)

	fixtures, err := r.fetcher.ByDate(ctx, start, end)

	if err != nil {
		r.logger.Errorf("error fetching fixtures in team rating handler: %s", err.Error())
		return err
	}

	for _, fix := range fixtures {
		err := r.processor.ByFixture(ctx, fix)

		if err != nil {
			r.logger.Errorf("error processing fixtures in team rating handler: %s", err.Error())
			continue
		}
	}

	return nil
}

func NewHandler(f fixture.Fetcher, p RatingProcessor, c clockwork.Clock, l *logrus.Logger) RatingHandler {
	return RatingHandler{
		fetcher:   f,
		processor: p,
		clock:     c,
		logger:    l,
	}
}
