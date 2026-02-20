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

			diff := abs(t1 - t2)

			// biri en az 6 olmalı
			if t1 < 6 && t2 < 6 {
				return false
			}

			// en az 2 fark
			if diff < 2 {
				return false
			}

			// max 7 olabilir
			if t1 > 7 || t2 > 7 {
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

			if t1 < 10 && t2 < 10 {
				return false
			}

			if diff < 2 {
				return false
			}

			return true
		})
	}
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
