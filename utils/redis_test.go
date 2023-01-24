package utils

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis"
)

type MockRedisGetDAL struct {
	client *redis.Client
}

func (r MockRedisGetDAL) Get(key string) (string, error) {
	return "", errors.New("big redis error")
}

func (r MockRedisGetDAL) Scan(cursor uint64, match string, count int64) *redis.ScanCmd {
	return r.client.Scan(cursor, match, count)
}

func (r MockRedisGetDAL) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return r.client.Set(key, value, expiration)
}

func TestGetAllRedisKeysAndValues(t *testing.T) {
	server := miniredis.RunT(t)

	mockRdbClient := InitRedisClient(&redis.Options{Addr: server.Addr()})
	mockErroredRdbClient := InitRedisClient(&redis.Options{Addr: "fake:0000"})

	mockGetFn := &MockRedisGetDAL{client: redis.NewClient(&redis.Options{Addr: server.Addr()})}

	tests := []struct {
		name      string
		rdb       DataAccessLayer
		want      []RedisKeyVal
		wantErr   bool
		setValues bool
	}{
		{
			"GIVEN a redis cluster with multiple keys THEN return all the keys",
			mockRdbClient,
			[]RedisKeyVal{{"foo", "bar"}, {"test", "v1"}},
			false,
			true,
		},
		{
			"GIVEN a redis cluster which raises an error THEN return an error",
			mockErroredRdbClient,
			nil,
			true,
			false,
		},
		{
			"GIVEN a redis cluster which returns an error when trying to get a key THEN return an error",
			mockGetFn,
			nil,
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setValues {
				for _, val := range tt.want {
					tt.rdb.Set(val.Repo, val.Version, 0)
				}
			}

			got, err := GetAllRedisKeysAndValues(tt.rdb)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllRedisKeysAndValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllRedisKeysAndValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
