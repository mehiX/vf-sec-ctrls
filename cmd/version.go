package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"text/template"

	"github.com/mehix/vf-sec-ctrls/pkg/version"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionTemplate = `Version:      {{.Version}}
Go version:   {{.GoVersion}}
Built:        {{.BuildTime}}
OS/Arch:      {{.Os}}/{{.Arch}}`

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the current software version",
	Run: func(cmd *cobra.Command, args []string) {
		if err := PrintVersion(os.Stdout); err != nil {
			log.Fatal(err)
		}
		fmt.Println()
	},
}

func PrintVersion(w io.Writer) error {
	t, err := template.New("").Parse(versionTemplate)
	if err != nil {
		return err
	}

	v := struct {
		Version   string
		GoVersion string
		BuildTime string
		Os        string
		Arch      string
	}{
		Version:   version.Version,
		GoVersion: runtime.Version(),
		BuildTime: version.BuildDate,
		Os:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}

	return t.Execute(w, v)
}
