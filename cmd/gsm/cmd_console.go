package main

import (
	"github.com/spf13/cobra"

	"github.com/iperalta7/gsm/internal/console"
)

func consoleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "console",
		Short:        "Attach to the game server console",
		SilenceUsage: true,
	}
	configPath := addConfigFlag(cmd)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig(*configPath)
		if err != nil {
			return err
		}
		return console.New(cfg).Attach()
	}
	return cmd
}
