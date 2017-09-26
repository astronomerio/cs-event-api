package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	Version    = "0.0.0"
	CommitHash = ""
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of the clickstream-api server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
		fmt.Println(CommitHash)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
