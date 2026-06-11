package scoreboard

import "context"

type UseCase struct {
	r Repository
}

func NewUseCase(r Repository) *UseCase {
	return &UseCase{r: r}
}

func (u *UseCase) GetScoreBoard(ctx context.Context, leagueId string) ([]*ScoreBoard, error) {
	return u.r.GetScoreBoard(ctx, leagueId)
}
