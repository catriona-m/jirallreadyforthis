package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type FlagData struct {
	JiraToken    string
	JiraUrl      string
	UserName     string
	Jql          string
	GHToken      string
	IssueKeys    []string
	DryRun       bool
	Transitions  []string
	Debug        bool
	SprintId     int
	CustomFields []string
	Linked       bool
	NotCommented int
}

func configureFlags(root *cobra.Command) error {
	flags := FlagData{}
	pflags := root.PersistentFlags()

	pflags.StringVarP(&flags.JiraUrl, "jira-url", "j", "", "The base jira url eg 'https://readyforthis.atlassian.net/'")
	pflags.StringVarP(&flags.UserName, "jira-user", "u", "", "User name associated with the jira token")
	pflags.StringVarP(&flags.UserName, "token-jira", "", "", "Jira API token")
	pflags.StringVarP(&flags.UserName, "token-gh", "", "", "Github API token")
	pflags.StringVarP(&flags.Jql, "jql", "", "", "Jql query string to filter issues on")
	pflags.BoolVarP(&flags.DryRun, "dry-run", "", true, "Print a simulation of what is expected without making actual changes. Defaults to true.")

	pflags.StringSliceVarP(&flags.IssueKeys, "issue-keys", "", []string{}, "List of issue keys to process")
	pflags.StringSliceVarP(&flags.Transitions, "transitions", "", []string{}, "List of transition workflows in order based on status names eg 'to do;in progress;done,blocked;in progress;done")

	pflags.IntVarP(&flags.SprintId, "sprint-id", "", 0, "The id of the sprint to move issues to")
	pflags.StringSliceVarP(&flags.CustomFields, "custom-fields", "f", []string{}, "A list of custom fields to search for links in")

	pflags.IntVarP(&flags.NotCommented, "not-commented", "", 0, "Filter issues based on whether they have been commented on in a specified number of days.")
	pflags.BoolVarP(&flags.Linked, "linked", "", true, "Only list jira issues with either github issues that are closed or pull requests that are merged. Defaults to true.")

	// binding map for viper/pflag -> env
	m := map[string]string{
		"jira-url":      "JIRA_URL",
		"jira-user":     "JIRA_USER",
		"token-jira":    "JIRA_TOKEN",
		"token-gh":      "GITHUB_TOKEN",
		"jql":           "",
		"dry-run":       "",
		"issue-keys":    "",
		"transitions":   "",
		"sprint-id":     "",
		"custom-fields": "",
		"not-commented": "",
		"linked":        "",
	}

	for name, env := range m {
		if err := viper.BindPFlag(name, pflags.Lookup(name)); err != nil {
			return fmt.Errorf("error binding '%s' flag: %w", name, err)
		}

		if env != "" {
			if err := viper.BindEnv(name, env); err != nil {
				return fmt.Errorf("error binding '%s' to env '%s' : %w", name, env, err)
			}
		}
	}
	return nil
}

func GetFlags() FlagData {
	return FlagData{
		JiraToken:    viper.GetString("token-jira"),
		JiraUrl:      viper.GetString("jira-url"),
		UserName:     viper.GetString("jira-user"),
		Jql:          viper.GetString("jql"),
		GHToken:      viper.GetString("token-gh"),
		IssueKeys:    viper.GetStringSlice("issue-keys"),
		DryRun:       viper.GetBool("dry-run"),
		Transitions:  viper.GetStringSlice("transitions"),
		Debug:        viper.GetBool("debug"),
		SprintId:     viper.GetInt(""),
		CustomFields: nil,
		Linked:       false,
		NotCommented: 0,
	}
}
