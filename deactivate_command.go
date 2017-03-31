package main

import (
	"flag"
	"fmt"
	"sync"
)

// BetaAccessHeader required to use inactive statuses in GitHub deployments api
const BetaAccessHeader = "application/vnd.github.ant-man-preview+json"

var description string

func init() {
	flag.StringVar(&description, "description", "", "description for the inactive status")
}

func deactivateCommand(deploymentIDs []string) error {
	err := validateGlobalArgs()
	if err != nil {
		return err
	}

	client := NewClient(GitHubToken)
	client.AcceptHeader = BetaAccessHeader

	wg := sync.WaitGroup{}
	errors := make(chan error, len(deploymentIDs))
	defer close(errors)
	for _, id := range deploymentIDs {
		wg.Add(1)
		go func(deploymentID string) {
			url := fmt.Sprintf("/repos/%s/deployments/%s/statuses", GitHubRepository, deploymentID)
			data := struct {
				State       string `json:"state"`
				Description string `json:"description"`
			}{
				State:       "inactive",
				Description: description,
			}
			err := client.Post(url, data)
			if err != nil {
				errors <- err
			}
			wg.Done()
		}(id)
	}
	wg.Wait()
	if len(errors) > 0 {
		return <-errors
	}
	return nil
}
