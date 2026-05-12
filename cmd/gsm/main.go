package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/iperalta7/gsm/internal/config"
	"github.com/iperalta7/gsm/internal/process"
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
	)
	return root
}

// ── Helpers ───────────────────────────────────────────────────────────────────

const defaultConfigPath = "/etc/gsm/config.yaml"

func loadConfig(path string) (*config.Config, error) {
	return config.Load(path)
}

func addConfigFlag(cmd *cobra.Command) *string {
	var p string
	cmd.Flags().StringVarP(&p, "config", "c", defaultConfigPath, "path to server config")
	return &p
}

// ── validate ──────────────────────────────────────────────────────────────────

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

// ── start ─────────────────────────────────────────────────────────────────────

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

// ── stop ──────────────────────────────────────────────────────────────────────

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

// ── restart ───────────────────────────────────────────────────────────────────

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

// ── status ────────────────────────────────────────────────────────────────────

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
