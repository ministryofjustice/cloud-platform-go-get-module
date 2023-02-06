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

func TestInitGetAll(t *testing.T) {
	gin.SetMode("test")

	server := miniredis.RunT(t)

	mockRdbClient := utils.InitRedisClient(&redis.Options{Addr: server.Addr()})
	mockErroredRdbClient := utils.InitRedisClient(&redis.Options{Addr: "fake:0000"})

	tests := []struct {
		name             string
		rdb              utils.DataAccessLayer
		expectedStatus   int
		expectedResponse string
	}{
		{
			"GIVEN the correct param AND the param is in redis THEN return repos version data",
			mockRdbClient,
			200,
			"[{\"repo\":\"foo\",\"currentVersion\":\"bar\"},{\"repo\":\"ping\",\"currentVersion\":\"pong\"},{\"repo\":\"sing\",\"currentVersion\":\"song\"}]",
		},
		{
			"GIVEN the correct param BUT redis is down THEN return an error",
			mockErroredRdbClient,
			500,
			"{\"error\":\"Reading from Redis\"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			InitGetAll(r, tt.rdb)
			tt.rdb.Set("foo", "bar", 0)
			tt.rdb.Set("ping", "pong", 0)
			tt.rdb.Set("sing", "song", 0)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedResponse, w.Body.String())
		})
	}
}
