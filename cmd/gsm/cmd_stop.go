package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/iperalta7/gsm/internal/process"
)

func stopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop the game server",
	}
	configPath := addConfigFlag(cmd)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig(*configPath)
		if err != nil {
			return err
		}
		fmt.Printf("stopping %s...\n", cfg.Name)
		if err := process.New().Stop(cfg.Service.Name); err != nil {
			return err
		}
		fmt.Println("stopped")
		return nil
	}
	return cmd
}
