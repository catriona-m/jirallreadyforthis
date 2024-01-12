/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/jirallreadyforthis/cli"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// sprintAddCmd represents the sprintAdd command
var sprintAddCmd = &cobra.Command{
	Use:   "sprint-add",
	Short: "Add issues to a sprint",
	Long:  `Add issues to a sprint based on input issue keys (eg 'IPL-000') or issues found with an input jql query`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sprint-add called")

		jiraUrl, _ := cmd.Flags().GetString("jira-url")
		userName, _ := cmd.Flags().GetString("user-name")
		jql, _ := cmd.Flags().GetString("jql")
		sprintId, _ := cmd.Flags().GetInt("sprint-id")
		issueKeys, _ := cmd.Flags().GetStringSlice("issue-keys")
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		// env vars
		viper.AutomaticEnv()
		jiraToken := viper.GetString("JRFT_JIRA_TOKEN")
		if jiraToken == "" {
			fmt.Println("Missing required environment variable `JRFT_JIRA_TOKEN`")
			os.Exit(1)
		}

		ghToken := viper.GetString("JRFT_GH_TOKEN")
		if ghToken == "" {
			fmt.Println("Missing required environment variable `JRFT_GH_TOKEN`")
			os.Exit(1)
		}

		s := cli.SprintAdd{
			JiraToken: jiraToken,
			JiraUrl:   jiraUrl,
			UserName:  userName,
			Jql:       jql,
			GHToken:   ghToken,
			SprintId:  sprintId,
			IssueKeys: issueKeys,
			DryRun:    dryRun,
		}
		err := s.AddIssuesToSprint()
		if err != nil {
			fmt.Sprintf("error adding issues to sprint: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(sprintAddCmd)

	sprintAddCmd.Flags().StringP("jira-url", "j", "", "The base jira url eg 'https://readyforthis.atlassian.net/'")
	sprintAddCmd.Flags().StringP("user-name", "u", "", "User name associated with the jira token")
	sprintAddCmd.Flags().StringP("jql", "", "", "Jql query string to filter issues to add to the sprint")
	sprintAddCmd.Flags().IntP("sprint-id", "", 0, "The id of the sprint to move issues to")
	sprintAddCmd.Flags().StringSliceP("issue-keys", "", []string{}, "List of issue keys to add to the sprint")
	sprintAddCmd.Flags().BoolP("dry-run", "", true, "Print a simulation of what is expected without making actual changes")

	sprintAddCmd.MarkFlagRequired("sprint-id")
	sprintAddCmd.MarkFlagsMutuallyExclusive("issue-keys", "jql")
}
