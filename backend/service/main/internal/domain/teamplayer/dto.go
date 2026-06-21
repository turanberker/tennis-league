package teamplayer

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
