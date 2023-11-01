package cli

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	c "github.com/gookit/color"
	"github.com/jirallreadyforthis/lib/gh"
	"github.com/jirallreadyforthis/lib/jira"
)

type List struct {
	JiraToken    string
	JiraUrl      string
	UserName     string
	Jql          string
	CustomFields []string
	Linked       bool
	GHToken      string
}

func (l List) ListJiraTickets() error {

	p := jira.Project{
		Token:    l.JiraToken,
		UserName: l.UserName,
		JiraUrl:  l.JiraUrl,
	}

	issues, err := p.ListIssues(l.Jql)
	if err != nil {
		return err
	}

	count := 0
	for _, issue := range issues {
		if l.Linked {
			githubLinks := make([]string, 0)
			if len(l.CustomFields) > 0 {
				for _, field := range l.CustomFields {
					if issue.Fields.Unknowns != nil {
						fieldValue, exists := issue.Fields.Unknowns.Value(field)
						if exists && fieldValue != nil {
							githubLinks = append(githubLinks, findGithubLinks(fieldValue.(string))...)
						}
					}
				}
			}
			githubLinks = append(githubLinks, findGithubLinks(issue.Fields.Description)...)

			// search issue comments for links
			issueWithComments, err := p.GetIssue(issue.ID)
			if err != nil {
				return err
			}

			if issueWithComments != nil {
				if len(issueWithComments.Fields.Comments.Comments) > 0 {
					for _, comment := range issueWithComments.Fields.Comments.Comments {
						githubLinks = append(githubLinks, findGithubLinks(comment.Body)...)
					}
				}
			}

			ghClosedOrMerged := make([]string, 0)
			githubLinks = removeDuplicates(githubLinks)
			for _, link := range githubLinks {
				if s := l.closedOrMerged(link); s != "" {
					ghClosedOrMerged = append(ghClosedOrMerged, s)
				}
			}
			if len(ghClosedOrMerged) > 0 {
				count++
				c.Printf("\n\n<green>%s</>	%s\n", issue.Fields.Summary, l.getJiraHtmlUrl(issue.Key))
				c.Printf("\t%s", strings.Join(ghClosedOrMerged, "\t\n\t"))
			}
		} else {
			count++
			c.Printf("\n\n<green>%s</>	%s\n", issue.Fields.Summary, l.getJiraHtmlUrl(issue.Key))
		}
	}
	c.Info.Printf("Finished listing %d issues\n", count)
	return nil
}

func (l List) closedOrMerged(link string) string {
	closedOrMergedString := ""

	re := regexp.MustCompile("https://github\\.com/(?P<repoName>\\S+/\\S+)/(?:pull|issues)/(?P<number>\\d+)")
	matches := re.FindAllStringSubmatch(link, -1)

	repoName := ""
	number := ""
	if len(matches) > 0 {
		repoIndex := re.SubexpIndex("repoName")
		repoName = matches[0][repoIndex]
		numberIndex := re.SubexpIndex("number")
		number = matches[0][numberIndex]
	}

	if repoName != "" && number != "" {

		repo := gh.NewRepo(repoName, l.GHToken)
		i, _ := strconv.Atoi(number)

		if strings.Contains(link, "issues") {
			issue, err := repo.GetIssue(i)
			if err != nil {
				fmt.Printf("getting issue %s: %v\n", link, err)
				return closedOrMergedString
			}

			if issue != nil {
				if issue.IsPullRequest() {
					merged, err := repo.PullRequestIsMerged(i)
					if err != nil {
						fmt.Printf("Error checking if pr is merged %s: %v\n", link, err)
						return closedOrMergedString
					}

					if merged {
						closedOrMergedString = c.Sprintf("<lightMagenta>%s\t%s</>", issue.GetTitle(), link)
					}

				} else if issue.GetState() == "Closed" {
					closedOrMergedString = c.Sprintf("<lightRed>%s\t%s</>", issue.GetTitle(), link)
				}
			}

		} else if strings.Contains(link, "pull") {
			merged, err := repo.PullRequestIsMerged(i)
			if err != nil {
				c.Errorf("Error checking if pr is merged %s: %v\n", link, err)
				return closedOrMergedString
			}

			pr, err := repo.GetPullRequest(i)
			if err != nil {
				c.Errorf("Error getting pr %s: %v", link, err)
				return closedOrMergedString
			}

			if merged {
				closedOrMergedString = c.Sprintf("<lightMagenta>%s\t%s</>", pr.GetTitle(), link)
			}
		}
	}
	return closedOrMergedString
}

func (l List) getJiraHtmlUrl(issueKey string) string {
	return fmt.Sprintf("%s/browse/%s", l.JiraUrl, issueKey)
}

func findGithubLinks(text string) []string {
	re := regexp.MustCompile("https://github\\.com/\\S+/\\S+/(?:pull|issues)/\\d+")
	matches := re.FindAllString(text, -1)

	links := make([]string, 0)
	for _, match := range matches {
		match = strings.Split(match, "|")[0]
		links = append(links, match)
	}

	return links
}

func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	list := make([]string, 0)

	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
