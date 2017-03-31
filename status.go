package main

import "time"

// Status is a GitHub deployment status
type Status struct {
	ID          uint64    `json:"id"`
	State       string    `json:"state"`
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description"`
}
