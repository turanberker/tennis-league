package userhandler

import "github.com/turanberker/tennis-league-service/internal/domain/user"

type UserResponse struct {
	Id       string    `json:"id"`
	Name     string    `json:"name"`
	Surname  string    `json:"surname"`
	Email    string    `json:"email"`
	Approved bool      `json:"approved"`
	Role     user.Role `json:"role"`
	PlayerId *string   `json:"playerId,omitempty"`
}
