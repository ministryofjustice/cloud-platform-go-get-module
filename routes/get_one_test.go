package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/ministryofjustice/cloud-platform-go-get-module/utils"
	"github.com/stretchr/testify/assert"
)

func TestInitGetOne(t *testing.T) {
	gin.SetMode("test")
	server := miniredis.RunT(t)

	mockRdbClient := utils.InitRedisClient(&redis.Options{Addr: server.Addr()})
	mockErroredRdbClient := utils.InitRedisClient(&redis.Options{Addr: "fake:0000"})

	tests := []struct {
		name             string
		rdb              utils.DataAccessLayer
		expectedStatus   int
		expectedResponse string
		urlParam         string
	}{
		{
			"GIVEN the correct param AND the param is in redis THEN return repo version data",
			mockRdbClient,
			200,
			"{\"currentVersion\":\"bar\",\"repo\":\"foo\",\"sha\":\"abc123\"}",
			"/foo",
		},
		{
			"GIVEN the correct param BUT redis is down THEN return a 500 error",
			mockErroredRdbClient,
			500,
			"{\"error\":\"Reading from Redis\"}",
			"/foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			InitGetOne(r, tt.rdb)

			tt.rdb.HMSet("foo", map[string]interface{}{"currentVersion": "bar", "sha": "abc123"})
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tt.urlParam, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedResponse, w.Body.String())
		})
	}
}

// HGETALL returns an empty map and no error for a missing key, unlike GET which
// returns redis.Nil. Without an explicit emptiness check the handler 200s on a
// repo that isn't stored.
func TestInitGetOne_MissingRepoReturnsNotFound(t *testing.T) {
	gin.SetMode("test")
	server := miniredis.RunT(t)

	mockRdbClient := utils.InitRedisClient(&redis.Options{Addr: server.Addr()})

	// Seeded directly rather than through the DAL so this test can be written
	// before HMSet exists on the interface.
	server.HSet("cloud-platform-terraform-s3-bucket", "currentVersion", "4.9.1", "sha", "9f2c1ab")

	r := gin.New()
	InitGetOne(r, mockRdbClient)

	w := httptest.NewRecorder()
	// A repo that was never stored, while the keyspace is non-empty.
	req, _ := http.NewRequest("GET", "/cloud-platform-terraform-unknown", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "{\"error\":\"Repo not found: cloud-platform-terraform-unknown\"}", w.Body.String())
}
