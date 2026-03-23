package league

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"

	customerror "github.com/turanberker/tennis-league-service/internal/domain/error"
	"github.com/turanberker/tennis-league-service/internal/domain/leaguecoordinator"
	"github.com/turanberker/tennis-league-service/internal/domain/match"
	"github.com/turanberker/tennis-league-service/internal/domain/outbox"
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
	teamUseCase           *team.UseCase
	matchUc               *match.UseCase
	outboxRepository      outbox.Repository
	repo                  Repository
	teamRepo              team.Repository
	matchRepo             match.Repository
	scoreBoardRepo        scoreboard.Repository
	coordinatorRepository leaguecoordinator.Repository
}

func NewUsecase(
	tm *database.TransactionManager,
	teamUc *team.UseCase,
	matchUc *match.UseCase,
	userUseCase *user.Usecase,
	repo Repository,
	teamRepo team.Repository,
	matchRepo match.Repository,
	outboxRepository outbox.Repository,
	scoreBoardRepo scoreboard.Repository,
	coordinatorRepository leaguecoordinator.Repository,
) *Usecase {
	return &Usecase{repo: repo,
		teamUseCase:           teamUc,
		matchUc:               matchUc,
		teamRepo:              teamRepo,
		matchRepo:             matchRepo,
		scoreBoardRepo:        scoreBoardRepo,
		coordinatorRepository: coordinatorRepository,
		userUsecase:           userUseCase,
		tm:                    tm,
		outboxRepository:      outboxRepository,
	}
}

func (u *Usecase) ApproveMatchScore(ctx context.Context, leagueId string, matchId string) error {

	return u.tm.WithTransaction(ctx, func(txCtx context.Context) error {
		err := u.matchUc.ApproveScore(txCtx, match.MatchSource_TOURNAMENT, matchId)

		if err != nil {
			return err
		}

		event := LeagueMatchApprovedEvent{
			LeagueId: leagueId,
			MatchId:  matchId,
		}
		payload, _ := json.Marshal(event)
		outboxEntity := &outbox.PersistEntity{
			AggregateType: "match",
			AggregateID:   matchId,
			EventType:     "LeagueMatchApproved",
			Payload:       payload,
		}
		err = u.outboxRepository.Save(txCtx, outboxEntity)
		if err != nil {
			return err
		}

		return nil
	})

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

func (u *Usecase) CreateFixture(ctx context.Context, leagueId string) error {

	return u.tm.WithTransaction(ctx, func(txCtx context.Context) error {

		created, err := u.repo.IsFixtureCreated(txCtx, leagueId)
		if err != nil {
			return err
		}
		if created {
			return customerror.NewBussinnessError(http.StatusConflict,
				customerror.ErrLeagueAlreadyFixtureCreated,
				"Fikstür zaten oluşturulmuş")
		}
		//TODO Burada ligin Single- double olmasına göre işlem yapılacak
		teams, err := u.teamRepo.GetByLeagueId(txCtx, leagueId)
		var bulkInsert match.BulkInsertMatches
		var matches []match.SideIds
		var teamIds []string

		for i := 0; i < len(teams); i++ {

			teamIds = append(teamIds, teams[i].ID)
			for j := i + 1; j < len(teams); j++ { // j=i+1 → tekrar ve kendisiyle maç yok
				team1Id := teams[i].ID
				team2Id := teams[j].ID

				// 50% ihtimalle takımların yerini değiştir
				if rand.Intn(2) == 0 {
					team1Id, team2Id = team2Id, team1Id
				}

				match := match.SideIds{
					Side1: team1Id,
					Side2: team2Id,
				}

				matches = append(matches, match)
			}
		}
		//Maçların sırasını karıştır (Opsiyonel ama daha profesyonel bir fikstür sağlar)
		rand.Shuffle(len(matches), func(i, j int) {
			matches[i], matches[j] = matches[j], matches[i]
		})

		bulkInsert.Sides = matches
		bulkInsert.Type = match.MatchType{Id: &leagueId,
			Source: match.MatchSource_TOURNAMENT,
			Type:   match.MatchType_DOUBLE,
		}

		err = u.repo.StartLeague(txCtx, leagueId)
		if err != nil {
			return err
		}
		err = u.matchRepo.SaveBulkMatches(txCtx, &bulkInsert)
		if err != nil {
			return err
		}
		return u.scoreBoardRepo.SaveFixture(txCtx, leagueId, teamIds)
	})

}

func (u *Usecase) GetById(ctx context.Context, id string) (*League, error) {

	return u.repo.GetById(ctx, id)
}

func (u *Usecase) GetAll(ctx context.Context, status *LEAGUE_STATUS) ([]*LeagueListSelect, error) {
	return u.repo.GetAll(ctx, status)
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

func (u *Usecase) CreateTeam(ctx context.Context, createTeamDto *CreateTeamRequestDto) (*CreateTeamResponseDto, error) {

	var response CreateTeamResponseDto

	err := u.tm.WithTransaction(ctx, func(txCtx context.Context) error {

		teamId, err := u.teamUseCase.Save(txCtx, &team.CreateTeamRequest{
			LeagueID:  createTeamDto.LeagueId,
			Name:      createTeamDto.Name,
			PlayerIDs: createTeamDto.PlayerIDs,
		})
		if err != nil {
			return err
		}

		response.TeamId = *teamId

		totalAttendance, err := u.repo.IncreaseAttandanceCount(txCtx, createTeamDto.LeagueId)

		if err != nil {
			return err
		}
		response.TotalAttendance = *totalAttendance
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &response, nil
}
