package userhandler

import "tennis-league/common/security/dto"

type UserResponse struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Surname  string   `json:"surname"`
	Email    string   `json:"email"`
	Approved bool     `json:"approved"`
	Role     dto.Role `json:"role"`
	PlayerId *string  `json:"playerId,omitempty"`
}
