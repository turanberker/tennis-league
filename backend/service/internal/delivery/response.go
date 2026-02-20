package delivery

type Response struct {
	Success bool              `json:"success"`
	Data    interface{}       `json:"data,omitempty"`
	Error   string            `json:"error,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
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
		Error:   message,
	}
}

func NewValidationErrorResponse(errors map[string]string) *Response {
	return &Response{
		Success: false,
		Errors:  errors,
	}
}
