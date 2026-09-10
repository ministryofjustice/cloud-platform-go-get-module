package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/google/go-github/v50/github"
	"github.com/ministryofjustice/cloud-platform-go-get-module/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitPostOne(t *testing.T) {
	gin.SetMode("test")
	server := miniredis.RunT(t)
	githubClient := githubClientForPostOneTest(t)

	mockRdbClient := utils.InitRedisClient(&redis.Options{Addr: server.Addr()})
	mockErroredRdbClient := utils.InitRedisClient(&redis.Options{Addr: "fake:0000"})

	tests := []struct {
		name             string
		rdb              utils.DataAccessLayer
		githubClient     *github.Client
		url              string
		apiKey           string
		expectedStatus   int
		expectedResponse string
	}{
		{
			"GIVEN the correct params THEN update the repo version data",
			mockRdbClient,
			githubClient,
			"/update/foo/bar",
			"test-api-key",
			200,
			"{\"message\":\"foo updated to bar\"}",
		},
		{
			"GIVEN the correct params BUT redis is down THEN return an error",
			mockErroredRdbClient,
			githubClient,
			"/update/foo/bar",
			"test-api-key",
			500,
			"{\"error\":\"Writing to Redis\"}",
		},
		{
			"GIVEN an incorrect apiKey THEN return an unauthorised error",
			mockRdbClient,
			nil,
			"/update/foo/bar",
			"invalid-api-key",
			401,
			"{\"error\":\"Invalid API Key supplied\"}",
		},
		{
			"GIVEN an incorrect repo param THEN return an bad request error",
			mockRdbClient,
			nil,
			"/update//bar",
			"test-api-key",
			400,
			"{\"error\":\"Repo parameter is must be provided `/update/:repo_name/:updated_version_number`\"}",
		},
		{
			"GIVEN an incorrect version param THEN return an bad request error",
			mockRdbClient,
			nil,
			"/update/foo/ ",
			"test-api-key",
			400,
			"{\"error\":\"Update parameter is must be provided eg. `/update/:repo_name/:updated_version_number`\"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			InitPostOne(r, tt.rdb, tt.githubClient, "test-api-key")

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", tt.url, nil)
			req.Header.Set("X-API-Key", tt.apiKey)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedResponse, w.Body.String())
		})
	}
}

func githubClientForPostOneTest(t *testing.T) *github.Client {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]interface{}{
			"object": map[string]string{
				"type": "commit",
				"sha":  "test-commit-sha",
			},
		})
		require.NoError(t, err)
	}))
	t.Cleanup(server.Close)

	client := github.NewClient(server.Client())
	baseURL, err := url.Parse(server.URL + "/")
	require.NoError(t, err)
	client.BaseURL = baseURL

	return client
}
