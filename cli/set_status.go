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
	DryRun      bool
	Transitions []string
	Debug       bool
	CheckLog    bool
}

func (s SetStatus) SetStatus() error {
	p := jira.Project{
		Token:    s.JiraToken,
		UserName: s.UserName,
		JiraUrl:  s.JiraUrl,
	}

	count := 0
	if len(s.IssueKeys) > 0 {
		for _, issueKey := range s.IssueKeys {
			issue, err := getIssueFromKey(issueKey, p)
			if err != nil {
				return err
			}
			if s.DryRun {
				fmt.Printf("setting issue (key %s id %s) to new status\n", issueKey, issue.ID)
			} else {
				if s.CheckLog {
					transitions := strings.Split(s.Transitions[0], ",")
					// status to transition to
					status := transitions[len(transitions)-1]

					// if the status has recently changed from the status we are aiming to transition to we should avoid reverting this back
					if issueIsRecentlyTransitioned(issue.ID, status, p) {
						continue
					}
				}
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
				fmt.Printf("setting issue (key %s id %s) to new status\n", issue.Key, issue.ID)
			} else {
				if s.CheckLog {
					transitions := strings.Split(s.Transitions[0], ",")
					// status to transition to
					status := transitions[len(transitions)-1]

					// if the status has recently changed from the status we are aiming to transition to we should avoid reverting this back
					if issueIsRecentlyTransitioned(issue.ID, status, p) {
						continue
					}
				}
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

func (s SetStatus) transitionIssue(issue j.Issue, p jira.Project) error {
	if s.Debug {
		fmt.Printf("attempting to transition status on issue %s\n", issue.Key)
	}

	foundWorkflow := false
	currentStatus := strings.ToLower(issue.Fields.Status.Name)
	originalStatus := currentStatus
	transitioned := false

	for _, transition := range s.Transitions {
		workflow := strings.Split(transition, ";")
		for i, status := range workflow {
			// find where the issue is in the chain and keep transitioning to the next status until we get to the end of the workflow
			if currentStatus == strings.ToLower(status) {
				foundWorkflow = true
				if len(workflow) >= i+2 {
					transitioned = false

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

	if transitioned {
		c.Info.Sprintf("Transitioned issue %s from %s to %s", issue.Key, originalStatus, currentStatus)
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

func issueIsRecentlyTransitioned(issueId string, status string, p jira.Project) bool {

	issue, err := p.GetIssueWithChangeLog(issueId)
	if err != nil {
		c.Errorf("retrieving issueId %s with changelog", issueId)
		return false
	}

	if issue != nil {
		if changelog := issue.Changelog; changelog != nil {
			if histories := changelog.Histories; histories != nil {
				for _, history := range histories {
					if items := history.Items; items != nil {
						// check only the most recent changelog entry
						if len(items) > 0 {
							if items[0].Field == "status" {
								if strings.ToLower(items[0].FromString) == strings.ToLower(status) {
									c.Warn.Sprintf("NOT updating issue %s as it was updated from status %q to %q on %s", issue.Key, items[0].FromString, items[0].ToString, history.Created)
									return true
								}
							}
						}
					}
					return false
				}
			}
		}
	}
	return false
}
