package cache

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type CacheManager struct {
	rdb *goredis.Client
}

func NewCacheManager(rdb *goredis.Client) *CacheManager {
	return &CacheManager{rdb: rdb}
}

type cacheKey string

func (cm *CacheManager) PrepareCacheKey(parts ...string) cacheKey {
	return cacheKey("cache:" + strings.Join(parts, ":"))
}

func (cm *CacheManager) Invalidate(ctx context.Context, key cacheKey) error {

	return cm.rdb.Del(ctx, string(key)).Err()
}

func (cm *CacheManager) InvalidateByPrefix(ctx context.Context, keyPrefix string) error {
	iter := cm.rdb.Scan(ctx, 0, string(keyPrefix)+"*", 0).Iterator()
	iter.Next(ctx)
	for iter.Next(ctx) {
		err := cm.rdb.Del(ctx, iter.Val()).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func PutCache[T any](cm *CacheManager, ctx context.Context, key cacheKey, value T, ttl time.Duration) error {
	return putCache(cm, ctx, key, value, ttl)
}

func putCache[T any](cm *CacheManager, ctx context.Context, key cacheKey, value T, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return cm.rdb.Set(ctx, string(key), data, ttl).Err()
}

func Cacheable[T any](cm *CacheManager, ctx context.Context, key cacheKey, ttl time.Duration, fn func() (T, error)) (T, error) {
	var result T

	// 1. Önbellekte var mı bak?
	val, err := cm.rdb.Get(ctx, string(key)).Result()
	if err == nil {
		// Önbellekte bulundu, serileştirilmiş veriyi nesneye dönüştür
		// --- KRİTİK NOKTA ---
		// Eğer T bir pointer ise ve nil ise, yeni bir instance oluşturmalıyız
		rv := reflect.ValueOf(&result).Elem()
		if rv.Kind() == reflect.Ptr && rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}

		err = json.Unmarshal([]byte(val), &result)
		if err == nil {
			return result, nil
		}
		// Unmarshal hatası alırsan loglayıp devam edebilirsin (cache miss gibi davran)
	}

	// 2. Önbellekte yoksa fonksiyonu çalıştır (Asıl kaynak: DB vb.)
	result, err = fn()
	if err != nil {
		return result, err
	}

	err = putCache(cm, ctx, key, result, ttl)

	if err != nil {
		return result, nil
	}

	return result, nil
}
