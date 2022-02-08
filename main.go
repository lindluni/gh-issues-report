package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cli/go-gh"
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
		return
	}

	ctx := context.Background()
	client := github.NewClient(httpClient)

	var workflowRuns []*github.WorkflowRun
	opts := &github.ListWorkflowRunsOptions{
		Created:     fmt.Sprintf(">%d-%d-01", time.Now().Year(), time.Now().Month()+1),
		ListOptions: github.ListOptions{PerPage: 100},
	}
	for {
		runs, resp, err := client.Actions.ListRepositoryWorkflowRuns(ctx, repo.Owner(), repo.Name(), opts)
		if err != nil {
			fmt.Printf("Unable to list repository workflow runs: %s\n", err.Error())
			os.Exit(1)
		}
		workflowRuns = append(workflowRuns, runs.WorkflowRuns...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	fmt.Println("Found", len(workflowRuns), "workflow runs")
}

// For more examples of using go-gh, see:
// https://github.com/cli/go-gh/blob/trunk/example_gh_test.go
