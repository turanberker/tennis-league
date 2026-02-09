package delivery

type Response struct {
	Success      bool        `json:"success"`
	ErrorMessage string      `json:"errorDetail,omitempty"`
	Data         interface{} `json:"data,omitempty"`
}

func NewSuccessResponse(data interface{}) *Response {
	return &Response{
		Success: true,
		Data:    data,
	}
}

func NewErrorResponse(errorMessage string) *Response {
	return &Response{
		Success:      false,
		ErrorMessage: errorMessage,
	}
}
