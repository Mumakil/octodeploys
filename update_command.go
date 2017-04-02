package main

import (
	"flag"
	"fmt"
	"strconv"
	"sync"
)

// BetaAccessHeader required to use inactive statuses in GitHub deployments api
const BetaAccessHeader = "application/vnd.github.ant-man-preview+json"

// Description description for the state
var Description string

// NewState new state for the deployment
var NewState string

func init() {
	flag.StringVar(&Description, "description", "", "description for the status")
	flag.StringVar(&NewState, "newState", "", "new state for the deployment")
}

func updateCommand(rawDeploymentIDs []string) error {
	err := validateGlobalArgs()
	if err != nil {
		return err
	}
	if NewState == "" {
		return fmt.Errorf("missing newState")
	}

	deploymentIDs := make([]uint64, 0, len(rawDeploymentIDs))
	for _, rawID := range rawDeploymentIDs {
		id, err := strconv.ParseUint(rawID, 10, 64)
		if err != nil {
			return fmt.Errorf("deactivate command: error parsing deployment id: %s", err.Error())
		}
		deploymentIDs = append(deploymentIDs, id)
	}

	client := NewClient(GitHubToken)
	client.AcceptHeader = BetaAccessHeader

	return updateAll(client, deploymentIDs)
}

func updateAll(client *Client, deploymentIDs []uint64) error {
	wg := sync.WaitGroup{}
	errors := make(chan error, len(deploymentIDs))
	defer close(errors)
	for _, id := range deploymentIDs {
		wg.Add(1)
		go func(deploymentID uint64) {
			err := updateDeployment(client, deploymentID)
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

func updateDeployment(client *Client, deploymentID uint64) error {
	url := fmt.Sprintf("/repos/%s/deployments/%d/statuses", GitHubRepository, deploymentID)
	data := struct {
		State       string `json:"state"`
		Description string `json:"description"`
	}{
		State:       NewState,
		Description: Description,
	}
	return client.Post(url, data)
}
