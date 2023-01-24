package utils

import (
	"github.com/go-redis/redis"
)

type RedisKeyVal struct {
	Repo    string `json:"repo"`
	Version string `json:"version"`
}

func InitRedisClient(secret string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "redis-helm-master.cloud-platform-go-get-module.svc.cluster.local:6379",
		Password: secret,
		DB:       0,
	})
}

func GetAllRedisKeysAndValues(rdb *redis.Client) ([]RedisKeyVal, error) {
	var keys []RedisKeyVal
	iter := rdb.Scan(0, "*", 0).Iterator()
	for iter.Next() {
		val, err := rdb.Get(iter.Val()).Result()
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
