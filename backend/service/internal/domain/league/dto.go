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
