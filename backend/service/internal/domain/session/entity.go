package session

type Session struct {
	SessionId string
	UserId    string
	Role      string
	PlayerId  *string
}

type StartSessionInput struct {
	UserId   string
	Role     string
	PlayerId *string
}
