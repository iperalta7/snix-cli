package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/iperalta7/gsm/internal/process"
)

func statusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show game server status",
	}
	configPath := addConfigFlag(cmd)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig(*configPath)
		if err != nil {
			return err
		}
		s, err := process.New().GetStatus(cfg.Service.Name)
		if err != nil {
			return err
		}
		fmt.Println(s)
		return nil
	}
	return cmd
}
