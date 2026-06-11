package httpcache

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func AddCacheControlHeader(maxAge int, cacheType Type) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Dinamik olarak header değerini oluşturuyoruz
		headerValue := fmt.Sprintf("%s, max-age=%d", cacheType, maxAge)

		c.Header("Cache-Control", headerValue)

		// Bir sonraki handler'a geçişi sağlar
		c.Next()
	}
}

type Type string

const (
	TYPE_PUBLIC  Type = "public"
	TYPE_PRIVATE Type = "private"
)
