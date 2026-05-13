package main

import (
	"github.com/spf13/cobra"

	"github.com/iperalta7/gsm/internal/console"
)

func cmdCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "cmd <command>",
		Short:        "Send a command to the game server",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
	}
	configPath := addConfigFlag(cmd)
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig(*configPath)
		if err != nil {
			return err
		}
		return console.New(cfg).Send(args[0])
	}
	return cmd
}
