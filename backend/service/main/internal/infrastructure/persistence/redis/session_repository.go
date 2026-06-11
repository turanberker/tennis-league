package redis

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"tennis-league/common/security/dto"
	"tennis-league/common/security/repository"
	"tennis-league/service/internal/domain/session"

	goredis "github.com/redis/go-redis/v9"
)

type SessionRepository struct {
	repository.SessionGetterRepository
	rdb *goredis.Client
}

func NewSessionRepository(sessionGetterRepository repository.SessionGetterRepository, rdb *goredis.Client) *SessionRepository {
	return &SessionRepository{SessionGetterRepository: sessionGetterRepository, rdb: rdb}
}

const sessionTTL = 7 * 24 * time.Hour

func sessionKey(id string) string {
	return "session:" + id
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (r *SessionRepository) Start(ctx context.Context, startSessionInput *session.StartSessionInput) (*dto.Session, error) {
	id, err := generateSessionID()
	if err != nil {
		return nil, err
	}

	key := sessionKey(id)

	err = r.rdb.HSet(ctx, key, map[string]any{
		"user_id":   startSessionInput.UserId,
		"role":      startSessionInput.Role,
		"player_id": startSessionInput.PlayerId,
	}).Err()
	if err != nil {
		return nil, err
	}

	if err := r.rdb.Expire(ctx, key, sessionTTL).Err(); err != nil {
		return nil, err
	}

	return &dto.Session{
		SessionId: id,
		UserId:    startSessionInput.UserId,
		Role:      startSessionInput.Role,
		PlayerId:  startSessionInput.PlayerId,
	}, nil
}

func (r *SessionRepository) Delete(ctx context.Context, sessionID string) error {
	return r.rdb.Del(ctx, sessionKey(sessionID)).Err()
}

func (r *SessionRepository) Refresh(ctx context.Context, sessionID string) error {
	return r.rdb.Expire(ctx, sessionKey(sessionID), sessionTTL).Err()
}
