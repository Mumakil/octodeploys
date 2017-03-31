package main

import (
	"flag"
	"fmt"
	"sync"
)

var limit uint64
var environment string
var state string
var includeStatuses bool

func init() {
	flag.Uint64Var(&limit, "limit", 100, "how many deployments to list")
	flag.StringVar(&environment, "environment", "", "filter by given environment")
	flag.StringVar(&state, "state", "", "filter by given state")
	flag.BoolVar(&includeStatuses, "statuses", false, "fetch also statuses")
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

	if includeStatuses || state != "" {
		err = fetchStatuses(client, deployments)
		if err != nil {
			return err
		}
	}
	if state != "" {
		deployments = deployments.FilterByState(state)
	}

	fmt.Println(deployments.String(includeStatuses))

	return nil
}

func fetchDeployments(client *Client) (Deployments, error) {
	var deployments Deployments
	path := fmt.Sprintf("/repos/%s/deployments", GitHubRepository)
	query := map[string]string{
		"limit":       fmt.Sprintf("%d", limit),
		"environment": environment,
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
