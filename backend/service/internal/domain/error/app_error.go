package customerror

import "fmt"

type BusinnesException struct {
	StatusCode int // 4
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
