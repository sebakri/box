package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"box/internal/config"
	"box/internal/installer"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:       "generate <type>",
	Short:     "Generates configuration files",
	Long:      `Generates configuration files for shell integration, such as direnv.`,
	ValidArgs: []string{"direnv"},
	Args:      cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		genType := args[0]
		if genType != "direnv" {
			fmt.Printf("Unknown generation type: %s\n", genType)
			os.Exit(1)
		}

		configFile := "box.yml"
		cfg, err := config.Load(configFile)
		if err != nil {
			cfg = &config.Config{}
		}

		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get current working directory: %v", err)
		}

		mgr := installer.New(cwd, cfg.Env)
		mgr.Output = io.Discard
		if err := mgr.EnsureEnvrc(); err != nil {
			log.Fatalf("Failed to generate .envrc: %v", err)
		}

		fmt.Printf("%s Generated .envrc\n", successStyle.Render("✅"))
		if err := mgr.AllowDirenv(); err != nil {
			fmt.Printf("%s Failed to run direnv allow: %v\n", warnStyle.Render("⚠️"), err)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
