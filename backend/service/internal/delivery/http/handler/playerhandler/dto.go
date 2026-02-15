package playerhandler

import "github.com/turanberker/tennis-league-service/internal/domain/player"

type PlayerResponse struct {
	ID      string     `json:"id"`
	Name    string     `json:"name"`
	Sex     player.Sex `json:"sex"`
	Surname string     `json:"surname"`
	UserId  *int64     `json:"user_id"`
}
