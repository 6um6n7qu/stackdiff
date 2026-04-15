package main

import (
	"fmt"
	"os"

	"github.com/stackdiff/stackdiff/internal/config"
	"github.com/stackdiff/stackdiff/internal/diff"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: stackdiff <config-a> <config-b>\n")
		fmt.Fprintf(os.Stderr, "Example: stackdiff service-a.env service-b.env\n")
		os.Exit(1)
	}

	pathA := os.Args[1]
	pathB := os.Args[2]

	cfgA, err := config.LoadFromFile(pathA)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading %s: %v\n", pathA, err)
		os.Exit(1)
	}

	cfgB, err := config.LoadFromFile(pathB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading %s: %v\n", pathB, err)
		os.Exit(1)
	}

	result := diff.Compare(cfgA, cfgB)
	diff.Print(result, pathA, pathB)

	if len(result) > 0 {
		os.Exit(2)
	}
}
