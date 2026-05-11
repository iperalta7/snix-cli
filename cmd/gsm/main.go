package main

import (
	"fmt"
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
	root.AddCommand(validateCmd())
	return root
}

func validateCmd() *cobra.Command {
	var configPath string

	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate a server config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(configPath)
			if err != nil {
				return err
			}
			if err := cfg.Validate(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			fmt.Printf("config valid: %s", cfg.Name)
			if cfg.Game != "" {
				fmt.Printf(" (%s)", cfg.Game)
			}
			fmt.Printf(" — %s runtime\n", cfg.Runtime)
			return nil
		},
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", "/etc/gsm/config.yaml", "path to server config")
	return cmd
}
