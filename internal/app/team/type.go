package team

import "time"

type Rating struct {
	TeamID    uint64
	FixtureID uint64
	SeasonID  uint64
	Attack    Points
	Defence   Points
	Timestamp time.Time
}

type Points struct {
	Total      uint8
	Difference uint8
}
