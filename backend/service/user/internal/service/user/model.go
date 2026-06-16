package user

import (
	"errors"
	"tennis-league/common/security/dto"
	"time"
)

type LoginUserCheck struct {
	ID           string
	Email        string
	Name         string
	Surname      string
	PasswordHash string
	Role         dto.Role
	PlayerId     *string
}

type PersistUser struct {
	Email        string
	Name         string
	Surname      string
	PasswordHash string
	Role         dto.Role
}
type UserRepository struct {
}

type UserData struct {
	ID           string
	Email        string
	Phone        *string
	Name         string
	Surname      string
	Role         dto.Role
	PasswordHash string
	CreatedAt    time.Time
	Approved     bool
}

type User struct {
	Id       string
	Name     string
	Surname  string
	Email    string
	Approved bool
	Role     dto.Role
	PlayerId *string
}

var USER_EXISTS_ERROR = errors.New("Bu Kullanıcı Mevcut")
