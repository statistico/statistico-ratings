package team_test

import (
	"github.com/statistico/statistico-ratings/internal/app/team"
	"github.com/statistico/statistico-ratings/internal/app/test"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRatingRepository_Insert(t *testing.T) {
	conn, cleanUp := test.GetConnection(t, []string{"team_rating"})
	writer := team.NewRatingWriter(conn)

	t.Run("increases table count", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		s := []struct {
			Rating *team.Rating
			Count  int8
		}{
			{
				&team.Rating{
					TeamID:    1,
					FixtureID: 55,
					SeasonID:  17462,
					Attack: team.Points{
						Total:      1728,
						Difference: -3,
					},
					Defence: team.Points{
						Total:      1241,
						Difference: 4,
					},
					FixtureDate: time.Unix(1630343736, 0),
					Timestamp: time.Now(),
				},
				1,
			},
			{
				&team.Rating{
					TeamID:    2,
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
					FixtureDate: time.Unix(1630243736, 0),
					Timestamp: time.Now(),
				},
				2,
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
					FixtureDate: time.Unix(1630143736, 0),
					Timestamp: time.Now(),
				},
				3,
			},
		}

		for _, st := range s {
			if err := writer.Insert(st.Rating); err != nil {
				t.Fatalf("Expected nil, got %s", err.Error())
			}

			var count int8

			row := conn.QueryRow("select count(*) from team_rating")

			if err := row.Scan(&count); err != nil {
				t.Errorf("Error when scanning rows returned by the database: %s", err.Error())
			}

			assert.Equal(t, st.Count, count)
		}
	})

	t.Run("returns a DuplicationError if record exists for team, fixture and season", func(t *testing.T) {
		t.Helper()
		defer cleanUp()

		r := &team.Rating{
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
			FixtureDate: time.Now(),
			Timestamp: time.Now(),
		}

		if err := writer.Insert(r); err != nil {
			t.Fatalf("Expected nil, got %s", err.Error())
		}

		err := writer.Insert(r)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		assert.Equal(t, "team rating exists for team 1, fixture 120 and season 17462", err.Error())
	})
}
