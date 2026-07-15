package playerhandler

import (
	"tennis-league/common/http/router"
	"tennis-league/user-interface/constants"
	"tennis-league/user-service/internal/service/player"

	"tennis-league/common/security/authmiddleware"
	"tennis-league/common/security/dto"
)

type PlayerHandler struct {
	uc *player.Usecase
}

func NewPlayerHandler(uc *player.Usecase) *PlayerHandler {
	return &PlayerHandler{uc: uc}
}

func (h *PlayerHandler) RegisterRoutes(r *router.CustomRouterGroup) {

	group := r.Group("/player")
	{
		group.GET("/list", h.getAll)
		group.POST("", authmiddleware.RequireRole(dto.RoleAdmin), h.save)
		group.PUT("/:id/assign-to-user", authmiddleware.RequireRole(dto.RoleAdmin), h.assignToUser)
		group.GET("/unassigned-players", h.unassignedPlayers)
		group.GET("/:id/statistics", h.getPlayerStatistics)
	}

}

func (h *PlayerHandler) save(c *router.CustomContext) {

	var req struct {
		Name    string        `json:"name" binding:"min=3,max=75,required"`
		Surname string        `json:"surname" binding:"min=3,max=75,required"`
		Sex     constants.Sex `json:"sex" binding:"required,oneof=M F"`
	}

	if !c.BindJSONOrAbort(&req) {
		return
	}

	persistPlayer := &player.PersistPlayer{
		Name:    req.Name,
		Surname: req.Surname,
		Sex:     req.Sex,
	}

	playerId, err := h.uc.Save(c.Request.Context(), persistPlayer)

	if err != nil {
		c.ErrorComplete(err)
	} else {
		c.OkComplete(playerId)
	}
}

func (h *PlayerHandler) assignToUser(c *router.CustomContext) {
	playerId := c.Param("id")

	var req struct {
		UserId string `form:"userId" binding:"required"`
	}
	if !c.BindQueryOrAbort(&req) {
		return
	}
	var err error
	ctx := c.Request.Context()

	err = h.uc.AssignToUser(ctx, playerId, req.UserId)

	if err != nil {
		c.ErrorComplete(err)
		return
	}
	c.OkComplete("İşlem başarılı")
}

func (h *PlayerHandler) getAll(c *router.CustomContext) {
	var req struct {
		Name *string        `form:"name" binding:"omitempty"`
		Sex  *constants.Sex `form:"sex" binding:"omitempty,oneof=M F"`
	}

	// Gin otomatik olarak URL'deki ?name=...&sex=... kısımlarını struct'a doldurur
	if !c.BindQueryOrAbort(&req) {
		return
	}

	players, _ := h.uc.List(c.Request.Context(), player.ListQueryParameters{Name: req.Name,
		Sex: req.Sex})

	playersResponse := make([]*PlayerResponse, 0, len(players))

	for _, l := range players {
		playersResponse = append(playersResponse, toPlayerResponse(l))
	}
	c.OkComplete(playersResponse)
}

func (h *PlayerHandler) unassignedPlayers(c *router.CustomContext) {

	var req struct {
		Sex constants.Sex `form:"sex" binding:"oneof=M F"`
	}
	// Gin otomatik olarak URL'deki ?name=...&sex=... kısımlarını struct'a doldurur
	if c.BindQueryOrAbort(&req) {
		return
	}

	isFalse := false
	players, err := h.uc.List(c.Request.Context(),
		player.ListQueryParameters{Sex: &req.Sex,
			HasUser: &isFalse})

	if err != nil {
		c.ErrorComplete(err)
		return
	}

	playerResponseList := make([]*PlayerResponse, 0, len(players))

	for _, l := range players {
		playerResponseList = append(playerResponseList, toPlayerResponse(l))
	}
	c.OkComplete(playerResponseList)
}

func toPlayerResponse(l *player.Player) *PlayerResponse {
	if l == nil {
		return nil
	}

	return &PlayerResponse{
		ID:           l.ID,
		Name:         l.Name,
		Surname:      l.Surname,
		Sex:          l.Sex,
		UserId:       l.UserId,
		DoublePoints: l.DoublePoints,
		SinglePoints: l.SinglePoints,
	}
}

func (h *PlayerHandler) getPlayerStatistics(c *router.CustomContext) {
	playerId := c.Param("id")
	var req struct {
		Limit *int `form:"limit" binding:"omitempty,numeric"`
	}
	if !c.BindQueryOrAbort(&req) {
		return
	}

	statistics, err := h.uc.GetPlayerStatistics(c.Request.Context(), player.PlayerStatisticsRequest{
		PlayerId: playerId,
		Limit:    req.Limit,
	})

	if err != nil {
		c.Abort()
		return
	}

	var response struct {
		EarnedSinglePoints int `json:"earnedSinglePoints"`
		EarnedDoublePoints int `json:"earnedDoublePoints"`
		SinglePoints       int `json:"singlePoints"`
		DoublePoints       int `json:"doublePoints"`
	}
	response.EarnedDoublePoints = statistics.LastDoublePointsSum
	response.EarnedSinglePoints = statistics.LastSinglePointsSum
	response.SinglePoints = statistics.CurrentSinglePoint
	response.DoublePoints = statistics.CurrentDoublePoint

	c.OkComplete(response)
}
