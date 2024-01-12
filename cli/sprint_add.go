package cli

import (
	"fmt"

	c "github.com/gookit/color"
	"github.com/jirallreadyforthis/lib/jira"
)

type SprintAdd struct {
	JiraToken string
	JiraUrl   string
	UserName  string
	Jql       string
	GHToken   string
	IssueKeys []string
	SprintId  int
	DryRun    bool
}

func (s SprintAdd) AddIssuesToSprint() error {
	p := jira.Project{
		Token:    s.JiraToken,
		UserName: s.UserName,
		JiraUrl:  s.JiraUrl,
	}

	issueIds := make([]string, 0)
	count := 0
	if len(s.IssueKeys) > 0 {
		for _, issueKey := range s.IssueKeys {
			issueId, err := getIssueIdFromKey(issueKey, p)
			if err != nil {
				return err
			}
			issueIds = append(issueIds, issueId)
			if s.DryRun {
				fmt.Printf("adding issue (key %s id %s) to sprint (id %d)\n", issueKey, issueId, s.SprintId)
			}
			count++
		}
	} else if s.Jql != "" {
		issues, err := p.ListIssues(s.Jql)
		if err != nil {
			return err
		}

		for _, issue := range issues {
			issueIds = append(issueIds, issue.ID)
			if s.DryRun {
				fmt.Printf("adding issue (key %s id %s) to sprint (id %d)\n", issue.Key, issue.ID, s.SprintId)
			}
			count++
		}

	}

	if !s.DryRun {
		if len(issueIds) > 0 {
			err := p.AddToSprint(s.SprintId, issueIds)
			if err != nil {
				return err
			}
		}
	}

	c.Info.Printf("\n Finished adding %d issues to sprint %d\n", count, s.SprintId)

	return nil
}

func getIssueIdFromKey(key string, p jira.Project) (string, error) {
	jql := fmt.Sprintf("issueKey = %s", key)
	issues, err := p.ListIssues(jql)
	if err != nil {
		return "", err
	}
	if len(issues) != 1 || issues[0].ID == "" {
		return "", fmt.Errorf("getting issue from key %s", key)
	}

	return issues[0].ID, nil
}
