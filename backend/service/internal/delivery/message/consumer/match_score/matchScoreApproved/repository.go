package matchScoreApproved

import "context"

type Repository interface {
	AddPlayerPoint(ctx context.Context, AddPlayerPoint *AddPlayerPoint) error
}
