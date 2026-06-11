package repository

import (
	"context"
	"tennis-league/common/security/dto"

	"time"

	goredis "github.com/redis/go-redis/v9"
)

type SessionGetterRepository interface {
	Get(ctx context.Context, sessionID string) (*dto.Session, error)
}

const sessionTTL = 7 * 24 * time.Hour

func SessionKey(id string) string {
	return "session:" + id
}

type SessionGetterRepositoryImpl struct {
	rdb *goredis.Client
}

func NewSessionGetterRepositoryImpl(rdb *goredis.Client) *SessionGetterRepositoryImpl {
	return &SessionGetterRepositoryImpl{rdb: rdb}
}
func (r *SessionGetterRepositoryImpl) Get(ctx context.Context, sessionID string) (*dto.Session, error) {
	key := SessionKey(sessionID)

	m, err := r.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if len(m) == 0 {
		return nil, nil
	}
	playerId := m["player_id"]
	return &dto.Session{
		SessionId: sessionID,
		UserId:    m["user_id"],
		Role:      m["role"],
		PlayerId:  &playerId,
	}, nil
}
