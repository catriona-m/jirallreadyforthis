package gh

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v52/github"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"golang.org/x/oauth2"
)

type Token struct {
	Token *string
}

type Repo struct {
	Owner string
	Name  string
	Token
}

func NewRepo(repo, token string) Repo {
	repoName := strings.Split(repo, "/")
	r := Repo{
		Owner: repoName[0],
		Name:  repoName[1],
		Token: Token{
			Token: nil,
		},
	}

	if token != "" {
		r.Token.Token = &token
	}

	return r
}

func (t Token) NewClient() *github.Client {

	userCacheDir, _ := os.UserCacheDir()
	cache := diskcache.New(filepath.Join(userCacheDir, "autoReviewCache"))

	tc := &http.Client{
		Transport: &oauth2.Transport{
			Base: httpcache.NewTransport(cache),
		},
	}

	if t.Token != nil {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: *t.Token},
		)
		tc = &http.Client{
			Transport: &oauth2.Transport{
				Base:   httpcache.NewTransport(cache),
				Source: ts,
			},
		}
	}

	return github.NewClient(tc)
}
