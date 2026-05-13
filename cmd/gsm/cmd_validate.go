package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func validateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "validate",
		Short:        "Validate a server config file",
		SilenceUsage: true,
	}
	configPath := addConfigFlag(cmd)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig(*configPath)
		if err != nil {
			return err
		}
		if err := cfg.Validate(); err != nil {
			return err
		}
		fmt.Printf("config valid: %s", cfg.Name)
		if cfg.Game != "" {
			fmt.Printf(" (%s)", cfg.Game)
		}
		fmt.Printf(" — %s runtime\n", cfg.Runtime)
		return nil
	}
	return cmd
}
