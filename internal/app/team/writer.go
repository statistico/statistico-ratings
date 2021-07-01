package team

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/statistico/statistico-ratings/internal/app"
)

type RatingWriter interface {
	Insert(r *Rating) error
}

type ratingWriter struct {
	connection *sql.DB
}

func (r *ratingWriter) Insert(x *Rating) error {
	var exists bool

	b := queryBuilder(r.connection)

	err := b.
		Select("id").
		From("team_rating").
		Prefix("SELECT exists (").
		Where(sq.Eq{"team_id": x.TeamID}).
		Where(sq.Eq{"season_id": x.SeasonID}).
		Where(sq.Eq{"fixture_id": x.FixtureID}).
		Suffix(")").
		Scan(&exists)

	if err != nil {
		return err
	}

	if exists {
		return &app.DuplicationError{
			TeamID:    x.TeamID,
			FixtureID: x.FixtureID,
			SeasonID:  x.SeasonID,
		}
	}

	_, err = b.
		Insert("team_rating").
		Columns(
			"team_id",
			"fixture_id",
			"season_id",
			"attack_total",
			"attack_points",
			"defence_total",
			"defence_points",
			"timestamp").
		Values(
			x.TeamID,
			x.FixtureID,
			x.SeasonID,
			x.Attack.Total,
			x.Attack.Difference,
			x.Defence.Total,
			x.Defence.Difference,
			x.Timestamp.Unix(),
		).
		Exec()

	return err
}

func queryBuilder(c *sql.DB) sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(c)
}

func NewRatingWriter(c *sql.DB) RatingWriter {
	return &ratingWriter{connection: c}
}
