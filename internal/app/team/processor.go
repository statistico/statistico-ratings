package team

import (
	"context"
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-ratings/internal/app"
)

type RatingProcessor interface {
	ByFixture(ctx context.Context, f *statistico.Fixture) error
}

type ratingProcessor struct {
	reader RatingReader
	writer RatingWriter
	calculator RatingCalculator
}

func (r *ratingProcessor) ByFixture(ctx context.Context, f *statistico.Fixture) error {
	home, err := r.fetchRating(f.HomeTeam.Id)

	if err != nil {
		return err
	}

	away, err := r.fetchRating(f.AwayTeam.Id)

	if err != nil {
		return err
	}

	newHome, newAway, err := r.calculator.ForFixture(ctx, f, home, away)

	if err != nil {
		return err
	}

	err = r.writer.Insert(newHome)

	if err != nil {
		return err
	}

	err = r.writer.Insert(newAway)

	if err != nil {
		return err
	}

	return nil
}

func (r *ratingProcessor) fetchRating(teamID uint64) (*Rating, error) {
	rating, err := r.reader.Latest(teamID)

	switch err.(type) {
	case *app.NotFoundError:
		return &Rating{
			TeamID:    teamID,
			Attack:    Points{
				Total:      1500,
				Difference: 0,
			},
			Defence:   Points{
				Total:      1500,
				Difference: 0,
			},
		}, nil
	case nil:
		return rating, nil
	default:
		return nil, err
	}
}

func NewRatingProcessor(r RatingReader, w RatingWriter, c RatingCalculator) RatingProcessor {
	return &ratingProcessor{
		reader:     r,
		writer:     w,
		calculator: c,
	}
}
