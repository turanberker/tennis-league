package doubleteamhandler

import "tennis-league/user-interface/constants"

type PlayerResponse struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Sex          constants.Sex `json:"sex"`
	Surname      string        `json:"surname"`
	UserId       *string       `json:"userId"`
	DoublePoints int           `json:"doublePoints"`
	SinglePoints int           `json:"singlePoints"`
}
