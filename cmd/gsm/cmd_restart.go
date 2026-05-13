package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/iperalta7/gsm/internal/process"
)

func restartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restart",
		Short: "Restart the game server",
	}
	configPath := addConfigFlag(cmd)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig(*configPath)
		if err != nil {
			return err
		}
		fmt.Printf("restarting %s...\n", cfg.Name)
		if err := process.New().Restart(cfg.Service.Name); err != nil {
			return err
		}
		fmt.Println("restarted")
		return nil
	}
	return cmd
}
