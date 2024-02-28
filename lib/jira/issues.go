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

func (p Project) GetIssueWithChangeLog(issueId string) (*j.Issue, error) {
	client, err := p.NewClient()
	if err != nil {
		return nil, fmt.Errorf("creating jira client: %v: ", err)
	}

	opts := j.GetQueryOptions{
		Expand: "changelog",
	}

	issue, _, err := client.Issue.Get(issueId, &opts)
	if err != nil {
		return nil, fmt.Errorf("getting jira issue %s: %v: ", issueId, err)
	}

	return issue, nil
}

func (p Project) TransitionIssueStatus(issueId string, transitionID string) error {
	client, err := p.NewClient()
	if err != nil {
		return fmt.Errorf("creating jira client: %v: ", err)
	}

	_, err = client.Issue.DoTransition(issueId, transitionID)

	if err != nil {
		return fmt.Errorf("transitioning issue id %s to status id %s : %v: ", issueId, transitionID, err)
	}

	return nil
}

func (p Project) GetPossibleIssueTransitions(issueId string) ([]j.Transition, error) {
	client, err := p.NewClient()
	if err != nil {
		return nil, fmt.Errorf("creating jira client: %v: ", err)
	}

	transitions, _, err := client.Issue.GetTransitions(issueId)
	if err != nil {
		return nil, fmt.Errorf("getting possible issue transitions from id %s", issueId)
	}

	return transitions, nil
}
