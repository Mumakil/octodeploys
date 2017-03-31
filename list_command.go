package main

import (
	"flag"
	"fmt"
	"sync"
)

// Limit how many deployments to fetch from GitHub API
var Limit uint64

// Environment filter by environment
var Environment string

// State filter by state
var State string

// IncludeStatuses whether to include statuses when reporting
var IncludeStatuses bool

func init() {
	flag.Uint64Var(&Limit, "limit", 100, "how many deployments to list")
	flag.StringVar(&Environment, "environment", "", "filter by given environment")
	flag.StringVar(&State, "state", "", "filter by given state")
	flag.BoolVar(&IncludeStatuses, "statuses", false, "fetch also statuses")
}

func listCommand(args []string) error {
	err := validateGlobalArgs()
	if err != nil {
		return err
	}

	client := NewClient(GitHubToken)

	deployments, err := fetchDeployments(client)
	if err != nil {
		return err
	}

	if IncludeStatuses || State != "" {
		err = fetchStatuses(client, deployments)
		if err != nil {
			return err
		}
	}
	if State != "" {
		deployments = deployments.FilterByState(State)
	}

	fmt.Println(deployments.String(IncludeStatuses))

	return nil
}

func fetchDeployments(client *Client) (Deployments, error) {
	var deployments Deployments
	path := fmt.Sprintf("/repos/%s/deployments", GitHubRepository)
	query := map[string]string{
		"limit":       fmt.Sprintf("%d", Limit),
		"environment": Environment,
	}

	err := client.Get(path, query, &deployments)
	return deployments, err
}

func fetchStatuses(client *Client, deployments Deployments) error {
	wg := sync.WaitGroup{}
	errors := make(chan error, len(deployments))
	defer close(errors)
	for _, d := range deployments {
		wg.Add(1)
		go func(deployment *Deployment) {
			status, err := fetchLastStatus(client, deployment.ID)
			if err != nil {
				errors <- err
			}
			deployment.Status = status
			wg.Done()
		}(d)
	}
	wg.Wait()
	if len(errors) > 0 {
		return <-errors
	}
	return nil
}

func fetchLastStatus(client *Client, deploymentID uint64) (*Status, error) {
	var statuses []Status
	path := fmt.Sprintf("/repos/%s/deployments/%d/statuses", GitHubRepository, deploymentID)
	query := map[string]string{
		"limit": "1",
	}
	err := client.Get(path, query, &statuses)
	if err != nil {
		return nil, err
	}
	if len(statuses) == 0 {
		return nil, nil
	}
	return &statuses[0], nil
}
