package playerhandler

import "github.com/turanberker/tennis-league-service/internal/domain/player"

type PlayerResponse struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Sex          player.Sex `json:"sex"`
	Surname      string     `json:"surname"`
	UserId       *string    `json:"user_id"`
	DoublePoints int        `json:"double_points"`
	SinglePoints int        `json:"single_points"`
}
