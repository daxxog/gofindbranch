package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

func authenticate() *github.Client {
	pat := os.Getenv("GH_PAT")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: pat},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client
}

func filterBranches(branches []*github.Branch, filterRegex string) []*github.Branch {
	filteredBranches := []*github.Branch{}
	filter := regexp.MustCompile(filterRegex)

	for _, branch := range branches {
		if filter.MatchString(branch.GetName()) {
			filteredBranches = append(filteredBranches, branch)
		}
	}

	return filteredBranches
}

func main() {
	// Parse command-line arguments
	filter := flag.String("filter", "", "Branch filter (regex)")
	flag.Parse()

	if *filter == "" {
		log.Fatal("Please provide a branch filter using the -filter flag")
	}

	// Read the list of repositories from standard input
	repoList := []string{}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		repoList = append(repoList, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("Failed to read input:", err)
	}

	// Authenticate with the GitHub API
	client := authenticate()

	// Retrieve branches for each repository
	for _, repo := range repoList {
		owner, repoName := parseRepo(repo)
		if owner == "" || repoName == "" {
			log.Printf("Invalid repository format: %s", repo)
			continue
		}

		branches, _, err := client.Repositories.ListBranches(context.Background(), owner, repoName, nil)
		if err != nil {
			log.Printf("Failed to retrieve branches for repository %s: %v", repo, err)
			continue
		}

		// Filter branches based on the filter input
		filteredBranches := filterBranches(branches, *filter)

		// Display the filtered branch information
		for _, branch := range filteredBranches {
			printBranch(repo, branch)
		}
	}
}

func parseRepo(repo string) (string, string) {
	// Assuming the format is "<owner>/<repository>"
	split := regexp.MustCompile(`\s*/\s*`).Split(repo, -1)
	if len(split) != 2 {
		return "", ""
	}

	return split[0], split[1]
}

func printBranch(repo string, branch *github.Branch) {
	branchInfo := struct {
		Repository string `json:"repository"`
		Branch     string `json:"branch"`
	}{
		Repository: fmt.Sprintf("%s/%s", repo, branch.GetCommit().GetCommit().GetAuthor().GetLogin()),
		Branch:     branch.GetName(),
	}

	jsonData, err := json.Marshal(branchInfo)
	if err != nil {
		log.Printf("Failed to marshal JSON: %v", err)
		return
	}

	fmt.Println(string(jsonData))
}
