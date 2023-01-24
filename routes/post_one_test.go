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

func TestInitPostOne(t *testing.T) {
	gin.SetMode("test")
	server := miniredis.RunT(t)

	mockRdbClient := utils.InitRedisClient(&redis.Options{Addr: server.Addr()})
	mockErroredRdbClient := utils.InitRedisClient(&redis.Options{Addr: "fake:0000"})

	tests := []struct {
		name             string
		rdb              utils.DataAccessLayer
		url              string
		apiKey           string
		expectedStatus   int
		expectedResponse string
	}{
		{
			"GIVEN the correct params THEN update the repo version data",
			mockRdbClient,
			"/update/foo/bar",
			"test-api-key",
			200,
			"{\"message\":\"foo updated to bar\"}",
		},
		{
			"GIVEN the correct params BUT redis is down THEN return an error",
			mockErroredRdbClient,
			"/update/foo/bar",
			"test-api-key",
			500,
			"{\"error\":\"Writing to Redis\"}",
		},
		{
			"GIVEN an incorrect apiKey THEN return an unauthorised error",
			mockRdbClient,
			"/update/foo/bar",
			"invalid-api-key",
			401,
			"{\"error\":\"Invalid API Key supplied\"}",
		},
		{
			"GIVEN an incorrect repo param THEN return an bad request error",
			mockRdbClient,
			"/update//bar",
			"test-api-key",
			400,
			"{\"error\":\"Repo parameter is must be provided `/update/:repo_name/:updated_version_number`\"}",
		},
		{
			"GIVEN an incorrect version param THEN return an bad request error",
			mockRdbClient,
			"/update/foo/ ",
			"test-api-key",
			400,
			"{\"error\":\"Update parameter is must be provided eg. `/update/:repo_name/:updated_version_number`\"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			InitPostOne(r, tt.rdb, "test-api-key")

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", tt.url, nil)
			req.Header.Set("X-API-Key", tt.apiKey)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedResponse, w.Body.String())
		})
	}
}
