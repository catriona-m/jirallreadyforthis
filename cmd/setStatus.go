package cmd

import (
	"fmt"
	"os"

	"github.com/jirallreadyforthis/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// setStatusCmd represents the setStatus command
var setStatusCmd = &cobra.Command{
	Use:   "set-status",
	Short: "Change the status on issues",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("setStatus called")

		// env vars
		viper.AutomaticEnv()
		jiraToken := viper.GetString("JIRA_TOKEN")
		if jiraToken == "" {
			fmt.Println("Missing required environment variable `JIRA_TOKEN`")
			os.Exit(1)
		}

		ghToken := viper.GetString("GITHUB_TOKEN")
		if ghToken == "" {
			fmt.Println("Missing required environment variable `GITHUB_TOKEN`")
			os.Exit(1)
		}

		jiraUrl := viper.GetString("JIRA_URL")
		userName := viper.GetString("JIRA_USER")

		jiraUrl, _ = cmd.Flags().GetString("jira-url")
		userName, _ = cmd.Flags().GetString("user-name")
		jql, _ := cmd.Flags().GetString("jql")
		issueKeys, _ := cmd.Flags().GetStringSlice("issue-keys")
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		transitions, _ := cmd.Flags().GetStringSlice("transitions")

		if jiraUrl == "" {
			fmt.Println("Missing required variable - please set environment variable `JIRA_URL or --jira-url`")
			os.Exit(1)
		}

		if userName == "" {
			fmt.Println("Missing required variable - please set environment variable `JIRA_USER or --user-name`")
			os.Exit(1)
		}

		s := cli.SetStatus{
			JiraToken:   jiraToken,
			JiraUrl:     jiraUrl,
			UserName:    userName,
			Jql:         jql,
			GHToken:     ghToken,
			IssueKeys:   issueKeys,
			DryRun:      dryRun,
			Transitions: transitions,
		}
		err := s.SetStatus()
		if err != nil {
			fmt.Printf("error setting issue statuses: %v\n\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(setStatusCmd)

	setStatusCmd.Flags().StringP("jira-url", "j", "", "The base jira url eg 'https://readyforthis.atlassian.net/'")
	setStatusCmd.Flags().StringP("user-name", "u", "", "User name associated with the jira token")
	setStatusCmd.Flags().StringP("jql", "", "", "Jql query string to filter issues to update the status of")
	setStatusCmd.Flags().StringSliceP("issue-keys", "", []string{}, "List of issue keys to update the status of")
	setStatusCmd.Flags().BoolP("dry-run", "", true, "Print a simulation of what is expected without making actual changes. Defaults to true.")
	setStatusCmd.Flags().StringSliceP("transitions", "", []string{}, "List of transition workflows in order based on status names eg 'to do;in progress;done,blocked;in progress;done")

	setStatusCmd.MarkFlagRequired("transitions")
	setStatusCmd.MarkFlagsMutuallyExclusive("issue-keys", "jql")
}
