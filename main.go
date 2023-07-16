package main

import (
	"fmt"
	"github.com/alecthomas/kong"
)

var CLI struct {
	Repos    ReposCommand    `cmd:"repos" help:"List repositories"`
	Branches BranchesCommand `cmd:"branches" help:"List branches"`
}

type ReposCommand struct {
	FilterRegex     string `help:"Repository filter (regex)" short:"f"`
	CurrentUserOnly bool   `help:"List repositories owned by the current logged-in user only"`
}

type BranchesCommand struct {
	FilterRegex string `help:"Branch filter (regex)" short:"f"`
	OpenPRSOnly bool   `help:"Filter for only open pull requests"`
}

func main() {
	ctx := kong.Parse(&CLI)
	switch ctx.Command() {
	case "repos":
		repos_main(CLI.Repos.FilterRegex, CLI.Repos.CurrentUserOnly)
	case "branches":
		branches_main(CLI.Branches.FilterRegex, CLI.Branches.OpenPRSOnly)
	default:
		ctx.FatalIfErrorf(fmt.Errorf("invalid subcommand: %s", ctx.Command()))
	}
}
