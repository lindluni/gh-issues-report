package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/repository"
	"github.com/google/go-github/v42/github"
)

func main() {
	repo, err := gh.CurrentRepository()
	if err != nil {
		fmt.Printf("Unable to determin current repository: %s\n", err.Error())
		os.Exit(1)
	}

	httpClient, err := gh.HTTPClient(nil)
	if err != nil {
		fmt.Printf("Unable to create GitHub client: %s\n", err.Error())
		os.Exit(1)
	}

	ctx := context.Background()
	client := github.NewClient(httpClient)

	since := time.Now()
	since = time.Date(since.Year(), since.Month(), 1, 0, 0, 0, 0, time.UTC)
	retrieveIssueStatistics(ctx, since, repo, client)

}

func retrieveIssueStatistics(ctx context.Context, since time.Time, repo repository.Repository, client *github.Client) {
	var issues []*github.Issue
	opts := &github.IssueListByRepoOptions{
		State:       "all",
		Since:       since,
		ListOptions: github.ListOptions{PerPage: 100},
	}
	fmt.Printf("Retrieving issue statistics for %s/%s...\n", repo.Owner(), repo.Name())
	for {
		list, resp, err := client.Issues.ListByRepo(ctx, repo.Owner(), repo.Name(), opts)
		if err != nil {
			fmt.Printf("Unable to list repository issues: %s\n", err.Error())
			os.Exit(1)
		}
		issues = append(issues, list...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	var canonicalIssues []*github.Issue
	var canonicalPullRequests []*github.Issue
	for _, issue := range issues {
		if issue.IsPullRequest() {
			canonicalPullRequests = append(canonicalPullRequests, issue)
			continue
		}
		canonicalIssues = append(canonicalIssues, issue)
	}

	closedIssues := 0
	for _, issue := range canonicalIssues {
		if issue.GetState() == "closed" {
			closedIssues++
		}
	}

	closedPullRequests := 0
	for _, issue := range canonicalPullRequests {
		if issue.GetState() == "closed" {
			closedPullRequests++
		}
	}

	fmt.Printf("Total issues opened since %s: %d\n", since.Format("2006-01-02"), len(canonicalIssues))
	fmt.Printf("Total issues closed since %s: %d\n", since.Format("2006-01-02"), closedIssues)

	fmt.Printf("Total pull requests opened since %s: %d\n", since.Format("2006-01-02"), len(canonicalPullRequests))
	fmt.Printf("Total pull requests closed since %s: %d\n", since.Format("2006-01-02"), closedPullRequests)
}
