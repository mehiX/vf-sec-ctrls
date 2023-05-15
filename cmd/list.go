package cmd

import (
	"log"

	"github.com/mehix/vf-sec-ctrls/pkg/service/categ"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&showCmd)
}

var showCmd = cobra.Command{
	Use:   "list",
	Short: "Print all the controls loaded from Excel",
	Run: func(cmd *cobra.Command, args []string) {
		h, err := categ.NewFromFile(xlsFile, xlsSheet)
		if err != nil {
			log.Println(err)
		}
		h.Print(true)
	},
}
