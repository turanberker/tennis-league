package player

type Player struct {
	ID      int64
	Uuid    string
	Name    string
	Surname string
	UserId  *int64
}

type PersistPlayer struct {
	Name    string
	Surname string
	UserId  *int64
}
