package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/iperalta7/gsm/internal/process"
)

func startCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the game server",
	}
	configPath := addConfigFlag(cmd)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig(*configPath)
		if err != nil {
			return err
		}
		fmt.Printf("starting %s...\n", cfg.Name)
		if err := process.New().Start(cfg.Service.Name); err != nil {
			return err
		}
		fmt.Println("started")
		return nil
	}
	return cmd
}
