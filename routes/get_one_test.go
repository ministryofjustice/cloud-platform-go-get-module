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
			"{\"currentVersion\":\"bar\",\"repo\":\"foo\"}",
			"/foo",
		},
		{
			"GIVEN a incorrect param AND the param is in redis THEN return a 500 error",
			mockRdbClient,
			500,
			"{\"error\":\"Reading from Redis, check the repo_name param for typos: error\"}",
			"/error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			InitGetOne(r, tt.rdb)

			tt.rdb.Set("foo", "bar", 0)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tt.urlParam, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedResponse, w.Body.String())
		})
	}
}
