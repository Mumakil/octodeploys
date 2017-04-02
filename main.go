package main // import "github.com/Mumakil/github-deployments"

import (
	"flag"
	"fmt"
	"os"
)

// GitHubToken access token used
var GitHubToken string

// GitHubRepository repo to manipulate
var GitHubRepository string

const help = `github-deployments manages GitHub deployments for a repository

Available commands:

	help - print this help
	list - list open deployments
	deactivate - deactivate provided deployments
`

func validateGlobalArgs() error {
	if GitHubRepository == "" {
		return fmt.Errorf("missing repository - provide one with flag -repository or environment variable GITHUB_REPOSITORY")
	}
	if GitHubToken == "" {
		return fmt.Errorf("missing GitHub access token - provide one with flag -token or environment variable GITHUB_TOKEN")
	}
	return nil
}

func init() {
	flag.StringVar(&GitHubRepository, "repository", "", "Github repository to use")
	flag.StringVar(&GitHubToken, "token", "", "GitHub access token")
}

func main() {
	flag.Parse()

	if GitHubToken == "" {
		GitHubToken = os.Getenv("GITHUB_TOKEN")
	}
	if GitHubRepository == "" {
		GitHubRepository = os.Getenv("GITHUB_REPOSITORY")
	}

	command := flag.Arg(0)
	args := flag.Args()[1:]

	var err error

	switch command {
	case "help":
		fmt.Println(help)
	case "list":
		err = listCommand(args)
	case "update":
		err = updateCommand(args)
	case "updateByState":
		err = updateByStateCommand(args)
	default:
		fmt.Printf("Unrecognized command \"%s\"\n\n", command)
		fmt.Println(help)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running command \"%s\": %s\n", command, err.Error())
		os.Exit(1)
	}
}
