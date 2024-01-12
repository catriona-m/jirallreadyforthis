package cli

import (
	"fmt"
	"strings"

	j "github.com/andygrunwald/go-jira"
	c "github.com/gookit/color"
	"github.com/jirallreadyforthis/lib/jira"
)

type SetStatus struct {
	JiraToken   string
	JiraUrl     string
	UserName    string
	Jql         string
	GHToken     string
	IssueKeys   []string
	StatusId    string
	DryRun      bool
	Transitions []string
	Debug       bool
}

func (s SetStatus) SetStatus() error {
	p := jira.Project{
		Token:    s.JiraToken,
		UserName: s.UserName,
		JiraUrl:  s.JiraUrl,
	}

	statuses, err := statusIdToNameMap(p)
	if err != nil {
		return err
	}
	_, ok := statuses[s.StatusId]
	if !ok {
		return fmt.Errorf("id %s is not a valid status id", s.StatusId)
	}

	count := 0
	if len(s.IssueKeys) > 0 {
		for _, issueKey := range s.IssueKeys {
			issue, err := getIssueFromKey(issueKey, p)
			if err != nil {
				return err
			}
			if s.DryRun {
				fmt.Printf("setting issue (key %s id %s) to status ( %s)\n", issueKey, issue.ID, s.StatusId)
			} else {
				err = s.transitionIssue(*issue, p)
				if err != nil {
					return err
				}
			}
			count++
		}
	} else if s.Jql != "" {
		issues, err := p.ListIssues(s.Jql)
		if err != nil {
			return err
		}

		for _, issue := range issues {
			if s.DryRun {
				fmt.Printf("setting issue (key %s id %s) to status (%s)\n", issue.Key, issue.ID, s.StatusId)
			} else {
				err = s.transitionIssue(issue, p)
				if err != nil {
					return err
				}
			}
			count++
		}
	}

	c.Info.Printf("\n Finished updating the status on %d issues\n", count)

	return nil
}

func statusIdToNameMap(p jira.Project) (map[string]string, error) {
	statusIdToName := make(map[string]string, 0)

	statuses, err := p.ListStatuses()
	if err != nil {
		return nil, err
	}

	for _, status := range statuses {
		statusIdToName[status.ID] = status.Name
	}
	return statusIdToName, nil
}

func (s SetStatus) transitionIssue(issue j.Issue, p jira.Project) error {
	if s.Debug {
		fmt.Printf("attempting to transition status on issue %s\n", issue.Key)
	}

	foundWorkflow := false
	currentStatus := strings.ToLower(issue.Fields.Status.Name)
	for _, transition := range s.Transitions {

		workflow := strings.Split(transition, ";")
		for i, status := range workflow {
			// find where the issue is in the chain and keep transitioning to the next status until we get to the end of the workflow
			if currentStatus == strings.ToLower(status) {
				foundWorkflow = true
				if len(workflow) >= i+2 {
					transitioned := false

					// get a list of status transitions that are currently possible for this issue and check them
					// against the next transition name in the input list
					possibleTransitions, err := p.GetPossibleIssueTransitions(issue.ID)
					if err != nil {
						return err
					}
					for _, pt := range possibleTransitions {
						if s.Debug {
							fmt.Printf("checking possible transition '%s' against input transition '%s'\n", pt.Name, workflow[i+1])
						}

						if strings.ToLower(pt.Name) == strings.ToLower(workflow[i+1]) {
							if s.Debug {
								fmt.Printf("transitioning %s from %s to status %s\n", issue.Key, currentStatus, pt.Name)
							}
							err := p.TransitionIssueStatus(issue.ID, pt.ID)
							if err != nil {
								return err
							}
							currentStatus = strings.ToLower(pt.Name)
							transitioned = true
							break
						}
					}
					if !transitioned {
						c.Warn.Printf("it was not possible to transition status '%s' to '%s'\n", currentStatus, workflow[i+1])
						c.Warn.Printf("possible transitions are:\n")
						for _, pt := range possibleTransitions {
							c.Warn.Printf("%s\n", pt.Name)
						}
					}
				}
			}
		}
		if foundWorkflow {
			break
		}
	}

	return nil
}

func getIssueFromKey(key string, p jira.Project) (*j.Issue, error) {
	jql := fmt.Sprintf("issueKey = %s", key)
	issues, err := p.ListIssues(jql)
	if err != nil {
		return nil, err
	}
	if len(issues) != 1 {
		return nil, fmt.Errorf("getting issue from key %s", key)
	}

	return &issues[0], nil
}
