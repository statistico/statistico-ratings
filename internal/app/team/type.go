package team

import "time"

type Rating struct {
	TeamID    uint8
	FixtureID int8
	Attack    Points
	Defence   Points
	Timestamp time.Time
}

type Points struct {
	Total      uint64
	Difference uint64
}
