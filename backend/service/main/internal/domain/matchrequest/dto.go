package matchrequest

import (
	customtype "tennis-league/common/lib/type"
)

type Type string

var (
	TypeSingle Type = "SINGLE"
	TypeDouble Type = "DOUBLE"
	TypeTeam   Type = "TEAM"
)

type RequestDto struct {
	LeagueId  *string
	PlayerId  string
	Date      customtype.Date
	StartHour int
	Type      Type
}
