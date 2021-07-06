package team

import (
	"context"
	"github.com/statistico/statistico-proto/go"
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
	home, err := r.reader.Latest(f.HomeTeam.Id)

	if err != nil {
		return err
	}

	away, err := r.reader.Latest(f.AwayTeam.Id)

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

func NewRatingProcessor(r RatingReader, w RatingWriter, c RatingCalculator) RatingProcessor {
	return &ratingProcessor{
		reader:     r,
		writer:     w,
		calculator: c,
	}
}
