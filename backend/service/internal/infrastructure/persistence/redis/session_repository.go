package redis

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"time"

"github.com/turanberker/tennis-league-service/internal/domain/session"
	goredis "github.com/redis/go-redis/v9"
)

type SessionRepository struct {
	rdb *goredis.Client
}

func NewSessionRepository(rdb *goredis.Client) *SessionRepository {
	return &SessionRepository{rdb: rdb}
}

const sessionTTL = 1 * time.Hour

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



func (r *SessionRepository) Start(ctx context.Context, userID int64, role string) (*session.Session, error) {
	id, err := generateSessionID()
	if err != nil {
		return nil, err
	}

	key := sessionKey(id)

	err = r.rdb.HSet(ctx, key, map[string]any{
		"user_id": strconv.FormatInt(userID, 10),
		"role":    role,
	}).Err()
	if err != nil {
		return nil, err
	}

	if err := r.rdb.Expire(ctx, key, sessionTTL).Err(); err != nil {
		return nil, err
	}

	return &session.Session{
		SessionId: id,
		UserId:    userID,
		Role:      role,
	}, nil
}

func (r *SessionRepository) Get(ctx context.Context, sessionID string) (*session.Session, error) {
	key := sessionKey(sessionID)

	m, err := r.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if len(m) == 0 {
		return nil, nil
	}

	userID, _ := strconv.ParseInt(m["user_id"], 10, 64)

	return &session.Session{
		SessionId: sessionID,
		UserId:    userID,
		Role:      m["role"],
	}, nil
}

func (r *SessionRepository) Delete(ctx context.Context, sessionID string) error {
	return r.rdb.Del(ctx, sessionKey(sessionID)).Err()
}


func (r *SessionRepository) Refresh(ctx context.Context, sessionID string) error {
	return r.rdb.Expire(ctx, sessionKey(sessionID), sessionTTL).Err()
}
