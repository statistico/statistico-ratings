package fixture_test

import (
	"context"
	"errors"
	"github.com/jonboulle/clockwork"
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-ratings/internal/app/fixture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestFetcher_ByCompetition(t *testing.T) {
	t.Run("fetches and returns fixtures for a competition", func(t *testing.T) {
		t.Helper()

		competitions := []uint64{8}
		fixtureClient := new(MockFixtureClient)
		seasonClient := new(MockSeasonClient)
		clock := clockwork.NewFakeClock()

		ctx := context.Background()

		fetcher := fixture.NewFetcher(competitions, fixtureClient, seasonClient, clock)

		seasonClient.On("ByCompetitionID", ctx, uint64(8), "name_desc").Return(seasonResponse(), nil)

		req := mock.MatchedBy(func(r *statistico.FixtureSearchRequest) bool {
			assert.Equal(t, []uint64{2, 3}, r.SeasonIds)
			assert.Equal(t, "1984-04-04T00:00:00Z", r.DateBefore.GetValue())
			assert.Equal(t, "date_asc", r.Sort.GetValue())
			return true
		})

		response := fixtureResponse()

		fixtureClient.On("Search", ctx, req).Return(response, nil)

		fixtures, err := fetcher.ByCompetition(ctx, uint64(8), int8(2))

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		assert.Equal(t, response, fixtures)
		seasonClient.AssertExpectations(t)
		fixtureClient.AssertExpectations(t)
	})

	t.Run("returns an error if returned by season client", func(t *testing.T) {
		t.Helper()

		competitions := []uint64{8}
		fixtureClient := new(MockFixtureClient)
		seasonClient := new(MockSeasonClient)
		clock := clockwork.NewFakeClock()

		ctx := context.Background()

		fetcher := fixture.NewFetcher(competitions, fixtureClient, seasonClient, clock)

		e := errors.New("season client error")

		seasonClient.On("ByCompetitionID", ctx, uint64(8), "name_desc").Return([]*statistico.Season{}, e)

		fixtureClient.AssertNotCalled(t, "Search")

		_, err := fetcher.ByCompetition(ctx, uint64(8), int8(2))

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "season client error", err.Error())
		seasonClient.AssertExpectations(t)
		fixtureClient.AssertExpectations(t)
	})

	t.Run("returns an error if returned by fixture client", func(t *testing.T) {
		t.Helper()

		competitions := []uint64{8}
		fixtureClient := new(MockFixtureClient)
		seasonClient := new(MockSeasonClient)
		clock := clockwork.NewFakeClock()

		ctx := context.Background()

		fetcher := fixture.NewFetcher(competitions, fixtureClient, seasonClient, clock)

		e := errors.New("fixture client error")

		seasonClient.On("ByCompetitionID", ctx, uint64(8), "name_desc").Return(seasonResponse(), nil)

		req := mock.MatchedBy(func(r *statistico.FixtureSearchRequest) bool {
			assert.Equal(t, []uint64{2, 3}, r.SeasonIds)
			assert.Equal(t, "1984-04-04T00:00:00Z", r.DateBefore.GetValue())
			assert.Equal(t, "date_asc", r.Sort.GetValue())
			return true
		})

		fixtureClient.On("Search", ctx, req).Return([]*statistico.Fixture{}, e)

		_, err := fetcher.ByCompetition(ctx, uint64(8), int8(2))

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "fixture client error", err.Error())
		seasonClient.AssertExpectations(t)
		fixtureClient.AssertExpectations(t)
	})
}

