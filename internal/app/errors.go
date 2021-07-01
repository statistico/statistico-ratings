package app

import "fmt"

type DuplicationError struct {
	TeamID  uint64
	FixtureID  uint64
	SeasonID  uint64
}

func (d *DuplicationError) Error() string {
	return fmt.Sprintf(
		"team rating exists for team %d, fixture %d and season %d",
		d.TeamID,
		d.FixtureID,
		d.SeasonID,
	)
}

type NotFoundError struct {
	TeamID uint64
}

func (n *NotFoundError) Error() string {
	return fmt.Sprintf("team %d rating does not exist", n.TeamID)
}
