package delivery

import customerror "github.com/turanberker/tennis-league-service/internal/domain/error"

type Response struct {
	Success bool              `json:"success"`
	Data    interface{}       `json:"data,omitempty"`
	Error   ErrorDetail       `json:"error,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

type ErrorDetail struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewSuccessResponse(data interface{}) *Response {
	return &Response{
		Success: true,
		Data:    &data,
	}
}

func NewErrorResponse(message string) *Response {
	return &Response{
		Success: false,
		Error:   ErrorDetail{Code: "N/A", Message: message},
	}
}

func NewBusinnesErrorResponse(err customerror.BusinnesException) *Response {
	return &Response{
		Success: false,
		Error:   ErrorDetail{Code: err.ErrorCode, Message: err.Message},
	}
}

var (
	UnexpectedError = &Response{
		Success: false,
		Error:   ErrorDetail{Code: "INT-001", Message: "Beklenmedik bir hata oluştu"},
	}
)

func NewValidationErrorResponse(errors map[string]string) *Response {
	return &Response{
		Success: false,
		Errors:  errors,
	}
}
