package main

import (
	"log"

	"github.com/go-chi/jwtauth/v5"
	"github.com/turanberker/tennis-league-service/internal/delivery/http"
	"github.com/turanberker/tennis-league-service/internal/delivery/http/handler/leaguehandler"
	"github.com/turanberker/tennis-league-service/internal/delivery/http/handler/userhandler"
	"github.com/turanberker/tennis-league-service/internal/domain/league"
	"github.com/turanberker/tennis-league-service/internal/domain/user"
	"github.com/turanberker/tennis-league-service/internal/infrastructure/persistence/postgres"
	"github.com/turanberker/tennis-league-service/internal/platform/database"
)

func main() {
	db, err := database.NewPostgres()
	if err != nil {
		log.Fatal(err)
	}

	userRepo := postgres.NewUserRepository(db)
	userUC := user.NewUsecase(db, userRepo)

	tokenAuth := jwtauth.New("HS256", []byte("secret-key"), nil)
	userHandler := userhandler.NewUserHandler(userUC, tokenAuth)

	leagueRepository := postgres.NewLeagueRepository(db)

	leagueUseCase := league.NewUsecase(db, leagueRepository)
	leagueHandler := leaguehandler.NewHandler(leagueUseCase)

	r := http.NewRouter(userHandler, leagueHandler)

	log.Println("Server running on :8500")
	r.Run(":8500")
}
