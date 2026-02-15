package league

import (
	"context"
	"database/sql"
	"errors"

	"github.com/turanberker/tennis-league-service/internal/domain/match"
	"github.com/turanberker/tennis-league-service/internal/domain/team"
)

var ErrNameFieldRequired = errors.New("Name can not be null or empty string")
var ErrNameLenghtError = errors.New("Name size must between 5 and 75 characters")

type Usecase struct {
	db        *sql.DB
	repo      Repository
	teamRepo  team.Repository
	matchRepo match.Repository
}

func NewUsecase(db *sql.DB, repo Repository, teamRepo team.Repository, matchRepo match.Repository) *Usecase {
	return &Usecase{db: db, repo: repo, teamRepo: teamRepo, matchRepo: matchRepo}
}

func (u *Usecase) SetFitxtureCreatedDate(ctx context.Context, leagueId string) error {
	created, err := u.repo.IsFixtureCreated(ctx, leagueId)

	if created {
		return errors.New("Fikstür zaten oluşturulmuş")
	}

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	u.repo.SetFitxtureCreatedDate(ctx, tx, leagueId)
	teams, err := u.teamRepo.GetByLeagueId(ctx, leagueId)

	if err != nil {
		tx.Rollback()
		return err
	}

	var matches []match.PersistLeagueMatch
	for i := 0; i < len(teams); i++ {
		for j := i + 1; j < len(teams); j++ { // j=i+1 → tekrar ve kendisiyle maç yok
			match := match.PersistLeagueMatch{
				LeagueId: leagueId,
				Team1Id:  teams[i].ID,
				Team2Id:  teams[j].ID,
			}
			matches = append(matches, match)
		}
	}

	defer tx.Rollback()
	u.matchRepo.SaveLeagueMatches(ctx, tx, matches)
	tx.Commit()

	return nil
}

func (u *Usecase) GetById(ctx context.Context, id int64) (*League, error) {

	return u.repo.GetById(ctx, id)
}

func (u *Usecase) GetAll(ctx context.Context, name string) ([]*League, error) {

	return u.repo.GetAll(ctx, name)
}

func (u *Usecase) Save(ctx context.Context, persistLeague *PersistLeague) (*string, error) {
	if persistLeague.Name == "" {
		return nil, ErrNameFieldRequired
	}

	if len(persistLeague.Name) < 5 || len(persistLeague.Name) > 75 {
		return nil, ErrNameLenghtError
	}

	return u.repo.Save(ctx, persistLeague)
}
