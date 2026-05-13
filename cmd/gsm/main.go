package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/iperalta7/gsm/internal/config"
)

func main() {
	if err := rootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func rootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "gsm",
		Short: "Game server manager",
	}
	root.AddCommand(
		validateCmd(),
		startCmd(),
		stopCmd(),
		restartCmd(),
		statusCmd(),
		consoleCmd(),
		cmdCmd(),
	)
	return root
}

// ── Shared helpers ────────────────────────────────────────────────────────────

const defaultConfigPath = "/etc/gsm/config.yaml"

func loadConfig(path string) (*config.Config, error) {
	return config.Load(path)
}

func addConfigFlag(cmd *cobra.Command) *string {
	var p string
	cmd.Flags().StringVarP(&p, "config", "c", defaultConfigPath, "path to server config")
	return &p
}
