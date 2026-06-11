package dto

type Role string

const (
	RoleAdmin       Role = "ADMIN"
	RoleCoordinator Role = "COORDINATOR"
	RolePlayer      Role = "PLAYER"
)

type Session struct {
	SessionId string
	UserId    string
	Role      string
	PlayerId  *string
}
