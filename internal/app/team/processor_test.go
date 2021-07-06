package team_test

import (
	"context"
	"errors"
	"github.com/statistico/statistico-proto/go"
	"github.com/statistico/statistico-ratings/internal/app/team"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestRatingProcessor_ByFixture(t *testing.T) {
	fixture := statistico.Fixture{
		HomeTeam:             &statistico.Team{Id: 5},
		AwayTeam:             &statistico.Team{Id: 6},
	}

	ctx := context.Background()

	t.Run("processes new ratings for home and away teams for a fixture", func(t *testing.T) {
		t.Helper()

		reader := new(MockRatingReader)
		writer := new(MockRatingWriter)
		calc := new(MockRatingCalculator)

		processor := team.NewRatingProcessor(reader, writer, calc)

		home := team.Rating{}
		away := team.Rating{}

		reader.On("Latest", uint64(5)).Return(&home, nil)
		reader.On("Latest", uint64(6)).Return(&away, nil)

		newHome := team.Rating{}
		newAway := team.Rating{}

		calc.On("ForFixture", ctx, &fixture, &home, &away).Return(&newHome, &newAway, nil)

		writer.On("Insert", &newHome).Return(nil)
		writer.On("Insert", &newAway).Return(nil)

		err := processor.ByFixture(ctx, &fixture)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		reader.AssertExpectations(t)
		writer.AssertExpectations(t)
		calc.AssertExpectations(t)
	})

	t.Run("returns error if returned by rating reader", func(t *testing.T) {
		t.Helper()

		reader := new(MockRatingReader)
		writer := new(MockRatingWriter)
		calc := new(MockRatingCalculator)

		processor := team.NewRatingProcessor(reader, writer, calc)

		e := errors.New("rating reader error")

		reader.On("Latest", uint64(5)).Return(&team.Rating{}, e)

		reader.AssertNotCalled(t, "Latest", uint64(6))
		calc.AssertNotCalled(t, "ForFixture")
		writer.AssertNotCalled(t, "Insert")

		err := processor.ByFixture(ctx, &fixture)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "rating reader error", e.Error())
	})

	t.Run("returns error if returned by rating calculator", func(t *testing.T) {
		t.Helper()

		reader := new(MockRatingReader)
		writer := new(MockRatingWriter)
		calc := new(MockRatingCalculator)

		processor := team.NewRatingProcessor(reader, writer, calc)

		e := errors.New("rating calculator error")

		home := team.Rating{}
		away := team.Rating{}

		reader.On("Latest", uint64(5)).Return(&home, nil)
		reader.On("Latest", uint64(6)).Return(&away, nil)

		calc.On("ForFixture", ctx, &fixture, &home, &away).Return(&team.Rating{}, &team.Rating{}, e)

		writer.AssertNotCalled(t, "Insert")

		err := processor.ByFixture(ctx, &fixture)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		reader.AssertExpectations(t)
		calc.AssertExpectations(t)

		assert.Equal(t, "rating calculator error", e.Error())
	})

	t.Run("returns error if returned by rating writer", func(t *testing.T) {
		t.Helper()

		reader := new(MockRatingReader)
		writer := new(MockRatingWriter)
		calc := new(MockRatingCalculator)

		processor := team.NewRatingProcessor(reader, writer, calc)

		e := errors.New("rating writer error")

		home := team.Rating{}
		away := team.Rating{}

		reader.On("Latest", uint64(5)).Return(&home, nil)
		reader.On("Latest", uint64(6)).Return(&away, nil)

		newHome := team.Rating{}
		newAway := team.Rating{}

		calc.On("ForFixture", ctx, &fixture, &home, &away).Return(&newHome, &newAway, nil)

		writer.On("Insert", &newHome).Return(e)

		err := processor.ByFixture(ctx, &fixture)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		reader.AssertExpectations(t)
		calc.AssertExpectations(t)

		assert.Equal(t, "rating writer error", e.Error())
	})
}

type MockRatingReader struct {
	mock.Mock
}

func (m *MockRatingReader) Latest(teamID uint64) (*team.Rating, error) {
	args := m.Called(teamID)
	return args.Get(0).(*team.Rating), args.Error(1)
}

type MockRatingWriter struct {
	mock.Mock
}

func (m *MockRatingWriter) Insert(r *team.Rating) error {
	args := m.Called(r)
	return args.Error(0)
}

type MockRatingCalculator struct {
	mock.Mock
}

func (m *MockRatingCalculator) ForFixture(ctx context.Context, f *statistico.Fixture, h, a *team.Rating) (*team.Rating, *team.Rating, error) {
	args := m.Called(ctx, f, h, a)
	return args.Get(0).(*team.Rating), args.Get(1).(*team.Rating), args.Error(2)
}
