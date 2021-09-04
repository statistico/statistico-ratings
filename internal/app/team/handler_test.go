package team_test

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-ratings/internal/app/team"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestRatingHandler_ByCompetition(t *testing.T) {
	t.Run("fetches and processes fixtures", func(t *testing.T) {
		t.Helper()

		fetcher := new(MockFixtureFetcher)
		processor := new(MockTeamRatingProcessor)
		logger, _ := test.NewNullLogger()

		handler := team.NewHandler(fetcher, processor, logger)

		fix1 := statistico.Fixture{}
		fix2 := statistico.Fixture{}
		fix3 := statistico.Fixture{}

		fixtures := []*statistico.Fixture{
			&fix1,
			&fix2,
			&fix3,
		}

		ctx := context.Background()

		fetcher.On("ByCompetition", ctx, uint64(8), uint64(4)).Return(fixtures, nil)

		processor.On("ByFixture", ctx, &fix1).Once().Return(nil)
		processor.On("ByFixture", ctx, &fix2).Once().Return(nil)
		processor.On("ByFixture", ctx, &fix3).Once().Return(nil)

		handler.ByCompetition(ctx, uint64(8), uint64(4))

		fetcher.AssertExpectations(t)
		processor.AssertExpectations(t)
	})

	t.Run("logs an error if returned by fixture client", func(t *testing.T) {
		t.Helper()

		fetcher := new(MockFixtureFetcher)
		processor := new(MockTeamRatingProcessor)
		logger, hook := test.NewNullLogger()

		handler := team.NewHandler(fetcher, processor, logger)

		e := errors.New("fixture fetcher error")

		ctx := context.Background()

		fetcher.On("ByCompetition", ctx, uint64(8), uint64(4)).Return([]*statistico.Fixture{}, e)

		processor.AssertNotCalled(t, "ByFixture")

		handler.ByCompetition(ctx, uint64(8), uint64(4))

		assert.Equal(t, "error fetching fixtures in team rating handler: fixture fetcher error", hook.LastEntry().Message)
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
		fetcher.AssertExpectations(t)
	})

	t.Run("logs error and exits execution if error returned by processor", func(t *testing.T) {
		t.Helper()

		fetcher := new(MockFixtureFetcher)
		processor := new(MockTeamRatingProcessor)
		logger, hook := test.NewNullLogger()

		handler := team.NewHandler(fetcher, processor, logger)

		fix1 := statistico.Fixture{}
		fix2 := statistico.Fixture{}
		fix3 := statistico.Fixture{}

		fixtures := []*statistico.Fixture{
			&fix1,
			&fix2,
			&fix3,
		}

		ctx := context.Background()

		fetcher.On("ByCompetition", ctx, uint64(8), uint64(4)).Return(fixtures, nil)

		e := errors.New("team rating processing error")

		processor.On("ByFixture", ctx, &fix1).Once().Return(nil)
		processor.On("ByFixture", ctx, &fix2).Once().Return(e)

		processor.AssertNotCalled(t, "ByFixture", ctx, &fix3)

		handler.ByCompetition(ctx, uint64(8), uint64(4))

		assert.Equal(t, "error processing fixtures in team rating handler: team rating processing error", hook.LastEntry().Message)
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
		fetcher.AssertExpectations(t)
		processor.AssertExpectations(t)
	})
}

func TestRatingHandler_ByDate(t *testing.T) {
	t.Run("fetches and processes fixtures", func(t *testing.T) {
		t.Helper()

		fetcher := new(MockFixtureFetcher)
		processor := new(MockTeamRatingProcessor)
		logger, _ := test.NewNullLogger()

		handler := team.NewHandler(fetcher, processor, logger)

		fix1 := statistico.Fixture{}
		fix2 := statistico.Fixture{}
		fix3 := statistico.Fixture{}

		fixtures := []*statistico.Fixture{
			&fix1,
			&fix2,
			&fix3,
		}

		ctx := context.Background()
		date := time.Now()

		fetcher.On("ByDate", ctx, date).Return(fixtures, nil)

		processor.On("ByFixture", ctx, &fix1).Once().Return(nil)
		processor.On("ByFixture", ctx, &fix2).Once().Return(nil)
		processor.On("ByFixture", ctx, &fix3).Once().Return(nil)

		handler.ByDate(ctx, date)

		fetcher.AssertExpectations(t)
		processor.AssertExpectations(t)
	})

	t.Run("logs an error if returned by fixture client", func(t *testing.T) {
		t.Helper()

		fetcher := new(MockFixtureFetcher)
		processor := new(MockTeamRatingProcessor)
		logger, hook := test.NewNullLogger()

		handler := team.NewHandler(fetcher, processor, logger)

		e := errors.New("fixture fetcher error")

		ctx := context.Background()
		date := time.Now()

		fetcher.On("ByDate", ctx, date).Return([]*statistico.Fixture{}, e)

		processor.AssertNotCalled(t, "ByFixture")

		handler.ByDate(ctx, date)

		assert.Equal(t, "error fetching fixtures in team rating handler: fixture fetcher error", hook.LastEntry().Message)
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
		fetcher.AssertExpectations(t)
	})

	t.Run("logs error and continues execution if error returned by processor", func(t *testing.T) {
		t.Helper()

		fetcher := new(MockFixtureFetcher)
		processor := new(MockTeamRatingProcessor)
		logger, hook := test.NewNullLogger()

		handler := team.NewHandler(fetcher, processor, logger)

		fix1 := statistico.Fixture{}
		fix2 := statistico.Fixture{}
		fix3 := statistico.Fixture{}

		fixtures := []*statistico.Fixture{
			&fix1,
			&fix2,
			&fix3,
		}

		ctx := context.Background()
		date := time.Now()

		fetcher.On("ByDate", ctx, date).Return(fixtures, nil)

		e := errors.New("team rating processing error")

		processor.On("ByFixture", ctx, &fix1).Once().Return(nil)
		processor.On("ByFixture", ctx, &fix2).Once().Return(e)
		processor.On("ByFixture", ctx, &fix3).Once().Return(nil)

		handler.ByDate(ctx, date)

		assert.Equal(t, "error processing fixtures in team rating handler: team rating processing error", hook.LastEntry().Message)
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
		fetcher.AssertExpectations(t)
		processor.AssertExpectations(t)
	})
}

type MockFixtureFetcher struct {
	mock.Mock
}

func (m *MockFixtureFetcher) ByCompetition(ctx context.Context, competitionID, seasonID uint64) ([]*statistico.Fixture, error) {
	args := m.Called(ctx, competitionID, seasonID)
	return args.Get(0).([]*statistico.Fixture), args.Error(1)
}

func (m *MockFixtureFetcher) ByDate(ctx context.Context, date time.Time) ([]*statistico.Fixture, error) {
	args := m.Called(ctx, date)
	return args.Get(0).([]*statistico.Fixture), args.Error(1)
}

type MockTeamRatingProcessor struct {
	mock.Mock
}

func (m *MockTeamRatingProcessor) ByFixture(ctx context.Context, f *statistico.Fixture) error {
	args := m.Called(ctx, f)
	return args.Error(0)
}
