package main

import (
	"log"
	"tennis-league/common/http/router"
	"tennis-league/common/lib/cache"
	"tennis-league/common/lib/database"
	authmiddleware "tennis-league/common/security/authmiddleware"
	"tennis-league/common/security/repository"
	"tennis-league/service/internal/delivery/http/handler/dashboard"
	"tennis-league/service/internal/delivery/http/handler/doubleteamhandler"
	"tennis-league/service/internal/delivery/http/handler/leaguehandler"
	"tennis-league/service/internal/delivery/http/handler/matchhandler"
	"tennis-league/service/internal/domain/league"
	"tennis-league/service/internal/domain/match"
	"tennis-league/service/internal/domain/scoreboard"
	postgres2 "tennis-league/user-service/internal/repository/postgres"
	"tennis-league/user-service/internal/service/player"

	"tennis-league/service/internal/domain/team"
	"tennis-league/service/internal/domain/user"
	"tennis-league/service/internal/infrastructure/persistence/postgres"
	"tennis-league/service/internal/infrastructure/persistence/redis"
)

func main() {
	serverConfig := router.LoadServerConfig()
	matchhandler.RegisterSetValidations()

	db, err := database.NewPostgres()
	if err != nil {
		log.Fatal(err)
	}
	redisClient, err := cache.NewRedis()
	if err != nil {
		log.Fatal(err)
	}
	sessionGetterRepository := repository.NewSessionGetterRepositoryImpl(redisClient)
	sessionRepository := redis.NewSessionRepository(sessionGetterRepository, redisClient)
	transactionManager := database.NewTransactionManager(db)
	userRepo := postgres.NewUserRepository(db)
	userUC := user.NewUsecase(transactionManager, userRepo)
	cacheManager := cache.NewCacheManager(redisClient)

	leagueRepository := postgres.NewLeagueRepository(db)
	teamRepository := postgres.NewTeamRepository(db)
	teamPlayerRepository := postgres.NewTeamPlayerRepository(db)
	matchRepository := postgres.NewMatchRepository(db)
	matchSetRepository := postgres.NewMatchSetRepository(db)
	scoreBoardRepository := postgres.NewScoreBoardRepository(db)
	outboxRepository := postgres.NewOutboxRepository(db)
	playerRepository := postgres2.NewPlayerRepository(db)
	leagueCoordinatorRepository := postgres.NewLeagueCoordinatorRepository(db)

	teamUseCase := team.NewUseCase(transactionManager, cacheManager, teamRepository, teamPlayerRepository)
	matchUseCase := match.NewUseCase(transactionManager, cacheManager, matchRepository, matchSetRepository, outboxRepository)
	leagueUseCase := league.NewUsecase(transactionManager, cacheManager, teamUseCase, matchUseCase, userUC, leagueRepository, teamRepository,
		matchRepository, outboxRepository, scoreBoardRepository, leagueCoordinatorRepository)

	scoreBaordUc := scoreboard.NewUseCase(scoreBoardRepository)
	playerUc := player.NewUsecase(transactionManager, playerRepository, matchRepository)

	dashboardHandler := dashboard.NewDashboardHandler(playerUc)
	leagueHandler := leaguehandler.NewHandler(leagueUseCase, teamUseCase, scoreBaordUc, matchUseCase)

	matchHandler := matchhandler.NewMatchHandler(matchUseCase)
	doubleTeamHandler := doubleteamhandler.NewDoubleTeamHandler(teamUseCase)
	r := router.NewRouter(serverConfig, authmiddleware.NewAuthMiddleware("tennis", sessionRepository),
		dashboardHandler,
		leagueHandler,

		matchHandler,
		doubleTeamHandler)

	log.Println("Server running on :" + serverConfig.Port)
	r.Run(":" + serverConfig.Port)
}
