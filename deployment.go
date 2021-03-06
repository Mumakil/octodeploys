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
func (d *Deployment) String(includeStatus bool) string {
	var status string
	if includeStatus {
		if d.Status == nil {
			status = "\tno status"
		} else {
			status = fmt.Sprintf("\t%s", d.Status.State)
		}
	}
	// Make deploy by id a bit nicer
	ref := d.Ref
	if d.Ref == d.Sha {
		ref = d.Ref[:8]
	}
	return fmt.Sprintf(
		"%d%s\t%s\t%s\t%s\t%s\t%s",
		d.ID,
		status,
		ref,
		d.Sha[:8],
		d.Environment,
		d.CreatedAt.Format(time.RFC3339),
		d.Creator.Login,
	)
}

// String formats a string for output of the deployments
func (d Deployments) String(includeStatus bool) string {
	if len(d) == 0 {
		return "No deployments"
	}
	status := ""
	if includeStatus {
		status = "\tStatus"
	}
	res := fmt.Sprintf("ID%s\tRef\tSha\tEnvironment\tDate\tDeployer\n", status)
	for _, deployment := range d {
		res += deployment.String(includeStatus) + "\n"
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
