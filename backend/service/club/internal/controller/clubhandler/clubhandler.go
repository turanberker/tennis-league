package clubhandler

import (
	"tennis-league/common/http/router"
	"tennis-league/common/security/authmiddleware"
	"tennis-league/common/security/dto"

	"github.com/gin-gonic/gin"
)

type ClubHandler struct {
}

func NewClubHandler() *ClubHandler {
	return &ClubHandler{}
}

func (h *ClubHandler) RegisterRoutes(r *gin.Engine) {

	group := r.Group("/club")
	{
		group.GET("/list", router.Wrap(h.list))
		group.POST("", h.save, authmiddleware.RequireRole(dto.RoleAdmin))
	}
}

func (h *ClubHandler) list(c *router.CustomContext) {}
func (h *ClubHandler) save(c *gin.Context) {

	req, ok := router.RequestBody[PersistClub](c)

}
