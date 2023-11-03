package jira

import (
	"fmt"

	j "github.com/andygrunwald/go-jira"
)

func (p Project) ListIssues(jql string) ([]j.Issue, error) {

	client, err := p.NewClient()
	if err != nil {
		return nil, fmt.Errorf("creating jira client: %v: ", err)
	}

	last := 0
	var issues []j.Issue
	for {
		opt := &j.SearchOptions{
			MaxResults: 1000,
			StartAt:    last,
		}

		chunk, resp, err := client.Issue.Search(jql, opt)
		if err != nil {
			return nil, err
		}

		total := resp.Total
		if issues == nil {
			issues = make([]j.Issue, 0, total)
		}
		issues = append(issues, chunk...)
		last = resp.StartAt + len(chunk)
		if last >= total {
			return issues, nil
		}
	}

	return issues, nil
}

func (p Project) GetIssue(issueId string) (*j.Issue, error) {
	client, err := p.NewClient()
	if err != nil {
		return nil, fmt.Errorf("creating jira client: %v: ", err)
	}

	issue, _, err := client.Issue.Get(issueId, nil)
	if err != nil {
		return nil, fmt.Errorf("getting jira issue %s: %v: ", issueId, err)
	}

	return issue, nil
}
