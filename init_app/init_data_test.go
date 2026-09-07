package init_app

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/go-github/v50/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockDataAccessLayer struct {
	values map[string]map[string]interface{}
	err    error
}

func (m *mockDataAccessLayer) Scan(cursor uint64, match string, count int64) *redis.ScanCmd {
	return redis.NewScanCmd(nil, cursor, match, count)
}

func (m *mockDataAccessLayer) Get(key string) (string, error) {
	return "", nil
}

func (m *mockDataAccessLayer) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return redis.NewStatusResult("OK", nil)
}

func (m *mockDataAccessLayer) HGetAll(key string) (map[string]string, error) {
	return nil, nil
}

func (m *mockDataAccessLayer) HMSet(key string, values map[string]interface{}) *redis.StatusCmd {
	if m.values == nil {
		m.values = make(map[string]map[string]interface{})
	}
	m.values[key] = values

	return redis.NewStatusResult("OK", m.err)
}

func TestInitData(t *testing.T) {
	tests := []struct {
		name             string
		githubStatusCode int
		dataErr          error
		wantErr          bool
		wantValues       map[string]map[string]interface{}
	}{
		{
			"GIVEN GitHub returns repo release data THEN store latest version and sha",
			http.StatusOK,
			nil,
			false,
			map[string]map[string]interface{}{
				"cloud-platform-terraform-test": {
					"currentVersion": "v1.2.3",
					"sha":            "abc123",
				},
			},
		},
		{
			"GIVEN GitHub search returns an error THEN return an error",
			http.StatusInternalServerError,
			nil,
			true,
			nil,
		},
		{
			"GIVEN Redis returns an error THEN continue without returning an error",
			http.StatusOK,
			errors.New("redis error"),
			false,
			map[string]map[string]interface{}{
				"cloud-platform-terraform-test": {
					"currentVersion": "v1.2.3",
					"sha":            "abc123",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataClient := &mockDataAccessLayer{err: tt.dataErr}
			githubClient := githubClientForInitDataTest(t, tt.githubStatusCode)

			err := InitData(dataClient, githubClient)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.True(t, reflect.DeepEqual(tt.wantValues, dataClient.values), "got %v, want %v", dataClient.values, tt.wantValues)
		})
	}
}

func githubClientForInitDataTest(t *testing.T, searchStatusCode int) *github.Client {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/search/repositories":
			w.WriteHeader(searchStatusCode)
			if searchStatusCode == http.StatusOK {
				writeJSON(t, w, map[string]interface{}{
					"items": []map[string]string{{"name": "cloud-platform-terraform-test"}},
				})
			}
		case "/repos/ministryofjustice/cloud-platform-terraform-test/releases/latest":
			writeJSON(t, w, map[string]string{"tag_name": "v1.2.3"})
		case "/repos/ministryofjustice/cloud-platform-terraform-test/git/ref/tags/v1.2.3":
			writeJSON(t, w, map[string]interface{}{
				"object": map[string]string{"type": "commit", "sha": "abc123"},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)

	client := github.NewClient(server.Client())
	baseURL, err := url.Parse(server.URL + "/")
	require.NoError(t, err)
	client.BaseURL = baseURL

	return client
}

func writeJSON(t *testing.T, w http.ResponseWriter, data interface{}) {
	t.Helper()

	err := json.NewEncoder(w).Encode(data)
	require.NoError(t, err)
}
