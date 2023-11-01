package jira

import (
	j "github.com/andygrunwald/go-jira"
)

type Project struct {
	Token    string
	UserName string
	JiraUrl  string
}

func (p Project) NewClient() (*j.Client, error) {

	tp := j.BasicAuthTransport{
		Username: p.UserName,
		Password: p.Token,
	}

	client, err := j.NewClient(tp.Client(), p.JiraUrl)
	if err != nil {
		return nil, err
	}

	return client, nil
}
