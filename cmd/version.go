package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the software version",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: make it dynamic
		fmt.Println("1.0.0")
	},
}
