package player

import "tennis-league/user-interface/constants"

type Player struct {
	ID           string
	Name         string
	Surname      string
	Sex          constants.Sex
	UserId       *string
	DoublePoints int
	SinglePoints int
}

type PersistPlayer struct {
	Name    string
	Surname string
	Sex     constants.Sex
	UserId  *int64
}

type ListQueryParameters struct {
	Name    *string
	Sex     *constants.Sex
	HasUser *bool
}

type PlayerStatisticsRequest struct {
	PlayerId string
	Limit    *int
}

type PlayerStatistics struct {
	CurrentDoublePoint  int
	CurrentSinglePoint  int
	LastDoublePointsSum int
	LastSinglePointsSum int
}
