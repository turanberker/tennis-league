package main

import (
	"log"
	"tennis-league/common/http/router"
	"tennis-league/common/lib/cache"
	"tennis-league/common/lib/database"
	authmiddleware "tennis-league/common/security/auth"
	"tennis-league/common/security/repository"
	"tennis-league/user-service/internal/controller/authhandler"
	"tennis-league/user-service/internal/controller/userhandler"
	"tennis-league/user-service/internal/repository/postgres"
	"tennis-league/user-service/internal/repository/redis"
	"tennis-league/user-service/internal/service/auth"
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

	tokenService := token.NewTokenService("tennis", sessionRepository, serverConfig)
	userUC := user.NewUsecase(transactionManager, userRepo)
	authUC := auth.NewUsecase(db, userRepo, sessionRepository)

	authHandler := authhandler.NewAuthHandler(authUC, tokenService)
	userHandler := userhandler.NewUserHandler(userUC)

	r := router.NewRouter(serverConfig,
		authmiddleware.NewAuthMiddleware("tennis", sessionRepository),

		authHandler,

		userHandler)

	log.Println("Server running on :" + serverConfig.Port)
	r.Run(":" + serverConfig.Port)
}
