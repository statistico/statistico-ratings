package team_test

import (
	"github.com/statistico/statistico-ratings/internal/app/team"
	"github.com/statistico/statistico-ratings/internal/app/test"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRatingReader_Latest(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, []string{"team_rating"})
	writer := team.NewRatingWriter(conn)
	reader := team.NewRatingReader(conn)

	t.Run("returns last rating record for a team", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		s := []struct {
			Rating *team.Rating
		}{
			{
				&team.Rating{
					TeamID:    1,
					FixtureID: 65,
					SeasonID:  17462,
					Attack: team.Points{
						Total:      1728,
						Difference: -3,
					},
					Defence: team.Points{
						Total:      1241,
						Difference: 4,
					},
					Timestamp: time.Unix(1625169423, 0),
				},
			},
			{
				&team.Rating{
					TeamID:    1,
					FixtureID: 55,
					SeasonID:  17462,
					Attack: team.Points{
						Total:      1901,
						Difference: 24,
					},
					Defence: team.Points{
						Total:      1023,
						Difference: 2,
					},
					Timestamp: time.Unix(1624169423, 0),
				},
			},
			{
				&team.Rating{
					TeamID:    1,
					FixtureID: 120,
					SeasonID:  17462,
					Attack: team.Points{
						Total:      2810,
						Difference: 13,
					},
					Defence: team.Points{
						Total:      1100,
						Difference: -25,
					},
					Timestamp: time.Unix(1622169423, 0),
				},
			},
		}

		for _, st := range s {
			if err := writer.Insert(st.Rating); err != nil {
				t.Fatalf("Expected nil, got %s", err.Error())
			}
		}

		fetched, err := reader.Latest(1)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		a := assert.New(t)

		a.Equal(uint64(1), fetched.TeamID)
		a.Equal(uint64(65), fetched.FixtureID)
		a.Equal(uint64(17462), fetched.SeasonID)
		a.Equal(float64(1728), fetched.Attack.Total)
		a.Equal(float64(-3), fetched.Attack.Difference)
		a.Equal(float64(1241), fetched.Defence.Total)
		a.Equal(float64(4), fetched.Defence.Difference)
		a.Equal(time.Unix(1625169423, 0), fetched.Timestamp)
	})

	t.Run("returns a NotFoundError if team rating does not exist", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		_, err := reader.Latest(1)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "team 1 rating does not exist", err.Error())
	})
}

func TestRatingReader_Get(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, []string{"team_rating"})
	writer := team.NewRatingWriter(conn)
	reader := team.NewRatingReader(conn)

	t.Run("returns a slice of Rating struct", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		team1 := uint64(1)
		team2 := uint64(2)
		season := uint64(10)
		date := time.Unix(1625162423, 0)

		s := []struct {
			Query  *team.ReaderQuery
			Count  int
		}{
			{
				&team.ReaderQuery{TeamID: &team1},
				3,
			},
			{
				&team.ReaderQuery{
					TeamID:   &team2,
					SeasonID: &season,
				},
				1,
			},
			{
				&team.ReaderQuery{
					TeamID:   &team2,
					Before:   &date,
				},
				1,
			},
		}

		insertRatings(t, writer)

		for _, st := range s {
			ratings, err := reader.Get(st.Query)

			if err != nil {
				t.Fatalf("Expected nil, got %s", err.Error())
			}

			assert.Equal(t, st.Count, len(ratings))
		}
	})

	t.Run("returns an sorted list of team rating struct", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		insertRatings(t, writer)

		teamID := uint64(1)

		query := team.ReaderQuery{
			TeamID:   &teamID,
			Sort:     "timestamp_desc",
		}

		ratings, err := reader.Get(&query)

		if err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		a := assert.New(t)

		a.Equal(3, len(ratings))
		a.Equal(uint64(67), ratings[0].FixtureID)
		a.Equal(uint64(66), ratings[1].FixtureID)
		a.Equal(uint64(65), ratings[2].FixtureID)
	})
}

func insertRatings(t *testing.T, r team.RatingWriter) {
	ratings := []*team.Rating{
		{
			TeamID:    1,
			FixtureID: 65,
			SeasonID:  9,
			Attack: team.Points{
				Total:      1728,
				Difference: -3,
			},
			Defence: team.Points{
				Total:      1241,
				Difference: 4,
			},
			Timestamp: time.Unix(1625162423, 0),
		},
		{
			TeamID:    1,
			FixtureID: 66,
			SeasonID:  9,
			Attack: team.Points{
				Total:      1728,
				Difference: -3,
			},
			Defence: team.Points{
				Total:      1241,
				Difference: 4,
			},
			Timestamp: time.Unix(1625163423, 0),
		},
		{
			TeamID:    1,
			FixtureID: 67,
			SeasonID:  9,
			Attack: team.Points{
				Total:      1728,
				Difference: -3,
			},
			Defence: team.Points{
				Total:      1241,
				Difference: 4,
			},
			Timestamp: time.Unix(1625164423, 0),
		},
		{
			TeamID:    2,
			FixtureID: 66,
			SeasonID:  9,
			Attack: team.Points{
				Total:      1728,
				Difference: -3,
			},
			Defence: team.Points{
				Total:      1241,
				Difference: 4,
			},
			Timestamp: time.Unix(1625162423, 0),
		},
		{
			TeamID:    2,
			FixtureID: 70,
			SeasonID:  10,
			Attack: team.Points{
				Total:      1728,
				Difference: -3,
			},
			Defence: team.Points{
				Total:      1241,
				Difference: 4,
			},
			Timestamp: time.Unix(1625172423, 0),
		},
	}

	for _, rating := range ratings {
		if err := r.Insert(rating); err != nil {
			t.Fatalf("Error inserting team rating: %s", err.Error())
		}
	}
}