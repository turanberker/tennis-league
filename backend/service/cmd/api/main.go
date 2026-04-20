package main

import (
	"log"

	"github.com/turanberker/tennis-league-service/internal/delivery/http"
	"github.com/turanberker/tennis-league-service/internal/delivery/http/handler/authhandler"
	"github.com/turanberker/tennis-league-service/internal/delivery/http/handler/dashboard"
	"github.com/turanberker/tennis-league-service/internal/delivery/http/handler/leaguehandler"
	"github.com/turanberker/tennis-league-service/internal/delivery/http/handler/matchhandler"
	"github.com/turanberker/tennis-league-service/internal/delivery/http/handler/playerhandler"
	"github.com/turanberker/tennis-league-service/internal/delivery/http/handler/userhandler"
	"github.com/turanberker/tennis-league-service/internal/delivery/http/middleware"
	"github.com/turanberker/tennis-league-service/internal/domain/auth"
	"github.com/turanberker/tennis-league-service/internal/domain/league"
	"github.com/turanberker/tennis-league-service/internal/domain/match"
	"github.com/turanberker/tennis-league-service/internal/domain/player"
	"github.com/turanberker/tennis-league-service/internal/domain/scoreboard"
	"github.com/turanberker/tennis-league-service/internal/platform"

	"github.com/turanberker/tennis-league-service/internal/domain/team"
	"github.com/turanberker/tennis-league-service/internal/domain/user"
	"github.com/turanberker/tennis-league-service/internal/infrastructure/persistence/postgres"
	"github.com/turanberker/tennis-league-service/internal/infrastructure/persistence/redis"

	"github.com/turanberker/tennis-league-service/internal/platform/cache"
	"github.com/turanberker/tennis-league-service/internal/platform/database"
)

func main() {
	serverConfig := platform.LoadServerConfig()
	matchhandler.RegisterSetValidations()

	db, err := database.NewPostgres()
	if err != nil {
		log.Fatal(err)
	}
	redisClient, err := cache.NewRedis()
	if err != nil {
		log.Fatal(err)
	}
	sessionRepository := redis.NewSessionRepository(redisClient)
	transactionManager := database.NewTransactionManager(db)
	userRepo := postgres.NewUserRepository(db)
	userUC := user.NewUsecase(transactionManager, userRepo)
	cacheManager := cache.NewCacheManager(redisClient)

	tokenService := middleware.NewTokenService("tennis")

	leagueRepository := postgres.NewLeagueRepository(db)
	teamRepository := postgres.NewTeamRepository(db)
	teamPlayerRepository := postgres.NewTeamPlayerRepository(db)
	matchRepository := postgres.NewMatchRepository(db)
	matchSetRepository := postgres.NewMatchSetRepository(db)
	scoreBoardRepository := postgres.NewScoreBoardRepository(db)
	outboxRepository := postgres.NewOutboxRepository(db)
	playerRepository := postgres.NewPlayerRepository(db)
	leagueCoordinatorRepository := postgres.NewLeagueCoordinatorRepository(db)

	authUC := auth.NewUsecase(db, userRepo, sessionRepository)
	teamUseCase := team.NewUseCase(transactionManager, teamRepository, teamPlayerRepository)
	matchUseCase := match.NewUseCase(transactionManager, cacheManager, matchRepository, matchSetRepository, outboxRepository)
	leagueUseCase := league.NewUsecase(transactionManager, cacheManager, teamUseCase, matchUseCase, userUC, leagueRepository, teamRepository,
		matchRepository, outboxRepository, scoreBoardRepository, leagueCoordinatorRepository)

	scoreBaordUc := scoreboard.NewUseCase(scoreBoardRepository)
	playerUc := player.NewUsecase(transactionManager, playerRepository, matchRepository)

	dashboardHandler := dashboard.NewDashboardHandler(playerUc)
	leagueHandler := leaguehandler.NewHandler(leagueUseCase, teamUseCase, scoreBaordUc, matchUseCase)
	authHandler := authhandler.NewAuthHandler(authUC, tokenService)
	userHandler := userhandler.NewUserHandler(userUC)
	playerhandler := playerhandler.NewPlayerHandler(playerUc)
	matchHandler := matchhandler.NewMatchHandler(matchUseCase)

	r := http.NewRouter(serverConfig,
		middleware.NewAuthMiddleware("tennis", sessionRepository),
		dashboardHandler,
		authHandler,
		leagueHandler,
		playerhandler,
		matchHandler,
		userHandler)

	log.Println("Server running on :" + serverConfig.Port)
	r.Run(":" + serverConfig.Port)
}
