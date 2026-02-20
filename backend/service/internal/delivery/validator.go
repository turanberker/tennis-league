package delivery

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidationError(err error) map[string]string {
	errors := make(map[string]string)

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			field := strings.ToLower(fe.Field())
			errors[field] = validationMessage(fe)
		}
	}

	return errors
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "Bu alan zorunludur"
	case "gte":
		return "Değer çok küçük"
	case "lte":
		return "Değer çok büyük"
	default:
		return "Geçersiz değer"
	}
}