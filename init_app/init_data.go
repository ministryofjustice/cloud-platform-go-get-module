package init_app

import (
	"context"
	"fmt"

	"github.com/google/go-github/v50/github"
	"github.com/ministryofjustice/cloud-platform-go-get-module/githubutil"
	"github.com/ministryofjustice/cloud-platform-go-get-module/utils"
)

var owner = "ministryofjustice"

func InitDataClient(dataAddr, dataPassword string) utils.DataAccessLayer {
	return initRedis(dataAddr, dataPassword)
}

func InitData(dataClient utils.DataAccessLayer, githubClient *github.Client) error {
	repos, err := getRepos(githubClient)

	if err != nil {
		return fmt.Errorf("error getting repo data from github API: %v", err)
	}

	for _, repo := range repos {
		release, _, releaseErr := githubClient.Repositories.GetLatestRelease(context.Background(), owner, *repo.Name)
		if releaseErr != nil {
			fmt.Printf("error getting latest release: %v", releaseErr)
			continue
		}

		latestVersion := release.GetTagName()

		sha, shaErr := githubutil.GetTagCommitSHA(githubClient, owner, *repo.Name, latestVersion)
		if shaErr != nil {
			fmt.Printf("error getting sha: %v", shaErr)
			continue
		}

		dataErr := dataClient.HMSet(*repo.Name, map[string]interface{}{"currentVersion": latestVersion, "sha": sha}).Err()
		if dataErr != nil {
			fmt.Printf("error setting version: %v", dataErr)
			continue
		}
	}

	return nil
}

func getRepos(client *github.Client) ([]*github.Repository, error) {

	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{PerPage: 50},
	}

	// get all pages of results
	var allRepos []*github.Repository

	for {
		repos, resp, err := client.Search.Repositories(context.Background(), "cloud-platform-terraform- in:name archived:false is:public org:ministryofjustice", opt)

		if err != nil {
			return nil, err
		}

		allRepos = append(allRepos, repos.Repositories...)

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return allRepos, nil
}
