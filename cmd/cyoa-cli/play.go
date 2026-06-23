package main

import (
	"fmt"

	"github.com/Cryezidl/cyoa/cli"
	"github.com/spf13/cobra"
)

var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Start game in CLI",
	RunE: func(cmd *cobra.Command, args []string) error {
		if Story == nil {
			return fmt.Errorf("no story loaded: check the --path flag")
		}
		startingArc, _ := cmd.Flags().GetString("sa")
		return cli.StartGame(Story, startingArc, Logger)
	},
}

func init() {
	rootCmd.AddCommand(playCmd)
	playCmd.Flags().String("sa", "intro", "Starting arc for adventure")
}
