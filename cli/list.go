package cli

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v52/github"
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
	NotCommented int
	ClosedWithin int
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
		issueWithComments, err := p.GetIssue(issue.ID)
		if l.NotCommented > 0 {
			if len(issueWithComments.Fields.Comments.Comments) > 0 {
				lastComment := issueWithComments.Fields.Comments.Comments[len(issueWithComments.Fields.Comments.Comments)-1]
				date, err := time.Parse("2006-01-02", strings.Split(lastComment.Created, "T")[0])
				if err != nil {
					return fmt.Errorf("parsing comment creation time: %v", err)
				}

				if !date.Before(time.Now().AddDate(0, 0, -l.NotCommented)) {
					// found a comment after the specified time, so move to the next issue
					continue
				}
			}
		}

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
				ghIssue, repo := l.getIssueAndRepoFromLink(link)
				if ghIssue != nil && repo != nil {

					closed, err := closedOrMergedWithin(ghIssue, l.ClosedWithin)
					if err != nil {
						return err
					}
					if !closed {
						continue
					}
					s, err := l.closedOrMerged(ghIssue, *repo)
					if err != nil {
						return err
					}
					if s != "" {
						ghClosedOrMerged = append(ghClosedOrMerged, s)
					}
				}
			}

			createdTime := time.Time(issue.Fields.Created)
			date := strings.Split(createdTime.String(), " ")[0]
			ghClosedOrMerged = removeDuplicates(ghClosedOrMerged)
			if len(ghClosedOrMerged) > 0 {
				count++
				c.Printf("\n\n<green>%s\t%s\t%s</>\n", date, l.getJiraHtmlUrl(issue.Key), issue.Fields.Summary)
				c.Printf("\t%s", strings.Join(ghClosedOrMerged, "\t\n\t"))
			}
		} else {
			count++
			c.Printf("\n\n<green>%s</>	%s\n", l.getJiraHtmlUrl(issue.Key), issue.Fields.Summary)
		}
	}
	c.Info.Printf("\nFinished listing %d issues\n", count)
	return nil
}

func (l List) closedOrMerged(issue *github.Issue, repo gh.Repo) (string, error) {
	closedOrMergedString := ""

	if issue != nil {
		if issue.IsPullRequest() {
			merged, err := repo.PullRequestIsMerged(*issue.Number)
			if err != nil {
				c.Errorf("Error checking if pr %d is merged %s: %v\n", *issue.Number, err)
				return closedOrMergedString, nil
			}

			if merged {
				closedDate := strings.Split(issue.GetClosedAt().String(), " ")[0]
				if l.ClosedWithin > 0 {
					closed, err := closedOrMergedWithin(issue, l.ClosedWithin)
					if err != nil {
						return "", err
					}
					if !closed {
						return closedOrMergedString, nil
					}
				}
				closedOrMergedString = c.Sprintf("<lightMagenta>%s\t%s\t%s</>", closedDate, issue.GetHTMLURL(), issue.GetTitle())
			}
		} else if issue.GetState() == "closed" {
			closedDate := strings.Split(issue.GetClosedAt().String(), " ")[0]
			closedOrMergedString = c.Sprintf("<lightRed>%s\t%s\t%s</>", closedDate, issue.GetHTMLURL(), issue.GetTitle())
		}
	}

	return closedOrMergedString, nil
}

func (l List) getJiraHtmlUrl(issueKey string) string {
	return fmt.Sprintf("%s/browse/%s", l.JiraUrl, issueKey)
}

func findGithubLinks(text string) []string {
	re := regexp.MustCompile("https://github\\.com/[\\w-]+/[\\w-]+/(?:pull|issues)/\\d+")
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

func closedOrMergedWithin(issue *github.Issue, days int) (bool, error) {

	closedDate := strings.Split(issue.GetClosedAt().String(), " ")[0]
	parsedDate, err := time.Parse("2006-01-02", closedDate)
	if err != nil {
		return false, fmt.Errorf("parsing closed/merged time: %v", err)
	}

	if parsedDate.After(time.Now().AddDate(0, 0, -days)) {
		return true, nil
	}

	return false, nil
}

func (l List) getIssueAndRepoFromLink(link string) (*github.Issue, *gh.Repo) {
	re := regexp.MustCompile("https://github\\.com/(?P<repoName>[\\w-]+/[\\w-]+)/(?:pull|issues)/(?P<number>\\d+)")
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
				c.Errorf("\n Error getting issue from extracted link %s: %v\n", link, err)
				return nil, nil
			}
			return issue, &repo
		}
	}
	return nil, nil
}
