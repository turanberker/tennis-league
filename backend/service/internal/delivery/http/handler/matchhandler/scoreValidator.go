package matchhandler

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func RegisterSetValidations() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		// Normal Set
		v.RegisterValidation("tennis_set", func(fl validator.FieldLevel) bool {
			score := fl.Field().Interface().(SetScore)

			t1 := score.Team1Score
			t2 := score.Team2Score

			// biri en az 6 olmalı
			if t1 == t2 {
				return false
			}

			return true
		})

		// Super Tie
		v.RegisterValidation("super_tie", func(fl validator.FieldLevel) bool {
			score := fl.Field().Interface().(SetScore)

			t1 := score.Team1Score
			t2 := score.Team2Score

			diff := abs(t1 - t2)

			if diff < 2 {
				return false
			}

			return true
		})

		v.RegisterStructValidation(MatchScoreStructValidator, MatchScore{})
	}

}

// Fonksiyonun kendisi
func MatchScoreStructValidator(sl validator.StructLevel) {
	match := sl.Current().Interface().(MatchScore)

	// Set 1 kazananı
	s1Winner := match.Set1.Team1Score > match.Set1.Team2Score
	// Set 2 kazananı
	s2Winner := match.Set2.Team1Score > match.Set2.Team2Score

	isTie := s1Winner != s2Winner

	// Eğer setler 1-1 ise (kazananlar farklıysa)
	if isTie {
		if match.SuperTie == nil {
			sl.ReportError(match.SuperTie, "SuperTie", "superTie", "supertie_required", "")
		}
	} else {
		if match.SuperTie != nil {
			sl.ReportError(match.SuperTie, "SuperTie", "superTie", "supertie_must_be_null", "")
		}
	}
}

func abs(x int8) int8 {
	if x < 0 {
		return -x
	}
	return x
}
