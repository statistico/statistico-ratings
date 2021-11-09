package grpc

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-ratings/internal/app/team"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type TeamRatingService struct {
	reader team.RatingReader
	logger *logrus.Logger
	statistico.UnimplementedTeamRatingServiceServer
}

func (t *TeamRatingService) GetTeamRatings(ctx context.Context, r *statistico.TeamRatingRequest) (*statistico.TeamRatingResponse, error) {
	q, err := buildTeamReaderQuery(r)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ratings, err := t.reader.Get(q)

	if err != nil {
		t.logger.Errorf("Error fetching team ratings: %s", err.Error())
		return nil, status.Error(codes.Internal, "internal server error")
	}

	res := statistico.TeamRatingResponse{}

	for _, rt := range ratings {
		t.logger.Errorf("Attack points: %+v", rt.Attack.Total)
		st := statistico.TeamRating{
			TeamId:    rt.TeamID,
			FixtureId: rt.FixtureID,
			SeasonId:  rt.SeasonID,
			Attack: &statistico.Points{
				Points:     float32(rt.Attack.Total),
				Difference: float32(rt.Attack.Difference),
			},
			Defence: &statistico.Points{
				Points:     float32(rt.Defence.Total),
				Difference: float32(rt.Defence.Difference),
			},
			FixtureDate: timestamppb.New(rt.FixtureDate),
			Timestamp:   timestamppb.New(rt.Timestamp),
		}

		res.Ratings = append(res.Ratings, &st)
	}

	return &res, nil
}

func buildTeamReaderQuery(r *statistico.TeamRatingRequest) (*team.ReaderQuery, error) {
	q := team.ReaderQuery{
		TeamID: &r.TeamId,
		Sort:   r.Sort,
	}

	if r.SeasonId != nil {
		q.SeasonID = &r.SeasonId.Value
	}

	if r.DateBefore != nil {
		t, err := time.Parse(time.RFC3339, r.DateBefore.Value)

		if err != nil {
			return nil, err
		}

		q.Before = &t
	}

	return &q, nil
}

func NewTeamRatingService(r team.RatingReader, l *logrus.Logger) *TeamRatingService {
	return &TeamRatingService{reader: r, logger: l}
}
