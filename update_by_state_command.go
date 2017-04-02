package main

import (
	"flag"
	"fmt"
)

// Except don't touch this deployment id
var Except uint64

func init() {
	flag.Uint64Var(&Except, "except", 0, "don't touch this specific deployment id")
}

func updateByStateCommand(args []string) error {
	err := validateGlobalArgs()
	if err != nil {
		return err
	}
	if State == "" {
		return fmt.Errorf("missing state")
	}
	if NewState == "" {
		return fmt.Errorf("missing newState")
	}

	client := NewClient(GitHubToken)
	client.AcceptHeader = BetaAccessHeader

	deployments, err := fetchDeployments(client)
	if err != nil {
		return err
	}
	err = fetchStatuses(client, deployments)
	if err != nil {
		return err
	}

	filtered := deployments.FilterByState("success")
	deploymentIDs := make([]uint64, 0, len(filtered))
	for _, deployment := range filtered {
		if deployment.ID != Except {
			deploymentIDs = append(deploymentIDs, deployment.ID)
		}
	}
	return updateAll(client, deploymentIDs)
}
