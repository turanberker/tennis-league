package auth

import "tennis-league/common/security/dto"

type LoggedInUser struct {
	ID        string
	Name      string
	Surname   string
	Role      dto.Role
	SessionId string
	PlayerId  *string
}

type RegisterUserInput struct {
	Email    string
	Name     string
	Surname  string
	Password string
}
