package customerror

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type BusinnesException struct {
	StatusCode int //
	ErrorCode  string
	Message    string // Kullanıcıya gidecek mesaj
}

func (e *BusinnesException) Error() string {
	return fmt.Sprintf("[%s] %s", e.ErrorCode, e.Message)
}

func NewBussinnessError(statusCode int, errorCode string, message string) *BusinnesException {
	return &BusinnesException{
		StatusCode: statusCode, ErrorCode: errorCode, Message: message,
	}
}

func NewValidationError(err validator.ValidationErrors) *BusinnesException {
	var details []string
	for _, f := range err {
		details = append(details, fmt.Sprintf("%s alanı %s kuralına uymuyor", f.Field(), f.Tag()))
	}

	return &BusinnesException{
		StatusCode: http.StatusBadRequest,
		ErrorCode:  CONSTRAINT_VALIDATION, Message: strings.Join(details, ", "),
	}
}

type InternalException struct {
	RawError error // Loglanacak asıl teknik hata (SQL hatası vb.)
}

func (e *InternalException) Error() string {
	return e.Error()
}

func NewInternalError(rawError error) *InternalException {
	return &InternalException{
		RawError: rawError,
	}
}
