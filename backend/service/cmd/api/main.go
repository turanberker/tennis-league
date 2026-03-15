package main

import (
	"log"

	"github.com/turanberker/tennis-league-service/internal/delivery/http"
	"github.com/turanberker/tennis-league-service/internal/delivery/http/handler/authhandler"
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

	"github.com/turanberker/tennis-league-service/internal/domain/team"
	"github.com/turanberker/tennis-league-service/internal/domain/user"
	"github.com/turanberker/tennis-league-service/internal/infrastructure/persistence/postgres"
	"github.com/turanberker/tennis-league-service/internal/infrastructure/persistence/redis"

	"github.com/turanberker/tennis-league-service/internal/platform/cache"
	"github.com/turanberker/tennis-league-service/internal/platform/database"
)

func main() {

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
	userUC := user.NewUsecase(userRepo, transactionManager)

	authUC := auth.NewUsecase(db, userRepo, sessionRepository)
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

	leagueUseCase := league.NewUsecase(db, transactionManager, leagueRepository, teamRepository,
		matchRepository, scoreBoardRepository, leagueCoordinatorRepository, userUC)
	teamUseCase := team.NewUseCase(db, teamRepository, teamPlayerRepository)
	matchUseCase := match.NewUseCase(db, matchRepository, matchSetRepository, outboxRepository)
	scoreBaordUc := scoreboard.NewUseCase(scoreBoardRepository)
	leagueHandler := leaguehandler.NewHandler(leagueUseCase, teamUseCase, scoreBaordUc)
	authHandler := authhandler.NewAuthHandler(authUC, tokenService)
	userHandler := userhandler.NewUserHandler(transactionManager, userUC)

	playerUc := player.NewUsecase(db, playerRepository)
	playerhandler := playerhandler.NewPlayerHandler(playerUc, transactionManager)
	matchHandler := matchhandler.NewMatchHandler(matchUseCase)
	r := http.NewRouter(middleware.NewAuthMiddleware("tennis", sessionRepository),
		authHandler,
		leagueHandler,
		playerhandler,
		matchHandler,
		userHandler)

	log.Println("Server running on :8500")
	r.Run(":8500")
}
