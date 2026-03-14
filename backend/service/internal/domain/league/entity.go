package league

import (
	"errors"
	"time"
)

type League struct {
	ID                 string
	Name               string
	FixtureCreatedDate *time.Time
}

type PersistLeague struct {
	Name string
}

var LEAGE_WITH_NAME_EXISTS = errors.New("league name already exists")
