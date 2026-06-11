package handler

import (
	"errors"
	"log"
	"net/http"
	customerror "tennis-league/common/lib/error"
	"tennis-league/common/lib/http/delivery"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Önce handler çalışsın

		// Eğer bir hata varsa
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// 1. Business Exception Kontrolü
			var businessErr *customerror.BusinnesException
			if errors.As(err, &businessErr) {
				c.JSON(businessErr.StatusCode, delivery.NewBusinnesErrorResponse(*businessErr))
				return // Yanıt gönderildikten sonra durmalı
			}

			// 2. Internal Exception Kontrolü
			var internalErr *customerror.InternalException
			if errors.As(err, &internalErr) {
				// appErr.Message kontrolü yerine direkt RawError kontrolü daha mantıklı olabilir
				if internalErr.RawError != nil {
					log.Printf("[SİSTEM HATASI]: %v", internalErr.RawError)
				}

				// Kullanıcıya teknik detay değil, InternalException içindeki güvenli mesajı dönüyoruz
				c.JSON(http.StatusInternalServerError, delivery.UnexpectedError)
				return
			}

			// 3. Beklenmedik Hatalar
			log.Printf("[KRİTİK HATA]: %v", err)
			c.JSON(http.StatusInternalServerError, delivery.UnexpectedError)
		}
	}
}