func TestFetcher_ByDate(t *testing.T) {
	t.Run("fetches and returns fixtures by date", func(t *testing.T) {
		t.Helper()

		competitions := []uint64{8}
		fixtureClient := new(MockFixtureClient)
		seasonClient := new(MockSeasonClient)
		clock := clockwork.NewFakeClock()

		ctx := context.Background()

		fetcher := fixture.NewFetcher(competitions, fixtureClient, seasonClient, clock)

		seasonClient.
			On("ByCompetitionID", ctx, uint64(8), "name_desc").
			Once().
			Return(seasonResponse(), nil)

		req := mock.MatchedBy(func(r *statistico.FixtureSearchRequest) bool {
			assert.Equal(t, []uint64{3}, r.SeasonIds)
			assert.Equal(t, "2021-07-11T00:00:00Z", r.DateAfter.GetValue())
			assert.Equal(t, "2021-07-11T23:59:59Z", r.DateBefore.GetValue())
			assert.Equal(t, "date_asc", r.Sort.GetValue())
			return true
		})

		response := fixtureResponse()

		fixtureClient.On("Search", ctx, req).Return(response, nil)

		fixtures, err := fetcher.ByDate(ctx, time.Unix(1626008664, 0))

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		assert.Equal(t, response, fixtures)
		seasonClient.AssertExpectations(t)
		fixtureClient.AssertExpectations(t)
	})

	t.Run("returns an error if returned by season client", func(t *testing.T) {
		t.Helper()

		competitions := []uint64{8}
		fixtureClient := new(MockFixtureClient)
		seasonClient := new(MockSeasonClient)
		clock := clockwork.NewFakeClock()

		ctx := context.Background()

		fetcher := fixture.NewFetcher(competitions, fixtureClient, seasonClient, clock)

		e := errors.New("season client error")

		seasonClient.
			On("ByCompetitionID", ctx, uint64(8), "name_desc").
			Once().
			Return([]*statistico.Season{}, e)

		fixtureClient.AssertNotCalled(t, "Search")

		_, err := fetcher.ByDate(ctx, time.Unix(1626008664, 0))

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "season client error", err.Error())
		seasonClient.AssertExpectations(t)
		fixtureClient.AssertExpectations(t)
	})

	t.Run("returns an error if returned by fixture client", func(t *testing.T) {
		t.Helper()

		competitions := []uint64{8}
		fixtureClient := new(MockFixtureClient)
		seasonClient := new(MockSeasonClient)
		clock := clockwork.NewFakeClock()

		ctx := context.Background()

		fetcher := fixture.NewFetcher(competitions, fixtureClient, seasonClient, clock)

		e := errors.New("fixture client error")

		seasonClient.
			On("ByCompetitionID", ctx, uint64(8), "name_desc").
			Once().
			Return(seasonResponse(), nil)

		req := mock.MatchedBy(func(r *statistico.FixtureSearchRequest) bool {
			assert.Equal(t, []uint64{3}, r.SeasonIds)
			assert.Equal(t, "2021-07-11T00:00:00Z", r.DateAfter.GetValue())
			assert.Equal(t, "2021-07-11T23:59:59Z", r.DateBefore.GetValue())
			assert.Equal(t, "date_asc", r.Sort.GetValue())
			return true
		})

		fixtureClient.On("Search", ctx, req).Return([]*statistico.Fixture{}, e)

		_, err := fetcher.ByDate(ctx, time.Unix(1626008664, 0))

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "fixture client error", err.Error())
		seasonClient.AssertExpectations(t)
		fixtureClient.AssertExpectations(t)
	})
}

func seasonResponse() []*statistico.Season {
	return []*statistico.Season{
		{
			Id: 1,
			Name: "2018/2019",
		},
		{
			Id: 2,
			Name: "2019/2020",
		},
		{
			Id: 3,
			Name: "2020/2021",
		},
	}
}

func fixtureResponse() []*statistico.Fixture {
	return []*statistico.Fixture{}
}

type MockFixtureClient struct {
	mock.Mock
}

func (m *MockFixtureClient) Search(ctx context.Context, req *statistico.FixtureSearchRequest) ([]*statistico.Fixture, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]*statistico.Fixture), args.Error(1)
}

func (m *MockFixtureClient) ByID(ctx context.Context, fixtureID uint64) (*statistico.Fixture, error) {
	args := m.Called(ctx, fixtureID)
	return args.Get(0).(*statistico.Fixture), args.Error(1)
}

type MockSeasonClient struct {
	mock.Mock
}

func (m *MockSeasonClient) ByTeamID(ctx context.Context, teamId uint64, sort string) ([]*statistico.Season, error) {
	args := m.Called(ctx, teamId, sort)
	return args.Get(0).([]*statistico.Season), args.Error(1)
}

func (m *MockSeasonClient) ByCompetitionID(ctx context.Context, competitionId uint64, sort string) ([]*statistico.Season, error) {
	args := m.Called(ctx, competitionId, sort)
	return args.Get(0).([]*statistico.Season), args.Error(1)
}
