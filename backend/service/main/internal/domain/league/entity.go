package league

import (
	"errors"
	"time"
)

type LeagueListSelect struct {
	ID                string
	Name              string
	Format            LEAGUE_FORMAT
	Category          LEAGUE_CATEGORY
	Type              LEAGUE_PROCESS_TYPE
	Status            LEAGUE_STATUS
	TotalAttentance   int32
	CoordinatorUserId []string
}

type League struct {
	ID              string
	Name            string
	Format          LEAGUE_FORMAT
	Category        LEAGUE_CATEGORY
	ProcessType     LEAGUE_PROCESS_TYPE
	Status          LEAGUE_STATUS
	TotalAttendance int32
	StartDate       *time.Time
	EndDate         *time.Time
}

type PersistLeague struct {
	Name        string
	Format      LEAGUE_FORMAT
	Categoty    LEAGUE_CATEGORY
	ProcessType LEAGUE_PROCESS_TYPE
}

type LEAGUE_FORMAT string
type LEAGUE_CATEGORY string
type LEAGUE_PROCESS_TYPE string
type LEAGUE_STATUS string

const (
	LeagueFormat_SINGLE LEAGUE_FORMAT = "SINGLE"
	LeagueFormat_DOUBLE LEAGUE_FORMAT = "DOUBLE"
	LeagueFormat_TEAM   LEAGUE_FORMAT = "TEAM"

	LeagueCategory_MIX   LEAGUE_CATEGORY = "MIX"
	LeagueCategory_ERKEK LEAGUE_CATEGORY = "ERKEK"
	LeagueCategory_KADIN LEAGUE_CATEGORY = "DOUBLE"

	LeagueProcessType_FIXTURE LEAGUE_PROCESS_TYPE = "FIXTURE"
	LeagueProcessType_DEFI    LEAGUE_PROCESS_TYPE = "DEFI"

	LeagueStatus_DRAFT     LEAGUE_STATUS = "DRAFT"
	LeagueStatus_ACTIVE    LEAGUE_STATUS = "ACTIVE"
	LeagueStatus_COMPLETED LEAGUE_STATUS = "COMPLETED"
)

var LEAGE_WITH_NAME_EXISTS = errors.New("league name already exists")

type LeagueMatchApprovedEvent struct {
	LeagueId string `json:"leagueId"`
	MatchId  string `json:"matchId"`
}
