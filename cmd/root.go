package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	xlsFile  string
	xlsSheet string
)

var rootCmd = &cobra.Command{
	Use:   "sec-controls [list | serve]",
	Short: "Show security controls loaded from an Excel file",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if xlsFile == "" {
			fmt.Println("Missing Excel file to read from. You can upload one")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	rootCmd.PersistentFlags().StringVarP(&xlsFile, "xls-file", "x", "", "Excel file to read the values from")
	rootCmd.PersistentFlags().StringVarP(&xlsSheet, "xls-sheet", "s", "", "Worksheet that contains the controls data")

	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
