package player

type Player struct {
	ID           string
	Name         string
	Surname      string
	Sex          Sex
	UserId       *string
	DoublePoints int
	SinglePoints int
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

type ListQueryParameters struct {
	Name    *string
	Sex     *Sex
	HasUser *bool
}
