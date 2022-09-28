# gh-discussions

> `gh-discussions` is [GitHub CLI](https://cli.github.com) extension that allows you to manage discussions on GitHub.

## Install

```bash
gh extension install antgrutta/gh-discussions
```

## Usage

```txt
gh discussions [command] [flags]
```

```txt
gh cli to manage discussions on GitHub.
        It supports a variety of operations including importing discussions from csv files.

Usage:
  gh-discussions [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  list        Print a list of discussions given a repository
  migration   Import discussions from a csv file.

Flags:
  -a, --access-token string   Github access token (default: "")
      --config string         config file (default is $HOME/.gh-discussions.yaml)
  -h, --help                  help for gh-discussions
  -l, --logfile string        Log file (default "error.log")
  -t, --toggle                Help message for toggle

Use "gh-discussions [command] --help" for more information about a command.
```

## License

- [MIT](./license) (c) [Anthony Grutta](https://github.com/antrgrutta)
- [Code of Conduct](./.github/code_of_conduct.md)
