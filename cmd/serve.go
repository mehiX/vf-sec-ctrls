package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mehix/vf-sec-ctrls/categ"
	"github.com/mehix/vf-sec-ctrls/server"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the API server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addr := args[0]

		h, err := categ.NewFromFile(xlsFile, xlsSheet)
		if xlsFile != "" && err != nil {
			log.Println(err)
		}

		fmt.Printf("Listen on %s\n", addr)
		srvr := http.Server{
			Addr:    addr,
			Handler: server.New().Handlers(h),
		}
		if err := srvr.ListenAndServe(); err != nil {
			log.Println(err)
		}
	},
}
