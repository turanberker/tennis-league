package player

type Player struct {
	ID      int64
	Uuid    string
	Name    string
	Surname string
	Sex     Sex
	UserId  *int64
}

type PersistPlayer struct {
	Name    string
	Surname string
	Sex     Sex
	UserId  *int64
}

type Sex string

func (s Sex) IsValid() bool {
	return s == SexMale || s == SexFemale
}

const (
	SexMale   Sex = "M"
	SexFemale Sex = "F"
)
