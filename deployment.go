package main

import (
	"fmt"
	"time"
)

// Deployment holds data about a deployment in GitHub
type Deployment struct {
	ID          uint64    `json:"id"`
	URL         string    `json:"url"`
	Sha         string    `json:"sha"`
	Ref         string    `json:"ref"`
	Environment string    `json:"environment"`
	CreatedAt   time.Time `json:"created_at"`
	Creator     User      `json:"creator"`
	Status      *Status   `json:"-"`
}

// Deployments list of deployments
type Deployments []*Deployment

// String formats a string for output of the deployment
func (d *Deployment) String() string {
	var status string
	if d.Status != nil {
		status = fmt.Sprintf(" - %s", d.Status.State)
	}
	return fmt.Sprintf(
		"%d%s - %s (%s) in %s at %s by %s",
		d.ID,
		status,
		d.Ref,
		d.Sha[:8],
		d.Environment,
		d.CreatedAt.Format(time.RFC3339),
		d.Creator.Login,
	)
}

// String formats a string for output of the deployments
func (d Deployments) String() string {
	if len(d) == 0 {
		return "No deployments"
	}
	res := ""
	for _, deployment := range d {
		res += deployment.String() + "\n"
	}
	return res
}

// FilterByState returns only deployments matching the specified state
func (d Deployments) FilterByState(state string) Deployments {
	deployments := Deployments{}
	for _, deployment := range d {
		if deployment.Status != nil && deployment.Status.State == state {
			deployments = append(deployments, deployment)
		}
	}
	return deployments
}
