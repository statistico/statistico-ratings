package team

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/statistico/statistico-ratings/internal/app"
	"time"
)

type RatingReader interface {
	Latest(teamID uint64) (*Rating, error)
	Get(q *ReaderQuery) ([]*Rating, error)
}

type ratingReader struct {
	connection *sql.DB
}

func (r *ratingReader) Latest(teamID uint64) (*Rating, error) {
	b := queryBuilder(r.connection)

	var rating Rating
	var timestamp int64

	row := b.
		Select(
			"team_id",
			"fixture_id",
			"season_id",
			"attack_total",
			"attack_points",
			"defence_total",
			"defence_points",
			"timestamp",
		).
		From("team_rating").
		Where(sq.Eq{"team_id": teamID}).
		OrderBy("timestamp DESC").
		OrderBy("id DESC").
		Limit(1).
		QueryRow().
		Scan(
			&rating.TeamID,
			&rating.FixtureID,
			&rating.SeasonID,
			&rating.Attack.Total,
			&rating.Attack.Difference,
			&rating.Defence.Total,
			&rating.Defence.Difference,
			&timestamp,
		)

	if row != nil {
		return nil, &app.NotFoundError{TeamID: teamID}
	}

	rating.Timestamp = time.Unix(timestamp, 0)

	return &rating, nil
}

func (r *ratingReader) Get(q *ReaderQuery) ([]*Rating, error) {
	ratings := []*Rating{}

	return ratings, nil
}

func NewRatingReader(c *sql.DB) RatingReader {
	return &ratingReader{connection: c}
}
