package leaguehandler

import (
	"net/http"
	customerror "tennis-league/common/lib/error"
	"tennis-league/common/security/authmiddleware"
	"tennis-league/common/security/dto"
	errorcodes "tennis-league/service/internal/domain/error_codes"
	"tennis-league/service/internal/domain/league"
	"tennis-league/service/internal/domain/match"

	"github.com/gin-gonic/gin"
)

type leagueHandlerMiddleware struct {
	uc      *league.Usecase
	matchUc *match.UseCase
}

func (h *leagueHandlerMiddleware) checkIfCoordinator(c *gin.Context) {
	roleValue, _ := c.Get("Role")
	leagueId := c.Param("id")
	userId, _ := authmiddleware.GetUserIdFromContext(c)

	if role, ok := roleValue.(dto.Role); ok {

		// 3. Karşılaştırma yap
		if role == dto.RoleCoordinator {
			coordinator, err := h.uc.IsUserCoordinator(c.Request.Context(), leagueId, userId)
			if err != nil {
				_ = c.Error(customerror.NewInternalError(err))
				c.Abort()
			}
			if coordinator {
				c.Next()
			} else {
				err := &customerror.BusinnesException{
					StatusCode: http.StatusForbidden,
					ErrorCode:  errorcodes.INSUFFICIENT_PERMISSIONS,
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
			ErrorCode:  errorcodes.INSUFFICIENT_PERMISSIONS,
			Message:    "Bu ligde yetkiniz yok",
		}
		_ = c.Error(err)
		c.Abort()
	}
}

func (h *leagueHandlerMiddleware) checkIfUserIsCoordinatAdminOrPlayer(c *gin.Context) {
	roleValue, _ := c.Get("Role")
	userId, _ := authmiddleware.GetUserIdFromContext(c)
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
func (h *leagueHandlerMiddleware) abortWithForbidden(c *gin.Context, message string) {
	err := &customerror.BusinnesException{
		StatusCode: http.StatusForbidden,
		ErrorCode:  errorcodes.INSUFFICIENT_PERMISSIONS,
		Message:    message,
	}
	_ = c.Error(err)
	c.Abort()
}

func (h *leagueHandlerMiddleware) isPlayerPlayedInMatch(c *gin.Context, matchId string) bool {
	playerId, exists := authmiddleware.GetPlayerIdFromContext(c)
	if !exists {
		return false
	}
	playedInMatch, _ := h.matchUc.IsUserPlayerOfMatch(c.Request.Context(), matchId, playerId)
	return playedInMatch
}

func (h *leagueHandlerMiddleware) checkIfMatchIsLeague(c *gin.Context) {
	matchId := c.Param("matchId")
	matchInfo, err := h.matchUc.GetMatchInfo(c.Request.Context(), matchId)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}
	if matchInfo.Source != match.MatchSource_LEAGUE {
		businessErr := &customerror.BusinnesException{
			StatusCode: http.StatusBadRequest,
			ErrorCode:  errorcodes.ErrorINVALID_MATCH_SOURCE,
			Message:    "Buradan sadece Lig maçları güncellenebilir",
		}

		_ = c.Error(businessErr)
		c.Abort()
		return
	}
	c.Next()
}
