package team

import (
	"context"
	"fmt"
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-ratings/internal/app"
)

type RatingProcessor interface {
	ByFixture(ctx context.Context, f *statistico.Fixture) error
}

type ratingProcessor struct {
	reader     RatingReader
	writer     RatingWriter
	calculator RatingCalculator
	competitionMapping map[uint64]uint16
}

func (r *ratingProcessor) ByFixture(ctx context.Context, f *statistico.Fixture) error {
	home, err := r.fetchRating(f.HomeTeam.Id, f.Competition.Id)

	if err != nil {
		return err
	}

	away, err := r.fetchRating(f.AwayTeam.Id, f.Competition.Id)

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

func (r *ratingProcessor) fetchRating(teamID, competitionID uint64) (*Rating, error) {
	rating, err := r.reader.Latest(teamID)

	switch err.(type) {
	case *app.NotFoundError:
		score, err := r.parseCompetitionScore(competitionID)

		if err != nil {
			return nil, err
		}

		return &Rating{
			TeamID: teamID,
			Attack: Points{
				Total:      float64(score),
				Difference: 0,
			},
			Defence: Points{
				Total:      float64(score),
				Difference: 0,
			},
		}, nil
	case nil:
		return rating, nil
	default:
		return nil, err
	}
}

func (r *ratingProcessor) parseCompetitionScore(competitionID uint64) (uint16, error) {
	for competition, score := range r.competitionMapping {
		if competitionID == competition {
			return score, nil
		}
	}

	return 0, fmt.Errorf("competition %d is not supported", competitionID)
}

func NewRatingProcessor(r RatingReader, w RatingWriter, c RatingCalculator, comp map[uint64]uint16) RatingProcessor {
	return &ratingProcessor{
		reader:     r,
		writer:     w,
		calculator: c,
		competitionMapping: comp,
	}
}
