package utils

import (
	"time"

	"github.com/go-redis/redis"
)

type DataAccessLayer interface {
	Scan(uint64, string, int64) *redis.ScanCmd
	Get(string) (string, error)
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

type RedisKeyVal struct {
	Repo    string `json:"repo"`
	Version string `json:"currentVersion"`
}

type RedisDAL struct {
	client *redis.Client
}

func InitRedisClient(options *redis.Options) DataAccessLayer {
	dal := &RedisDAL{client: redis.NewClient(options)}
	return dal
}

func (r *RedisDAL) Get(key string) (string, error) {
	value, err := r.client.Get(key).Result()
	return value, err
}

func (r *RedisDAL) Scan(cursor uint64, match string, count int64) *redis.ScanCmd {
	return r.client.Scan(cursor, match, count)
}

func (r *RedisDAL) Set(key string, value any, expiration time.Duration) *redis.StatusCmd {
	return r.client.Set(key, value, expiration)
}

func GetAllRedisKeysAndValues(rdb DataAccessLayer) ([]RedisKeyVal, error) {
	var keys []RedisKeyVal
	iter := rdb.Scan(0, "*", 0).Iterator()
	for iter.Next() {
		val, err := rdb.Get(iter.Val())
		if err != nil {
			return nil, err
		}

		keyVal := RedisKeyVal{iter.Val(), val}
		keys = append(keys, keyVal)
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}

	return keys, nil
}
