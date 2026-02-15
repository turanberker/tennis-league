package league

import "time"

type League struct {
	ID                 string
	Name               string
	FixtureCreatedDate *time.Time
}

type PersistLeague struct {
	Name string
}
