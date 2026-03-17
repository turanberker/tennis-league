package league

import (
	"context"
	"errors"
	"net/http"

	customerror "github.com/turanberker/tennis-league-service/internal/domain/error"
	"github.com/turanberker/tennis-league-service/internal/domain/leaguecoordinator"
	"github.com/turanberker/tennis-league-service/internal/domain/match"
	"github.com/turanberker/tennis-league-service/internal/domain/scoreboard"
	"github.com/turanberker/tennis-league-service/internal/domain/team"
	"github.com/turanberker/tennis-league-service/internal/domain/user"
	"github.com/turanberker/tennis-league-service/internal/platform/database"
)

var ErrNameFieldRequired = errors.New("Name can not be null or empty string")
var ErrNameLenghtError = errors.New("Name size must between 5 and 75 characters")

type Usecase struct {
	tm                    *database.TransactionManager
	userUsecase           *user.Usecase
	repo                  Repository
	teamRepo              team.Repository
	matchRepo             match.Repository
	scoreBoardRepo        scoreboard.Repository
	coordinatorRepository leaguecoordinator.Repository
}

func (u *Usecase) AddNewCoordinator(ctx context.Context, leagueId string, userId string) (*bool, error) {

	var isAdded bool

	// TransactionManager (tm) üzerinden süreci sarmalıyoruz
	err := u.tm.WithTransaction(ctx, func(txCtx context.Context) error {

		// 1. Koordinatörü ekle (txCtx kullanarak transaction'ı taşıyoruz)
		added, err := u.coordinatorRepository.Add(txCtx, leagueId, userId)
		if err != nil {
			return err
		}

		// 2. Eğer eklendiyse rolü güncelle
		if *added {
			// DİKKAT: userUsecase de txCtx almalı ki aynı transaction'da kalsın
			err = u.userUsecase.SetUserAsCoordinator(txCtx, userId)
			if err != nil {
				return err
			}
			isAdded = true
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &isAdded, nil

}

func (u *Usecase) IsUserCoordinator(context context.Context, leagueId string, userId string) (bool, error) {
	return u.coordinatorRepository.Exists(context, leagueId, userId)
}

func (u *Usecase) GetFixture(context context.Context, leagueId string) ([]*match.LeagueFixtureMatch, error) {
	return u.matchRepo.GetFixtureByLeagueId(context, leagueId)
}

func NewUsecase(
	tm *database.TransactionManager,
	repo Repository,
	teamRepo team.Repository,
	matchRepo match.Repository,
	scoreBoardRepo scoreboard.Repository,
	coordinatorRepository leaguecoordinator.Repository,
	userUseCase *user.Usecase) *Usecase {
	return &Usecase{repo: repo,
		teamRepo:              teamRepo,
		matchRepo:             matchRepo,
		scoreBoardRepo:        scoreBoardRepo,
		coordinatorRepository: coordinatorRepository,
		userUsecase:           userUseCase,
		tm:                    tm,
	}
}

func (u *Usecase) SetFitxtureCreatedDate(ctx context.Context, leagueId string) error {

	return u.tm.WithTransaction(ctx, func(txCtx context.Context) error {

		created, err := u.repo.IsFixtureCreated(txCtx, leagueId)
		if err != nil {
			return err
		}
		if created {
			return errors.New("Fikstür zaten oluşturulmuş")
		}

		u.repo.SetFitxtureCreatedDate(txCtx, leagueId)
		teams, err := u.teamRepo.GetByLeagueId(txCtx, leagueId)

		var matches []*match.PersistLeagueMatch
		var teamIds []string

		for i := 0; i < len(teams); i++ {

			teamIds = append(teamIds, teams[i].ID)
			for j := i + 1; j < len(teams); j++ { // j=i+1 → tekrar ve kendisiyle maç yok
				match := &match.PersistLeagueMatch{
					LeagueId: leagueId,
					Team1Id:  teams[i].ID,
					Team2Id:  teams[j].ID,
				}

				matches = append(matches, match)
			}
		}

		u.matchRepo.SaveLeagueMatches(txCtx, matches)
		u.scoreBoardRepo.SaveFixture(txCtx, leagueId, teamIds)

		return nil

	})

}

func (u *Usecase) GetById(ctx context.Context, id string) (*League, error) {

	return u.repo.GetById(ctx, id)
}

func (u *Usecase) GetAll(ctx context.Context, name string) ([]*League, error) {

	return u.repo.GetAll(ctx, &name)
}

func (u *Usecase) Save(ctx context.Context, persistLeague *PersistLeague) (*string, error) {
	id, err := u.repo.Save(ctx, persistLeague)
	if err != nil {
		if errors.Is(err, LEAGE_WITH_NAME_EXISTS) {
			return nil, customerror.NewBussinnessError(http.StatusConflict,
				customerror.ErrLeagueAlreadyExists, "Bu isimli bir lig tanımlıdır")
		}
	} else {
		return nil, customerror.NewInternalError(err)
	}

	return id, nil

}
