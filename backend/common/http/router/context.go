package router

import (
	customerror "tennis-league/common/lib/error"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CustomContext struct {
	*gin.Context
}

func RequestBody[T any](c *gin.Context) (*T, bool) {
	var req T

	if err := c.ShouldBindJSON(&req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			c.Error(customerror.NewValidationError(ve))
		} else {
			c.Error(customerror.NewInternalError(err))
		}
		c.Abort()
		return nil, false
	}

	return &req, true
}

// 3. Handler fonksiyonlarımızın Gin yerine CustomContext alabilmesi için bir tip tanımlıyoruz
type HandlerFunc func(*CustomContext)

func Wrap(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(&CustomContext{Context: c})
	}
}
