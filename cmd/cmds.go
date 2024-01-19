package cmd

import (
	"fmt"
	"os"

	"github.com/jirallreadyforthis/cli"
	"github.com/spf13/cobra"
)

func Make() (*cobra.Command, error) {
	root := &cobra.Command{
		Use:   "jirallreadyforthis",
		Short: "A cli tool for working with Jira issues",
		Long:  ``, // TODO
	}

	root.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "Lists Jira issues based on flag inputs",
		Long:  ``, // TODO
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("list called")

			f := GetFlags()
			l := cli.List{
				JiraToken:    f.JiraToken,
				JiraUrl:      f.JiraUrl,
				CustomFields: f.CustomFields,
				UserName:     f.UserName,
				Jql:          f.Jql,
				Linked:       f.Linked,
				GHToken:      f.GHToken,
				NotCommented: f.NotCommented,
			}
			err := l.ListJiraTickets()
			if err != nil {
				fmt.Printf("Error listing jira tickets: %v\n", err)
				os.Exit(1)
			}
		},
	})

	root.AddCommand(&cobra.Command{
		Use:   "set-status",
		Short: "Change the status on issues",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("setStatus called")

			f := GetFlags()

			s := cli.SetStatus{
				JiraToken:   f.JiraToken,
				JiraUrl:     f.JiraUrl,
				UserName:    f.UserName,
				Jql:         f.Jql,
				GHToken:     f.GHToken,
				IssueKeys:   f.IssueKeys,
				DryRun:      f.DryRun,
				Transitions: f.Transitions,
			}
			err := s.SetStatus()
			if err != nil {
				fmt.Printf("error setting issue statuses: %v\n\n", err)
				os.Exit(1)
			}
		},
	})

	root.AddCommand(&cobra.Command{
		Use:   "sprint-add",
		Short: "Add issues to a sprint",
		Long:  `Add issues to a sprint based on input issue keys (eg 'IPL-000') or issues found with an input jql query`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("sprint-add called")

			f := GetFlags()
			s := cli.SprintAdd{
				JiraToken: f.JiraToken,
				JiraUrl:   f.JiraUrl,
				UserName:  f.UserName,
				Jql:       f.Jql,
				GHToken:   f.GHToken,
				SprintId:  f.SprintId,
				IssueKeys: f.IssueKeys,
				DryRun:    f.DryRun,
			}

			err := s.AddIssuesToSprint()
			if err != nil {
				fmt.Printf("error adding issues to sprint: %v\n\n", err)
				os.Exit(1)
			}
		},
	})

	if err := configureFlags(root); err != nil {
		return nil, fmt.Errorf("unable to configure flags: %w", err)
	}

	return root, nil
}
