package githubutil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-github/v50/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTagCommitSHA(t *testing.T) {
	tests := []struct {
		name      string
		repo      string
		refObject map[string]string
		tagObject map[string]string
		want      string
		wantErr   bool
	}{
		{
			"GIVEN a lightweight tag THEN return the commit sha",
			"cloud-platform-terraform-test",
			map[string]string{"type": "commit", "sha": "commit-sha"},
			nil,
			"commit-sha",
			false,
		},
		{
			"GIVEN an annotated tag THEN return the tagged commit sha",
			"cloud-platform-terraform-annotated",
			map[string]string{"type": "tag", "sha": "tag-sha"},
			map[string]string{"type": "commit", "sha": "annotated-commit-sha"},
			"annotated-commit-sha",
			false,
		},
		{
			"GIVEN an unexpected tag object type THEN return an error",
			"cloud-platform-terraform-blob",
			map[string]string{"type": "blob", "sha": "blob-sha"},
			nil,
			"",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := githubClientForTagTest(t, tt.repo, tt.refObject, tt.tagObject)

			got, err := GetTagCommitSHA(client, "ministryofjustice", tt.repo, "v1.0.0")

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func githubClientForTagTest(t *testing.T, repo string, refObject, tagObject map[string]string) *github.Client {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/repos/ministryofjustice/" + repo + "/git/ref/tags/v1.0.0":
			writeGitObject(t, w, refObject)
		case "/repos/ministryofjustice/" + repo + "/git/tags/tag-sha":
			writeGitObject(t, w, tagObject)
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

func writeGitObject(t *testing.T, w http.ResponseWriter, object map[string]string) {
	t.Helper()

	err := json.NewEncoder(w).Encode(map[string]interface{}{
		"object": map[string]string{
			"type": object["type"],
			"sha":  object["sha"],
		},
	})
	require.NoError(t, err)
}
