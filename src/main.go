package main

import (
	"fmt"
	"os"

	"github.com/team-project-gopherit/menu"
)

func main() {
	startDir, err := startDirFromArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := menu.RunPolo(startDir); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func startDirFromArgs() (string, error) {
	if len(os.Args) > 1 {
		return os.Args[1], nil
	}

	return os.Getwd()
}
