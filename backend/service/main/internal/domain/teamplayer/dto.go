package teamplayer

type Player struct {
	ID           string
	Name         string
	Surname      string
	Sex          Sex
	UserId       *string
	DoublePoints int
	SinglePoints int
}
