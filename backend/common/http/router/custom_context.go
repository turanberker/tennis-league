package router

import (
	"net/http"
	"tennis-league/common/lib/http/delivery"

	"github.com/gin-gonic/gin"
)

type CustomContext struct {
	*gin.Context
}

// HandlerFunc projedeki tek geçerli handler tipidir
type HandlerFunc func(*CustomContext)

// Handler: Gin handler'larını bizim CustomContext'imize dönüştürür
func Handler(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(&CustomContext{Context: c})
	}
}

// BindQueryOrAbort query parametrelerini bind eder, hata varsa doğrudan BadRequest döner
func (c *CustomContext) BindQueryOrAbort(obj any) bool {
	if err := c.ShouldBindQuery(obj); err != nil {
		errorMessage := delivery.ValidationError(err)
		c.JSON(http.StatusBadRequest, delivery.NewValidationErrorResponse(errorMessage))
		c.Abort()
		return false
	}
	return true
}

// BindJSONOrAbort JSON body'sini bind eder, hata varsa BadRequest döner
func (c *CustomContext) BindJSONOrAbort(obj any) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		errorMessage := delivery.ValidationError(err)
		c.JSON(http.StatusBadRequest, delivery.NewValidationErrorResponse(errorMessage))
		c.Abort()
		return false
	}
	return true
}

func (c *CustomContext) OkComplete(data any) {
	res := delivery.NewSuccessResponse(data)
	c.JSON(http.StatusOK, res)
}

func (c *CustomContext) CurrentPlayerId() (string, bool) {
	playerIdIdValue, exists := c.Get("PlayerId")
	return playerIdIdValue.(string), exists
}

func (c *CustomContext) CurrentUserId() (string, bool) {
	userIdValue, exists := c.Get("UserId")
	return userIdValue.(string), exists
}
