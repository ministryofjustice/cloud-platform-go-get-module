package init_app

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v50/github"
	"github.com/ministryofjustice/cloud-platform-go-get-module/utils"
)

var owner = "ministryofjustice"

func InitDataClient(dataAddr, dataPassword string) utils.DataAccessLayer {
	return initRedis(dataAddr, dataPassword)
}

func InitData(dataClient utils.DataAccessLayer) {
	client := github.NewClient(nil)
	repos, err := getRepos(client)

	if err != nil {
		log.Println("Error getting repo names: ", err)
	}

	for _, repo := range repos {
		release, _, releaseErr := client.Repositories.GetLatestRelease(context.Background(), owner, *repo.Name)
		if releaseErr != nil {
			log.Println("Error getting latest release", releaseErr)
			continue
		}

		latestVersion := release.GetTagName()
		fmt.Println(*repo.Name, latestVersion)

		dataErr := dataClient.Set(*repo.Name, latestVersion, 0).Err()
		if dataErr != nil {
			log.Println("Error setting version: ", dataErr)
		}
	}
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
