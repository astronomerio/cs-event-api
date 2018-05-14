package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version    string
	gitCommit  string
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Event API version",
		Long:  "Event API version",
		Run:   printVersion,
	}
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

func printVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("Event API Version: %s\n", version)
	fmt.Printf("Git Commit: %s\n", gitCommit)
}
