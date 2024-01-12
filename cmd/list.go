package cmd

import (
	"fmt"
	"os"

	"github.com/jirallreadyforthis/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists Jira issues based on flag inputs",
	Long:  ``, // TODO
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list called")

		jiraUrl, _ := cmd.Flags().GetString("jira-url")
		customFields, _ := cmd.Flags().GetStringSlice("custom-fields")
		userName, _ := cmd.Flags().GetString("user-name")
		jql, _ := cmd.Flags().GetString("jql")
		lastCommented, _ := cmd.Flags().GetInt("not-commented")
		linked, _ := cmd.Flags().GetBool("linked")

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

		l := cli.List{
			JiraToken:    jiraToken,
			JiraUrl:      jiraUrl,
			CustomFields: customFields,
			UserName:     userName,
			Jql:          jql,
			Linked:       linked,
			GHToken:      ghToken,
			NotCommented: lastCommented,
		}
		err := l.ListJiraTickets()
		if err != nil {
			fmt.Sprintf("Error listing jira tickets: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringP("jira-url", "j", "", "The base jira url eg 'https://readyforthis.atlassian.net/' HI")
	listCmd.Flags().StringSliceP("custom-fields", "f", []string{}, "A list of custom fields to search for links in")
	listCmd.Flags().StringP("user-name", "u", "", "User name associated with the jira token")
	listCmd.Flags().StringP("jql", "", "", "Jql query string to filter issues")
	listCmd.Flags().IntP("not-commented", "", 0, "Filter issues based on whether they have been commented on in a specified number of days.")
	listCmd.Flags().BoolP("linked", "", true, "Only list jira issues with either github issues that are closed or pull requests that are merged. Defaults to true.")
}
