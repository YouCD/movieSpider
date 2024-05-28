package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

//nolint:gochecknoglobals
var (
	Version   string
	commitID  string
	buildTime string
	goVersion string
	buildUser string
)

//nolint:exhaustruct,gochecknoglobals,forbidigo
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version info of " + Name,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf("Version:   %s\n", Version)
		fmt.Printf("CommitID:  %s\n", commitID)
		fmt.Printf("BuildTime: %s\n", buildTime)
		fmt.Printf("GoVersion: %s\n", goVersion)
		fmt.Printf("BuildUser: %s\n", buildUser)
	},
}
