package constants

type Sex string

func (s Sex) IsValid() bool {
	return s == SexMale || s == SexFemale
}

const (
	SexMale   Sex = "M"
	SexFemale Sex = "F"
)
