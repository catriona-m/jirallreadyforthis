package jira

import (
	"fmt"

	j "github.com/andygrunwald/go-jira"
)

func (p Project) ListStatuses() ([]j.Status, error) {
	client, err := p.NewClient()
	if err != nil {
		return nil, fmt.Errorf("creating jira client: %v: ", err)
	}
	statuses, _, err := client.Status.GetAllStatuses()
	if err != nil {
		return nil, fmt.Errorf("getting jira statuses: %v", err)
	}

	return statuses, nil
}
