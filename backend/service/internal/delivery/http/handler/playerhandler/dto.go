package playerhandler

import "github.com/turanberker/tennis-league-service/internal/domain/player"

type PlayerResponse struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Sex          player.Sex `json:"sex"`
	Surname      string     `json:"surname"`
	UserId       *string    `json:"userId"`
	DoublePoints int        `json:"doublePoints"`
	SinglePoints int        `json:"singlePoints"`
}
