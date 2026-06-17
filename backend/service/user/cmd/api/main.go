package main

import (
	"log"
	"tennis-league/common/http/router"
	"tennis-league/common/lib/cache"
	"tennis-league/common/lib/database"
	"tennis-league/common/security/authmiddleware"
	"tennis-league/common/security/repository"
	"tennis-league/user-service/internal/controller/authhandler"
	"tennis-league/user-service/internal/controller/dashboardhandler"
	"tennis-league/user-service/internal/controller/playerhandler"
	"tennis-league/user-service/internal/controller/userhandler"
	"tennis-league/user-service/internal/repository/postgres"
	"tennis-league/user-service/internal/repository/redis"
	"tennis-league/user-service/internal/service/auth"
	"tennis-league/user-service/internal/service/player"
	"tennis-league/user-service/internal/service/token"
	"tennis-league/user-service/internal/service/user"
)

func main() {
	serverConfig := router.LoadServerConfig()

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
	playerRepository := postgres.NewPlayerRepository(db)

	tokenService := token.NewTokenService("tennis", sessionRepository, serverConfig)
	userUC := user.NewUsecase(transactionManager, userRepo)
	authUC := auth.NewUsecase(db, userRepo, sessionRepository)
	playerUc := player.NewUsecase(transactionManager, playerRepository)

	dashboardHandler := dashboardhandler.NewDashboardHandler(playerUc)
	playerHandler := playerhandler.NewPlayerHandler(playerUc)
	authHandler := authhandler.NewAuthHandler(authUC, tokenService)
	userHandler := userhandler.NewUserHandler(userUC)

	r := router.NewRouter(serverConfig,
		authmiddleware.NewAuthMiddleware("tennis", sessionRepository),
		dashboardHandler,
		playerHandler,
		authHandler,
		userHandler)

	log.Println("Server running on :" + serverConfig.Port)
	r.Run(":" + serverConfig.Port)
}
