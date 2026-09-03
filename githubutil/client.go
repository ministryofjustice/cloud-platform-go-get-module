package githubutil

import (
	"context"
	"fmt"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

func NewGitHubClient(token string) *github.Client {
	if token == "" {
		// fall back to unauthenticated client if no token is provided
		fmt.Println("No GitHub token provided, using unauthenticated client with limited rate and access.")
		return github.NewClient(nil)
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	return github.NewClient(tc)
}
