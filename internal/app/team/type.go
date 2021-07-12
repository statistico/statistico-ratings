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
	Total      float64
	Difference float64
}

type ReaderQuery struct {
	TeamID  *uint64
	SeasonID  *uint64
	Before   *time.Time
	Sort     string
}
