package session

type Session struct {
	SessionId string
	UserId    string
	Role      string
	PlayerId  *string
}
