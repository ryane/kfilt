package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var GitCommit, GitState, Version string

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Display the kfilt version",
		Run: func(cmd *cobra.Command, args []string) {

			if GitState != "clean" {
				fmt.Printf("%s (%s-%s)\n", Version, GitCommit, GitState)
			} else {
				fmt.Printf("%s (%s)\n", Version, GitCommit)
			}
		},
	}
}
