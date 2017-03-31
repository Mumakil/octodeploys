# github-deployments

CLI for managing deployments in GitHub repositories

## Installing

You'll need a [working Go environment](https://golang.org/doc/install).

Install simply by doing

```sh
go get github.com/Mumakil/github-deployments
```

## Usage

You need to specify repository and GitHub access token to use the cli. You can either use environment variables (`GITHUB_TOKEN` and `GITHUB_REPOSITORY`) or command line flags (`-token`, `-repository`).

Help

```sh
github-deployments help
```

List deployments

```sh
github-deployments list
```

List deployments, filter by state and environment, include statuses.

```sh
github-deployments -statuses -state error -environment production list
```

## License

MIT
