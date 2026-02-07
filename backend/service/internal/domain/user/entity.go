package user

type Role string

const (
	RoleAdmin       Role = "ADMIN"
	RoleCoordinator Role = "COORDINATOR"
	RolePlayer      Role = "PLAYER"
)

type User struct {
	ID           int64
	Email        string
	Phone        string
	Name         string
	Surname      string
	PasswordHash string
	Role         Role
}

type RegisterUserInput struct {
	Email    string
	Name     string
	Surname  string
	Password string
}
