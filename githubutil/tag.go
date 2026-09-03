package githubutil

import (
	"context"
	"fmt"

	"github.com/google/go-github/v50/github"
)

func GetTagCommitSHA(client *github.Client, owner, repo, tag string) (string, error) {
	data, _, releaseByTagErr := client.Git.GetRef(context.Background(), owner, repo, "tags/"+tag)
	fmt.Printf("release by tag data for repo %s: %+v\n", repo, data)

	if releaseByTagErr != nil {
		return "", releaseByTagErr
	}

	object := data.GetObject()
	objectType := object.GetType()

	switch objectType {
	case "commit":
		return object.GetSHA(), nil

	// in the case of an annotated tag rather than a lightweight tag
	case "tag":
		annotatedTag, _, err := client.Git.GetTag(context.Background(), owner, repo, object.GetSHA())
		if err != nil {
			return "", err
		}
		return annotatedTag.GetObject().GetSHA(), nil

	default:
		return "", fmt.Errorf(
			"tag %q in repo %q points to unexpected error object type %q", tag, repo, objectType,
		)
	}
}
