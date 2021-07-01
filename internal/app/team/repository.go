package team

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/statistico/statistico-ratings/internal/app"
)

type RatingRepository interface {
	Insert(r *Rating) error
}

type ratingRepository struct {
	connection *sql.DB
}

func (r *ratingRepository) Insert(x *Rating) error {
	var exists bool

	b := queryBuilder(r.connection)

	err := b.Select(`SELECT exists (SELECT id FROM team_rating where name = $1 and user_id = $2)`).
		Where(sq.Eq{"team_id": x.TeamID}).
		Where(sq.Eq{"season_id": x.SeasonID}).
		Where(sq.Eq{"fixture_id": x.FixtureID}).
		Scan(exists)

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
		Values().
		Exec()

	return err
}

func queryBuilder(c *sql.DB) sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(c)
}

func NewRatingRepository(c *sql.DB) RatingRepository {
	return &ratingRepository{connection: c}
}
