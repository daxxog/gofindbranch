package main

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/google/go-github/v39/github"
)

func filterRepositories(repos []*github.Repository, filterRegex string) []*github.Repository {
	filteredRepos := []*github.Repository{}
	filter := regexp.MustCompile(filterRegex)

	for _, repo := range repos {
		if filter.MatchString(repo.GetName()) {
			filteredRepos = append(filteredRepos, repo)
		}
	}

	return filteredRepos
}

func repos_main(filter string, currentUserOnly bool) {
	if filter == "" && !currentUserOnly {
		log.Fatal("Please provide a repository filter using the -filter flag or use -current-user-only flag to list repositories owned by the current user")
	}

	// Authenticate with the GitHub API
	client := authenticate()

	// Retrieve repositories
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	var allRepos []*github.Repository

	for {
		var repos []*github.Repository
		var resp *github.Response
		var err error

		if currentUserOnly {
			// Retrieve repositories owned by the current user
			user, _, err := client.Users.Get(context.Background(), "")
			if err != nil {
				log.Fatal("Failed to retrieve current user:", err)
			}
			repos, resp, err = client.Repositories.List(context.Background(), *user.Login, opt)
		} else {
			// Retrieve all repositories
			repos, resp, err = client.Repositories.List(context.Background(), "", opt)
		}

		if err != nil {
			log.Fatal("Failed to retrieve repositories:", err)
		}

		allRepos = append(allRepos, repos...)

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	// Filter repositories based on the filter input
	filteredRepos := filterRepositories(allRepos, filter)

	// Display the filtered repositories
	for _, repo := range filteredRepos {
		fmt.Println(repo.GetName())
	}
}
