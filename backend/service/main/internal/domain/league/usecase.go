package league

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"tennis-league/common/lib/cache"
	"tennis-league/common/lib/database"
	customerror "tennis-league/common/lib/error"

	errorcodes "tennis-league/service/internal/domain/error_codes"
	"tennis-league/service/internal/domain/leaguecoordinator"
	"tennis-league/service/internal/domain/match"
	"tennis-league/service/internal/domain/outbox"
	"tennis-league/service/internal/domain/team"
	"tennis-league/service/internal/domain/user"
)

type Usecase struct {
	tm                    *database.TransactionManager
	cacheManager          *cache.CacheManager
	userUsecase           *user.Usecase
	teamUseCase           *team.UseCase
	matchUc               *match.UseCase
	outboxRepository      outbox.Repository
	repo                  Repository
	teamRepo              team.Repository
	matchRepo             match.Repository
	coordinatorRepository leaguecoordinator.Repository
	participantRepository ParticipantRepository
}

func NewUsecase(
	tm *database.TransactionManager,
	cacheManager *cache.CacheManager,
	teamUc *team.UseCase,
	matchUc *match.UseCase,
	userUseCase *user.Usecase,
	repo Repository,
	teamRepo team.Repository,
	matchRepo match.Repository,
	outboxRepository outbox.Repository,
	coordinatorRepository leaguecoordinator.Repository,
	participantRepository ParticipantRepository,
) *Usecase {
	return &Usecase{repo: repo,
		teamUseCase:           teamUc,
		cacheManager:          cacheManager,
		matchUc:               matchUc,
		teamRepo:              teamRepo,
		matchRepo:             matchRepo,
		coordinatorRepository: coordinatorRepository,
		userUsecase:           userUseCase,
		tm:                    tm,
		outboxRepository:      outboxRepository,
		participantRepository: participantRepository,
	}
}

func (u *Usecase) ApproveMatchScore(ctx context.Context, leagueId string, matchId string) error {

	return u.tm.WithTransaction(ctx, func(txCtx context.Context) error {
		err := u.matchUc.ApproveScore(txCtx, match.MatchSource_LEAGUE, matchId)

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
			EventType:     RoutingName_LeagueMatchApproved,
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

func (u *Usecase) IsUserCoordinator(ctx context.Context, leagueId string, userId string) (bool, error) {
	key := u.cacheManager.PrepareCacheKey("league", leagueId, "isCoordinator", userId)

	// Go burada T'nin 'bool' olduğunu fn'in dönüşünden anlıyor
	return cache.Cacheable(u.cacheManager, ctx, key, 10*time.Minute, func() (bool, error) {
		return u.coordinatorRepository.Exists(ctx, leagueId, userId)
	})

}

func (u *Usecase) GetFixture(context context.Context, leagueId string, filterParam *FixtureFilterParam) ([]*match.LeagueFixtureMatch, error) {

	var filter match.FixtureFilter
	if filterParam != nil && filterParam.TeamId != nil {
		filter.TeamId = filterParam.TeamId
	}
	return u.matchRepo.GetFixtureByLeagueId(context, leagueId, &filter)
}

func (u *Usecase) Start(ctx context.Context, leagueId string) error {

	return u.tm.WithTransaction(ctx, func(txCtx context.Context) error {

		league, err := u.repo.GetById(ctx, leagueId)
		if err != nil {
			return err
		}
		if league.Status != LeagueStatus_DRAFT {
			return customerror.NewBusinessError(http.StatusConflict,
				errorcodes.ErrLeagueAlreadyFixtureCreated,
				"Lig başlamıştır")
		}

		starter, err := u.leagueStarterGetter(*league, u.teamRepo, u.matchRepo)
		err = starter.Start(txCtx, leagueId)
		if err != nil {
			return err
		}
		err = u.repo.StartLeague(txCtx, leagueId)
		if err != nil {
			return err
		}
		leageuCacheKey := u.cacheManager.PrepareCacheKey("league", leagueId)
		err = u.cacheManager.Invalidate(txCtx, leageuCacheKey)
		return err
	})

}

func (u *Usecase) GetById(ctx context.Context, id string) (*League, error) {

	cacheKey := u.cacheManager.PrepareCacheKey("league", id)

	return cache.Cacheable(u.cacheManager, ctx, cacheKey, 1*time.Hour, func() (*League, error) {
		return u.repo.GetById(ctx, id)
	})

}

func (u *Usecase) GetAll(ctx context.Context, status *LEAGUE_STATUS) ([]*LeagueListSelect, error) {
	return u.repo.GetAll(ctx, status)
}

func (u *Usecase) Save(ctx context.Context, persistLeague *PersistLeague) (*string, error) {
	id, err := u.repo.Save(ctx, persistLeague)
	if err != nil {
		if errors.Is(err, LEAGE_WITH_NAME_EXISTS) {
			return nil, customerror.NewBusinessError(http.StatusConflict,
				errorcodes.ErrLeagueAlreadyExists, "Bu isimli bir lig tanımlıdır")
		}
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

		err = u.participantRepository.AddTeamToLeague(txCtx, createTeamDto.LeagueId, *teamId)
		totalAttendance, err := u.repo.IncreaseAttendanceCount(txCtx, createTeamDto.LeagueId)
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

func (u *Usecase) GetPlayersByLeagueId(ctx context.Context, leagueId string) ([]SingleLeagueAttendance, error) {
	return u.participantRepository.SingleLeagueAttendanceList(ctx, leagueId)
}

func (u *Usecase) IsPlayerAttandedToLeague(ctx context.Context, leagueId string, playerId string) (bool, error) {
	return u.participantRepository.IsPlayerAttendedToLeague(ctx, leagueId, playerId)
}

func (u *Usecase) AddPlayerToLeague(ctx context.Context, leagueId string, playerId string) (*int32, error) {

	var response *int32
	err := u.tm.WithTransaction(ctx, func(txCtx context.Context) error {

		err := u.participantRepository.AddPlayerToLeague(txCtx, leagueId, playerId)
		if err != nil {
			return err
		}

		totalAttendance, err := u.repo.IncreaseAttendanceCount(txCtx, leagueId)

		if err != nil {
			return err
		}
		response = totalAttendance

		cacheKey := u.cacheManager.PrepareCacheKey("league", leagueId)
		err = u.cacheManager.Invalidate(txCtx, cacheKey)
		if err != nil {
			return err
		}
		return nil

	})
	if err != nil {
		return nil, err
	}

	return response, nil
}
