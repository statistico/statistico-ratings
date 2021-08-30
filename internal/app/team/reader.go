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
	b := queryBuilder(r.connection)

	query := b.
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
		From("team_rating")

	rows, err := buildQuery(query, q).Query()

	if err != nil {
		return []*Rating{}, err
	}

	return rowsToRatingSlice(rows)
}

func rowsToRatingSlice(rows *sql.Rows) ([]*Rating, error) {
	var ratings []*Rating
	var timestamp int64

	for rows.Next() {
		var rating Rating
		var attack Points
		var defence Points

		err := rows.Scan(
			&rating.TeamID,
			&rating.FixtureID,
			&rating.SeasonID,
			&attack.Total,
			&attack.Difference,
			&defence.Total,
			&defence.Difference,
			&timestamp,
		)

		if err != nil {
			return ratings, err
		}

		rating.Attack = attack
		rating.Defence = defence
		rating.Timestamp = time.Unix(timestamp, 0)

		ratings = append(ratings, &rating)
	}

	err := rows.Close()

	if err != nil {
		return ratings, err
	}

	return ratings, nil
}

func buildQuery(b sq.SelectBuilder, q *ReaderQuery) sq.SelectBuilder {
	if q.TeamID != nil {
		b = b.Where(sq.Eq{"team_id": q.TeamID})
	}

	if q.SeasonID != nil {
		b = b.Where(sq.Eq{"season_id": q.SeasonID})
	}

	if q.Before != nil {
		b = b.Where(sq.LtOrEq{"timestamp": q.Before.Unix()})
	}

	if q.Sort == "timestamp_asc" {
		b = b.OrderBy("timestamp ASC")
	}

	if q.Sort == "timestamp_desc" {
		b = b.OrderBy("timestamp DESC")
	}

	return b
}

func NewRatingReader(c *sql.DB) RatingReader {
	return &ratingReader{connection: c}
}
