package matchrequest

import "context"

type Repository interface {
	Create(ctx context.Context, matchRequest RequestDto) (string, error)
}
