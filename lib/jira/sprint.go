package jira

import (
	"fmt"
)

func (p Project) AddToSprint(sprintId int, issueIds []string) error {
	client, err := p.NewClient()
	if err != nil {
		return fmt.Errorf("creating jira client: %v: ", err)
	}

	_, err = client.Sprint.MoveIssuesToSprint(sprintId, issueIds)
	if err != nil {
		return fmt.Errorf("moving issues to sprint id %d", sprintId)
	}
	return nil
}
