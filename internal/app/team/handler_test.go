package team_test

import (
	"context"
	"errors"
	"github.com/jonboulle/clockwork"
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
		clock := clockwork.NewFakeClockAt(time.Unix(1615593600, 0))
		logger, _ := test.NewNullLogger()

		handler := team.NewHandler(fetcher, processor, clock, logger)

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
		clock := clockwork.NewFakeClockAt(time.Unix(1615593600, 0))
		logger, hook := test.NewNullLogger()

		handler := team.NewHandler(fetcher, processor, clock, logger)

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
		clock := clockwork.NewFakeClockAt(time.Unix(1615593600, 0))
		logger, hook := test.NewNullLogger()

		handler := team.NewHandler(fetcher, processor, clock, logger)

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

func TestRatingHandler_Today(t *testing.T) {
	t.Run("fetches and processes fixtures", func(t *testing.T) {
		t.Helper()

		fetcher := new(MockFixtureFetcher)
		processor := new(MockTeamRatingProcessor)
		clock := clockwork.NewFakeClockAt(time.Unix(1615629600, 0))
		logger, _ := test.NewNullLogger()

		handler := team.NewHandler(fetcher, processor, clock, logger)

		fix1 := statistico.Fixture{}
		fix2 := statistico.Fixture{}
		fix3 := statistico.Fixture{}

		fixtures := []*statistico.Fixture{
			&fix1,
			&fix2,
			&fix3,
		}

		ctx := context.Background()
		start := time.Date(2021, 03, 13, 0, 0, 0, 0, time.UTC)
		end := time.Date(2021, 03, 13, 5, 0, 0, 0, time.UTC)

		fetcher.On("ByDate", ctx, start, end).Return(fixtures, nil)

		processor.On("ByFixture", ctx, &fix1).Once().Return(nil)
		processor.On("ByFixture", ctx, &fix2).Once().Return(nil)
		processor.On("ByFixture", ctx, &fix3).Once().Return(nil)

		err := handler.Today(ctx, 5)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		fetcher.AssertExpectations(t)
		processor.AssertExpectations(t)
	})

	t.Run("logs an error if returned by fixture client", func(t *testing.T) {
		t.Helper()

		fetcher := new(MockFixtureFetcher)
		processor := new(MockTeamRatingProcessor)
		clock := clockwork.NewFakeClockAt(time.Unix(1615629600, 0))
		logger, hook := test.NewNullLogger()

		handler := team.NewHandler(fetcher, processor, clock, logger)

		e := errors.New("fixture fetcher error")

		ctx := context.Background()
		start := time.Date(2021, 03, 13, 0, 0, 0, 0, time.UTC)
		end := time.Date(2021, 03, 13, 5, 0, 0, 0, time.UTC)

		fetcher.On("ByDate", ctx, start, end).Return([]*statistico.Fixture{}, e)

		processor.AssertNotCalled(t, "ByFixture")

		err := handler.Today(ctx, 5)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "error fetching fixtures in team rating handler: fixture fetcher error", hook.LastEntry().Message)
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
		fetcher.AssertExpectations(t)
	})

	t.Run("logs error and continues execution if error returned by processor", func(t *testing.T) {
		t.Helper()

		fetcher := new(MockFixtureFetcher)
		processor := new(MockTeamRatingProcessor)
		clock := clockwork.NewFakeClockAt(time.Unix(1615629600, 0))
		logger, hook := test.NewNullLogger()

		handler := team.NewHandler(fetcher, processor, clock, logger)

		fix1 := statistico.Fixture{}
		fix2 := statistico.Fixture{}
		fix3 := statistico.Fixture{}

		fixtures := []*statistico.Fixture{
			&fix1,
			&fix2,
			&fix3,
		}

		ctx := context.Background()
		start := time.Date(2021, 03, 13, 0, 0, 0, 0, time.UTC)
		end := time.Date(2021, 03, 13, 5, 0, 0, 0, time.UTC)

		fetcher.On("ByDate", ctx, start, end).Return(fixtures, nil)

		e := errors.New("team rating processing error")

		processor.On("ByFixture", ctx, &fix1).Once().Return(nil)
		processor.On("ByFixture", ctx, &fix2).Once().Return(e)
		processor.On("ByFixture", ctx, &fix3).Once().Return(nil)

		err := handler.Today(ctx, 5)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		assert.Equal(t, "error processing fixtures in team rating handler: team rating processing error", hook.LastEntry().Message)
		assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
		fetcher.AssertExpectations(t)
		processor.AssertExpectations(t)
	})

	t.Run("returns an error if hour provided is greater than current time hour", func(t *testing.T) {
		t.Helper()

		fetcher := new(MockFixtureFetcher)
		processor := new(MockTeamRatingProcessor)
		clock := clockwork.NewFakeClockAt(time.Unix(1615629600, 0))
		logger, _ := test.NewNullLogger()

		handler := team.NewHandler(fetcher, processor, clock, logger)

		ctx := context.Background()

		fetcher.AssertNotCalled(t, "ByDate")
		processor.AssertNotCalled(t, "ByFixture")

		err := handler.Today(ctx, 23)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "hour provided is greater than the current time hour", err.Error())
	})
}

type MockFixtureFetcher struct {
	mock.Mock
}

func (m *MockFixtureFetcher) ByCompetition(ctx context.Context, competitionID, seasonID uint64) ([]*statistico.Fixture, error) {
	args := m.Called(ctx, competitionID, seasonID)
	return args.Get(0).([]*statistico.Fixture), args.Error(1)
}

func (m *MockFixtureFetcher) ByDate(ctx context.Context, from, to time.Time) ([]*statistico.Fixture, error) {
	args := m.Called(ctx, from, to)
	return args.Get(0).([]*statistico.Fixture), args.Error(1)
}

type MockTeamRatingProcessor struct {
	mock.Mock
}

func (m *MockTeamRatingProcessor) ByFixture(ctx context.Context, f *statistico.Fixture) error {
	args := m.Called(ctx, f)
	return args.Error(0)
}
