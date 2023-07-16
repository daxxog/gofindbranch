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
	"strings"

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
	openPRsOnly := flag.Bool("open-prs-only", false, "Filter for only open pull requests")
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

		// Display the filtered branch information
		for _, branch := range branches {
			printBranch(repo, branch, client, *openPRsOnly, owner, repoName)
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

func extractPullRequestNumber(url string) int {
	parts := strings.Split(url, "/")
	number := parts[len(parts)-1]

	return parseInt(number)
}

func parseInt(number string) int {
	var result int
	_, err := fmt.Sscan(number, &result)
	if err != nil {
		return -1
	}

	return result
}

func printBranch(repo string, branch *github.Branch, client *github.Client, openPRsOnly bool, owner, repoName string) {
	branchInfo := struct {
		Repository string `json:"repository"`
		Branch     string `json:"branch"`
		PR         string `json:"pr"`
	}{
		Repository: fmt.Sprintf("%s/%s", repo, branch.GetCommit().GetCommit().GetAuthor().GetLogin()),
		Branch:     branch.GetName(),
	}

	// Retrieve pull request information
	prs, _, err := client.PullRequests.List(context.Background(), owner, repoName, &github.PullRequestListOptions{})
	if err != nil {
		log.Printf("Failed to retrieve pull requests for branch %s: %v", branch.GetName(), err)
	}

	if openPRsOnly {
		// Check if there is an open pull request
		for _, pr := range prs {
			if pr.GetState() == "open" {
				branchInfo.PR = pr.GetHTMLURL()
				break
			}
		}
	} else if len(prs) > 0 {
		branchInfo.PR = prs[0].GetHTMLURL()
	}

	jsonData, err := json.Marshal(branchInfo)
	if err != nil {
		log.Printf("Failed to marshal JSON: %v", err)
		return
	}

	fmt.Println(string(jsonData))
}
