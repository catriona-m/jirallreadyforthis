package gh

import (
	"context"
	"fmt"

	"github.com/google/go-github/v52/github"
)

func (r Repo) GetIssue(issueNumber int) (*github.Issue, error) {
	client := r.NewClient()

	issue, _, err := client.Issues.Get(context.Background(), r.Owner, r.Name, issueNumber)
	if err != nil {
		return issue, fmt.Errorf("requesting issue %d in repo %s/%s: %v ", issueNumber, r.Owner, r.Name, err)
	}
	return issue, nil
}
