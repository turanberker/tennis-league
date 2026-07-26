package league

type CreateTeamRequestDto struct {
	LeagueId  string
	Name      string
	PlayerIDs []string
}

type CreateTeamResponseDto struct {
	TeamId          string
	TotalAttendance int32
}

type FixtureFilterParam struct {
	TeamId *string
}

type NewLeaguePlayerAttendance struct {
	LeagueId   string
	PlayerId   string
	PlayerName string
}
