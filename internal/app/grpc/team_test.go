package grpc_test

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-ratings/internal/app/grpc"
	"github.com/statistico/statistico-ratings/internal/app/team"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestTeamRatingService_GetTeamRatings(t *testing.T) {
	t.Run("calls team rating reader and returns statistico team rating response", func(t *testing.T) {
		t.Helper()

		reader := new(MockTeamRatingReader)
		logger, _ := test.NewNullLogger()

		service := grpc.NewTeamRatingService(reader, logger)

		req := statistico.TeamRatingRequest{
			TeamId:     5,
			SeasonId:   &wrappers.UInt64Value{Value: 99},
			DateBefore: &wrappers.StringValue{Value: "2021-03-12T12:00:00+00:00"},
			Sort:       "timestamp_asc",
		}

		ratings := []*team.Rating{
			{
				TeamID:    5,
				FixtureID: 10,
				SeasonID:  99,
				Attack: team.Points{
					Total:      1432.12,
					Difference: 10,
				},
				Defence: team.Points{
					Total:      234,
					Difference: -23,
				},
				FixtureDate: time.Unix(1627226510, 0),
				Timestamp:   time.Unix(1627226510, 0),
			},
			{
				TeamID:    5,
				FixtureID: 11,
				SeasonID:  99,
				Attack: team.Points{
					Total:      1452.45,
					Difference: -2,
				},
				Defence: team.Points{
					Total:      225,
					Difference: -3,
				},
				FixtureDate: time.Unix(1627226510, 0),
				Timestamp:   time.Unix(1627226510, 0),
			},
		}

		query := mock.MatchedBy(func(q *team.ReaderQuery) bool {
			a := assert.New(t)

			a.Equal(uint64(5), *q.TeamID)
			a.Equal(uint64(99), *q.SeasonID)
			a.Equal("2021-03-12T12:00:00Z", q.Before.Format(time.RFC3339))
			a.Equal("timestamp_asc", q.Sort)
			return true
		})

		reader.On("Get", query).Return(ratings, nil)

		res, err := service.GetTeamRatings(context.Background(), &req)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		one := res.Ratings[0]
		two := res.Ratings[1]

		a := assert.New(t)

		a.Equal(uint64(5), one.TeamId)
		a.Equal(uint64(99), one.SeasonId)
		a.Equal(uint64(10), one.FixtureId)
		a.Equal(float32(1432.12), one.Attack.Points)
		a.Equal(float32(10), one.Attack.Difference)
		a.Equal(float32(234), one.Defence.Points)
		a.Equal(float32(-23), one.Defence.Difference)
		a.Equal(int64(1627226510), one.FixtureDate.GetSeconds())
		a.Equal(int64(1627226510), one.Timestamp.GetSeconds())

		a.Equal(uint64(5), two.TeamId)
		a.Equal(uint64(99), two.SeasonId)
		a.Equal(uint64(11), two.FixtureId)
		a.Equal(float32(1452.45), two.Attack.Points)
		a.Equal(float32(-2), two.Attack.Difference)
		a.Equal(float32(225), two.Defence.Points)
		a.Equal(float32(-3), two.Defence.Difference)
		a.Equal(int64(1627226510), two.FixtureDate.GetSeconds())
		a.Equal(int64(1627226510), two.Timestamp.GetSeconds())

		reader.AssertExpectations(t)
	})

	t.Run("returns an invalid argument error if date provided in request is in the wrong format", func(t *testing.T) {
		t.Helper()

		reader := new(MockTeamRatingReader)
		logger, _ := test.NewNullLogger()

		service := grpc.NewTeamRatingService(reader, logger)

		req := statistico.TeamRatingRequest{
			TeamId:     5,
			SeasonId:   &wrappers.UInt64Value{Value: 99},
			DateBefore: &wrappers.StringValue{Value: "hello"},
			Sort:       "timestamp_asc",
		}

		reader.AssertNotCalled(t, "Get")

		_, err := service.GetTeamRatings(context.Background(), &req)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		a := assert.New(t)

		a.Equal("rpc error: code = InvalidArgument desc = parsing time \"hello\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"hello\" as \"2006\"", err.Error())
	})

	t.Run("logs error and returns internal server error if error returned by team rating reader", func(t *testing.T) {
		t.Helper()

		reader := new(MockTeamRatingReader)
		logger, hook := test.NewNullLogger()

		service := grpc.NewTeamRatingService(reader, logger)

		req := statistico.TeamRatingRequest{
			TeamId: 5,
			Sort:   "timestamp_asc",
		}

		query := mock.MatchedBy(func(q *team.ReaderQuery) bool {
			a := assert.New(t)

			a.Equal(uint64(5), *q.TeamID)
			a.Nil(q.SeasonID)
			a.Nil(q.Before)
			a.Equal("timestamp_asc", q.Sort)
			return true
		})

		reader.On("Get", query).Return([]*team.Rating{}, errors.New("oh no"))

		_, err := service.GetTeamRatings(context.Background(), &req)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		a := assert.New(t)

		a.Equal("rpc error: code = Internal desc = internal server error", err.Error())
		assert.Equal(t, "Error fetching team ratings: oh no", hook.LastEntry().Message)
		assert.Equal(t, 1, len(hook.Entries))
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	})
}

type MockTeamRatingReader struct {
	mock.Mock
}

func (m *MockTeamRatingReader) Get(q *team.ReaderQuery) ([]*team.Rating, error) {
	args := m.Called(q)
	return args.Get(0).([]*team.Rating), args.Error(1)
}

func (m *MockTeamRatingReader) Latest(teamID uint64) (*team.Rating, error) {
	args := m.Called(teamID)
	return args.Get(0).(*team.Rating), args.Error(1)
}
