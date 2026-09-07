package githubutil

import (
	"fmt"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v50/github"
)

type AppConfig struct {
	AppID          int64
	InstallationID int64
	PrivateKey     []byte
}

func NewGitHubClient(cfg AppConfig) (*github.Client, error) {
	transport, err := ghinstallation.New(
		http.DefaultTransport,
		cfg.AppID,
		cfg.InstallationID,
		cfg.PrivateKey,
	)
	if err != nil {
		return nil, fmt.Errorf("creating GitHub App transport: %w", err)
	}

	return github.NewClient(&http.Client{Transport: transport}), nil
}
