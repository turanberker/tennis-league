package team

type Team struct {
	ID       string
	LeagueID string
	Name     string
}

type LeagueTeam struct {
	ID    string
	Name  string
	Power int32
}

type PersistTeam struct {
	LeagueID string
	Name     string
}
