package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

func SaveToCache(ctx context.Context, rdb *redis.Client, rkey string, data any, ttl time.Duration) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, rkey, bytes, ttl).Err()
}

func GetFromCache(ctx context.Context, rdb *redis.Client, rkey string, dst any) (bool, error) {
	val, err := rdb.Get(ctx, rkey).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}

	err = json.Unmarshal([]byte(val), dst)
	if err != nil {
		return false, err
	}

	return true, nil
}

func DelFromCache(ctx context.Context, rdb *redis.Client, rkeys ...string) error {
	if err := rdb.Del(ctx, rkeys...).Err(); err != nil {
		return err
	}
	return nil
}

func SaveToBlacklist(ctx context.Context, rdb *redis.Client, token string, expiration time.Duration) error {
	rkey := "blacklist:" + token
	return rdb.Set(ctx, rkey, "1", expiration).Err()
}

func IsBlacklisted(ctx context.Context, rdb *redis.Client, token string) (bool, error) {
	rkey := "blacklist:" + token
	n, err := rdb.Exists(ctx, rkey).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}
