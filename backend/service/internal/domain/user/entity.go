package user

import "errors"

type Role string

const (
	RoleAdmin       Role = "ADMIN"
	RoleCoordinator Role = "COORDINATOR"
	RolePlayer      Role = "PLAYER"
)

type LoginUserCheck struct {
	ID           string
	Email        string
	Name         string
	Surname      string
	PasswordHash string
	Role         Role
	PlayerId     *string
}

type PersistUser struct {
	Email        string
	Name         string
	Surname      string
	PasswordHash string
	Role         Role
}

type LoggedInUser struct {
	ID        string
	Name      string
	Surname   string
	Role      Role
	SessionId string
	PlayerId  *string
}

type RegisterUserInput struct {
	Email    string
	Name     string
	Surname  string
	Password string
}

type User struct {
	Id       string
	Name     string
	Surname  string
	Email    string
	Approved bool
	Role     Role
	PlayerId *string
}

var USER_EXISTS_ERROR = errors.New("Bu Kullanıcı Mevcut")
