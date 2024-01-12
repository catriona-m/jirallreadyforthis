package gh

import (
	"context"
	"fmt"

	"github.com/google/go-github/v52/github"
)

func (r Repo) PullRequestIsMerged(prNumber int) (bool, error) {
	client := r.NewClient()
	isMerged, _, err := client.PullRequests.IsMerged(context.Background(), r.Owner, r.Name, prNumber)
	if err != nil {
		return false, fmt.Errorf("error checking if pull request %d is merged: %v", prNumber, err)
	}
	return isMerged, nil
}

func (r Repo) GetPullRequest(prNumber int) (*github.PullRequest, error) {
	client := r.NewClient()
	pr, _, err := client.PullRequests.Get(context.Background(), r.Owner, r.Name, prNumber)
	if err != nil {
		return nil, fmt.Errorf("error getting pull request #%d: %v", prNumber, err)
	}
	return pr, nil
}
