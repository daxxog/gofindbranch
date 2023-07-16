package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a subcommand: branches or repos")
		return
	}

	subcommand := os.Args[1]

	switch subcommand {
	case "branches":
		branches_main()
	case "repos":
		repos_main()
	default:
		fmt.Println("Invalid subcommand. Available options: branches, repos")
	}
}
