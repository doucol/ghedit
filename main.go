package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-github/v26/github"
	"golang.org/x/oauth2"
)

func main() {
	token := os.Getenv("GITHUB_ACCESS_TOKEN")
	if token == "" {
		panic("'GITHUB_ACCESS_TOKEN' env variable is required")
	}
	var orgs, cmd string
	flag.StringVar(&orgs, "orgs", "", "the org list may be comma separated, if you like")
	flag.StringVar(&cmd, "cmd", "", "supported cmds: wiki-off")
	flag.Parse()

	if orgs == "" || cmd == "" {
		fmt.Printf("no orgs or commands to execute")
		os.Exit(0)
	}

	ctx := context.Background()
	tc := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))
	gh := github.NewClient(tc)

	// list all repositories for the authenticated user
	options := github.RepositoryListByOrgOptions{
		Type:        "all",
		ListOptions: github.ListOptions{PerPage: 500},
	}

	// need to support paging at some point

	for _, org := range strings.Split(orgs, ",") {
		rl, _, err := gh.Repositories.ListByOrg(ctx, org, &options)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		} else if rl == nil {
			fmt.Printf("no repose exist for %v\n", org)
			continue
		}
		for _, r := range rl {
			fmt.Printf("Repo name: %v\n", *r.Name)
			switch cmd {
			case "wiki-off":
				*r.HasWiki = false
				r2, _, err := gh.Repositories.Edit(ctx, org, *r.Name, r)
				if err != nil {
					er := err.(*github.ErrorResponse)
					fmt.Printf("error: %v\n\n", er.Message)
				} else {
					fmt.Printf(fmt.Sprintf("Repo Wiki enabled: %v\n\n", strconv.FormatBool(*r2.HasWiki)))
				}
			default:
			}
		}
	}
}
