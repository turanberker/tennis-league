package leaguehandler

import "time"

type LeagueResponse struct {
	ID                 string     `json:"id"`
	Name               string     `json:"name"`
	FixtureCreatedDate *time.Time `json:"fixtureCreatedDate,omitempty"`
}
