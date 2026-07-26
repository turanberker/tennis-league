package main

import (
	"log"
	"tennis-league/common/http/router"
	"tennis-league/common/lib/cache"
	"tennis-league/common/lib/database"
	"tennis-league/common/security/authmiddleware"
	"tennis-league/common/security/repository"
	"tennis-league/service/internal/delivery/http/handler/dashboard"
	"tennis-league/service/internal/delivery/http/handler/doubleteamhandler"
	"tennis-league/service/internal/delivery/http/handler/leaguehandler"
	"tennis-league/service/internal/delivery/http/handler/matchhandler"
	"tennis-league/service/internal/domain/league"
	"tennis-league/service/internal/domain/match"
	"tennis-league/service/internal/domain/matchrequest"
	"tennis-league/service/internal/domain/scoreboard"
	"tennis-league/service/internal/domain/team"
	"tennis-league/service/internal/domain/user"
	"tennis-league/service/internal/infrastructure/persistence/postgres"
	"tennis-league/service/internal/infrastructure/persistence/postgres/fixturerepository"
	"tennis-league/service/internal/infrastructure/persistence/redis"
	"tennis-league/user-interface/grpc/pb/playerpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	attendanceRepository := postgres.NewAttendanceRepository(db)
	teamPlayerRepository := postgres.NewTeamPlayerRepository(db)
	matchRepository := postgres.NewMatchRepository(db)
	matchSetRepository := postgres.NewMatchSetRepository(db)
	scoreBoardRepository := fixturerepository.NewScoreBoardRepository(db)
	outboxRepository := postgres.NewOutboxRepository(db)
	matchRequestRepository := postgres.NewMatchRequestRepository(db)
	leagueCoordinatorRepository := postgres.NewLeagueCoordinatorRepository(db)

	userGrpcAddr := "localhost:50051"

	conn, err := grpc.NewClient(userGrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("User Service gRPC bağlantı hatası: %v", err)
	}
	// Bağlantıyı uygulama kapanırken kapatmayı unutmuyoruz
	defer conn.Close()

	// 2. Üretilen paketi kullanarak bir Player Service istemcisi (client) oluşturuyoruz
	playerClient := playerpb.NewPlayerServiceClient(conn)

	teamUseCase := team.NewUseCase(transactionManager, cacheManager, attendanceRepository, teamPlayerRepository)
	matchUseCase := match.NewUseCase(transactionManager, cacheManager, matchRepository, matchSetRepository, outboxRepository)
	leagueUseCase := league.NewUsecase(transactionManager, cacheManager, teamUseCase, matchUseCase, userUC, leagueRepository, attendanceRepository,
		matchRepository, outboxRepository, leagueCoordinatorRepository, attendanceRepository, playerClient)
	matchRequestUseCase := matchrequest.NewMatchRequestUseCase(transactionManager, matchRequestRepository)
	scoreBaordUc := scoreboard.NewUseCase(scoreBoardRepository)

	dashboardHandler := dashboard.NewDashboardHandler(matchUseCase)
	leagueHandler := leaguehandler.NewHandler(leagueUseCase, teamUseCase, scoreBaordUc, matchUseCase, matchRequestUseCase, playerClient)

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
