package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is the current version of box, set during build time.
var Version = "dev"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of box",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("box", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
