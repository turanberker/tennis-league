package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// CustomRouterGroup, standart gin.IRouter arayüzünü sarmalar
type CustomRouterGroup struct {
	gin.IRouter
}

// NewCustomGroup standart gin.Engine veya Group'u bizim yapımıza dönüştürür
func NewCustomGroup(router gin.IRouter) *CustomRouterGroup {
	return &CustomRouterGroup{IRouter: router}
}

// wrapHandlers: Middleware'leri öne, bizim sarmalanmış handler'ımızı en arkaya ekler.
func (cg *CustomRouterGroup) wrapHandlers(handlers ...any) []gin.HandlerFunc {

	ginHandlers := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {

		if h, ok := handler.(HandlerFunc); ok {
			ginHandlers[i] = Handler(h)

			// 2. Yol: Eğer yalın metot imzasıysa (h.getIncomingMatches doğrudan buraya düşer!)
		} else if h, ok := handler.(func(*CustomContext)); ok {
			ginHandlers[i] = Handler(h)

			// 3. Yol: Eğer standart Gin middleware'i ise (authmiddleware.RequireAuth vb.)[cite: 3]
		} else if h, ok := handler.(gin.HandlerFunc); ok {
			ginHandlers[i] = h

		} else if h, ok := handler.(func(*gin.Context)); ok {
			ginHandlers[i] = h

		} else {
			// Geliştirme aşamasında buraya düşen hatalı tipleri yakalamak için panik fırlatmak en iyisidir
			panic(fmt.Sprintf("HATA: Rota kaydedilirken geçersiz bir tip algılandı! Tip: %T", handler))
		}

		/*if h, ok := handler.(HandlerFunc); ok {
			ginHandlers[i] = Handler(h)
		} else if h, ok := handler.(gin.HandlerFunc); ok {
			ginHandlers[i] = h
		} else {
			fmt.Println("Bilinmeyen bir tip.")
		}*/

	}

	return ginHandlers
}

// GET, POST, PUT, DELETE metotlarımız artık tamamen tip güvenli.
// İlk sırada her zaman bizim HandlerFunc'ımız olmalı, sonrasında dilediğin kadar Gin middleware'i gelebilir.
func (cg *CustomRouterGroup) GET(relativePath string, handlers ...any) {
	cg.IRouter.GET(relativePath, cg.wrapHandlers(handlers...)...)
}

func (cg *CustomRouterGroup) POST(relativePath string, handlers ...any) {
	cg.IRouter.POST(relativePath, cg.wrapHandlers(handlers...)...)
}

func (cg *CustomRouterGroup) PUT(relativePath string, handlers ...any) {
	cg.IRouter.PUT(relativePath, cg.wrapHandlers(handlers...)...)
}

func (cg *CustomRouterGroup) DELETE(relativePath string, handlers ...any) {
	cg.IRouter.DELETE(relativePath, cg.wrapHandlers(handlers...)...)
}

func (cg *CustomRouterGroup) PATCH(relativePath string, handlers ...any) {
	cg.IRouter.PATCH(relativePath, cg.wrapHandlers(handlers...)...)
}

// Group metodu standart gin.HandlerFunc (middleware'leri) kabul eder.
// Alt gruplar da yine CustomRouterGroup olarak döner.
func (cg *CustomRouterGroup) Group(relativePath string, handlers ...any) *CustomRouterGroup {
	return &CustomRouterGroup{
		IRouter: cg.IRouter.Group(relativePath, cg.wrapHandlers(handlers...)...),
	}
}
