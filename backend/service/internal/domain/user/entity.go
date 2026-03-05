package user

type Role string

const (
	RoleAdmin       Role = "ADMIN"
	RoleCoordinator Role = "COORDINATOR"
	RolePlayer      Role = "PLAYER"
)

type User struct {
	ID           string
	Email        string
	Phone        string
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
