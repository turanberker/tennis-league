package team

type Team struct {
	ID       string
	LeagueID string
	Name     string
}

type PersistTeam struct {
	LeagueID string
	Name     string
}
