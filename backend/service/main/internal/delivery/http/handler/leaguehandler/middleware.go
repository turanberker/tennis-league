package leaguehandler

import (
	"fmt"
	"net/http"
	"tennis-league/common/http/router"
	customerror "tennis-league/common/lib/error"
	"tennis-league/common/security/authmiddleware"
	"tennis-league/common/security/dto"
	errorcodes "tennis-league/service/internal/domain/error_codes"
	"tennis-league/service/internal/domain/league"
	"tennis-league/service/internal/domain/match"
)

type leagueHandlerMiddleware struct {
	uc      *league.Usecase
	matchUc *match.UseCase
}

func (h *leagueHandlerMiddleware) checkIfCoordinator(c *router.CustomContext) {
	roleValue, _ := c.Get("Role")
	leagueId := c.Param("id")
	userId, _ := c.CurrentUserId()

	if role, ok := roleValue.(dto.Role); ok {

		// 3. Karşılaştırma yap
		if role == dto.RoleCoordinator {
			coordinator, err := h.uc.IsUserCoordinator(c.Request.Context(), leagueId, userId)
			if err != nil {
				c.ErrorComplete(err)
				return
			}
			if coordinator {
				c.Next()
			} else {
				err := &customerror.BusinnesException{
					StatusCode: http.StatusForbidden,
					ErrorCode:  authmiddleware.INSUFFICIENT_PERMISSIONS,
					Message:    "Bu ligde koordinatör değilsiniz",
				}
				_ = c.Error(err)
				c.Abort()
			}
		}

		if role == dto.RoleAdmin {
			c.Next()
		}
	} else {
		err := &customerror.BusinnesException{
			StatusCode: http.StatusForbidden,
			ErrorCode:  authmiddleware.INSUFFICIENT_PERMISSIONS,
			Message:    "Bu ligde yetkiniz yok",
		}
		_ = c.Error(err)
		c.Abort()
	}
}

func (h *leagueHandlerMiddleware) checkIfUserIsCoordinateAdminOrPlayer(c *router.CustomContext) {
	roleValue, _ := c.Get("Role")
	userId, _ := c.CurrentUserId()
	matchId := c.Param("matchId")
	// Not: Lig ID'si bu context'te farklı bir isimle (örn: leagueId) geliyorsa onu almalısın.
	// Eğer match üzerinden leagueId'ye gitmek gerekiyorsa usecase katmanında bu kontrolü yapabilirsin.
	leagueId := c.Param("id")

	role, ok := roleValue.(dto.Role)
	if !ok {
		h.abortWithForbidden(c, "Yetki bilgisi alınamadı")
		return
	}
	// 1. Durum: Admin ise sınırsız erişim
	if role == dto.RoleAdmin {
		c.Next()
		return
	}

	// 2. Durum: Koordinatör ise lig bazlı kontrol
	if role == dto.RoleCoordinator {
		coordinator, err := h.uc.IsUserCoordinator(c.Request.Context(), leagueId, userId)
		if err == nil && coordinator {
			c.Next()
			return
		}
		// Eğer koordinatör değilse hemen abort etmiyoruz, belki bu maçın oyuncusudur.
	}

	// 3. Durum: Oyuncu mu kontrolü (Admin veya Lig Koordinatörü değilse buraya düşer)
	playedInMatch := h.isPlayerPlayedInMatch(c, matchId)
	if playedInMatch {
		c.Next()
		return
	}

	// Hiçbir şart sağlanmadıysa erişimi reddet
	h.abortWithForbidden(c, "Bu işlem için yetkiniz bulunmamaktadır (Koordinatör, Admin veya Maçın Oyuncusu olmalısınız)")
}

// Yardımcı metod: Kod tekrarını önlemek için
func (h *leagueHandlerMiddleware) abortWithForbidden(c *router.CustomContext, message string) {
	err := &customerror.BusinnesException{
		StatusCode: http.StatusForbidden,
		ErrorCode:  authmiddleware.INSUFFICIENT_PERMISSIONS,
		Message:    message,
	}
	_ = c.Error(err)
	c.Abort()
}

func (h *leagueHandlerMiddleware) isPlayerPlayedInMatch(c *router.CustomContext, matchId string) bool {
	playerId, exists := c.CurrentPlayerId()
	if !exists {
		return false
	}
	playedInMatch, _ := h.matchUc.IsUserPlayerOfMatch(c.Request.Context(), matchId, playerId)
	return playedInMatch
}

func (h *leagueHandlerMiddleware) checkIfMatchIsLeague(c *router.CustomContext) {
	matchId := c.Param("matchId")
	matchInfo, err := h.matchUc.GetMatchInfo(c.Request.Context(), matchId)
	if err != nil {
		c.ErrorComplete(err)
		return
	}
	if matchInfo.Source != match.MatchSource_LEAGUE {
		businessErr := &customerror.BusinnesException{
			StatusCode: http.StatusBadRequest,
			ErrorCode:  errorcodes.ErrorINVALID_MATCH_SOURCE,
			Message:    "Buradan sadece Lig maçları güncellenebilir",
		}
		c.ErrorComplete(businessErr)

		return
	}
	c.Next()
}

func (h *leagueHandlerMiddleware) checkLeagueIsChallenging(c *router.CustomContext) {
	leagueId := c.Param("id")

	leagueData, err := h.uc.GetById(c.Request.Context(), leagueId)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	if leagueData.ProcessType == league.LeagueProcessType_DEFI {
		c.Next()
	} else {
		businessErr := &customerror.BusinnesException{
			StatusCode: http.StatusBadRequest,
			ErrorCode:  errorcodes.ErrorInvalid_ProcessType,
			Message:    fmt.Sprintf("%s tipinde ligler için maç talebinde bulunabilirsiniz!", league.LeagueProcessType_DEFI),
		}

		c.ErrorComplete(businessErr)
		return
	}

}

func (h *leagueHandlerMiddleware) userAttendedToLeague(c *router.CustomContext) {
	leagueId := c.Param("id")
	playerId, _ := c.CurrentPlayerId()
	attended, err := h.uc.IsPlayerAttandedToLeague(c.Request.Context(), playerId, leagueId)
	if err != nil {
		c.ErrorComplete(err)
		return
	}
	if !attended {
		h.abortWithForbidden(c, "Bu lige katılmadığınız için maç talebinde bulunamazsınız")
		return
	}

}
