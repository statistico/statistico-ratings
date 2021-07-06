-- +goose Up
-- +goose StatementBegin
CREATE TABLE team_rating (
  id SERIAL,
  team_id INTEGER NOT NULL,
  fixture_id INTEGER NOT NULL,
  season_id INTEGER NOT NULL,
  attack_total DECIMAL NOT NULL,
  attack_points DECIMAL NOT NULL,
  defence_total DECIMAL NOT NULL,
  defence_points DECIMAL NOT NULL,
  timestamp INTEGER NOT NULL
);

CREATE INDEX ON team_rating (team_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE team_rating;
-- +goose StatementEnd
