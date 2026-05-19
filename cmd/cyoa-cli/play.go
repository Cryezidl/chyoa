package main

import (
	"fmt"

	"github.com/Cryezidl/cyoa/cli"
	"github.com/spf13/cobra"
)

var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Запустить игру в терминале",
	Run: func(cmd *cobra.Command, args []string) {
		startingArc, _ := cmd.Flags().GetString("sa")
		fmt.Printf("[play.go] starting arc %s", startingArc)

		cli.StartGame(Story, startingArc, Logger)
	},
}

func init() {
	rootCmd.AddCommand(playCmd)
	playCmd.Flags().String("sa", "intro", "Starting arc for adventure")
}
