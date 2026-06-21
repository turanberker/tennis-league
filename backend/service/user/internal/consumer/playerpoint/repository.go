package playerpoint

import "context"

type EarnedPointRepository interface {
	AddPlayerPoint(ctx context.Context, AddPlayerPoint *AddPlayerPoint) error
}

type PlayerRepository interface {
	GetPlayerPoints(ctx context.Context, Ids []string) ([]PlayerPoints, error)

	DecreaseDoublePoint(ctx context.Context, playerId string, change int) (int, error)

	IncreaseDoublePoint(ctx context.Context, playerId string, change int) (int, error)
}
