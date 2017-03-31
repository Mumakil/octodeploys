package main

// User is a GitHub user
type User struct {
	ID    uint64 `json:"id"`
	Login string `json:"login"`
}
