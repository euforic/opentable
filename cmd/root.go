package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RootCmd is only exposed so it can be executed in
// main.go in the repos root package
var RootCmd = &cobra.Command{
	Use:   "opentable",
	Short: "A gRPC based opentable scraper",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
