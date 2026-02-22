package main

import (
	"log"

	"github.com/go-chi/jwtauth/v5"
	"github.com/turanberker/tennis-league-service/internal/delivery/http"
	"github.com/turanberker/tennis-league-service/internal/delivery/http/handler/leaguehandler"
	"github.com/turanberker/tennis-league-service/internal/delivery/http/handler/matchhandler"
	"github.com/turanberker/tennis-league-service/internal/delivery/http/handler/playerhandler"
	"github.com/turanberker/tennis-league-service/internal/delivery/http/handler/userhandler"
	"github.com/turanberker/tennis-league-service/internal/domain/league"
	"github.com/turanberker/tennis-league-service/internal/domain/match"
	"github.com/turanberker/tennis-league-service/internal/domain/player"
	"github.com/turanberker/tennis-league-service/internal/domain/scoreboard"
	"github.com/turanberker/tennis-league-service/internal/domain/team"
	"github.com/turanberker/tennis-league-service/internal/domain/user"
	"github.com/turanberker/tennis-league-service/internal/infrastructure/persistence/postgres"
	"github.com/turanberker/tennis-league-service/internal/platform/database"
)

func main() {

	matchhandler.RegisterSetValidations()

	db, err := database.NewPostgres()
	if err != nil {
		log.Fatal(err)
	}

	userRepo := postgres.NewUserRepository(db)
	userUC := user.NewUsecase(db, userRepo)

	tokenAuth := jwtauth.New("HS256", []byte("secret-key"), nil)
	userHandler := userhandler.NewUserHandler(userUC, tokenAuth)

	leagueRepository := postgres.NewLeagueRepository(db)
	teamRepository := postgres.NewTeamRepository(db)
	teamPlayerRepository := postgres.NewTeamPlayerRepository(db)
	matchRepository := postgres.NewMatchRepository(db)
	matchSetRepository := postgres.NewMatchSetRepository(db)
	scoreBoardRepository := postgres.NewScoreBoardRepository(db)

	leagueUseCase := league.NewUsecase(db, leagueRepository, teamRepository, matchRepository, scoreBoardRepository)
	teamUseCase := team.NewUseCase(db, teamRepository, teamPlayerRepository)
	matchUseCase := match.NewUseCase(db, matchRepository, matchSetRepository)
	scoreBaordUc := scoreboard.NewUseCase(scoreBoardRepository)
	leagueHandler := leaguehandler.NewHandler(leagueUseCase, teamUseCase, scoreBaordUc)

	playerRepository := postgres.NewPlayerRepository(db)
	playerUc := player.NewUsecase(db, playerRepository)
	playerhandler := playerhandler.NewPlayerHandler(playerUc)
	matchHandler := matchhandler.NewMatchHandler(matchUseCase)
	r := http.NewRouter(userHandler,
		leagueHandler,
		playerhandler,
		matchHandler)

	log.Println("Server running on :8500")
	r.Run(":8500")
}
