/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/Cryezidl/cyoa/cyoa"
	"github.com/spf13/cobra"
)

var Story cyoa.Story
var Logger *slog.Logger

var rootCmd = &cobra.Command{
	Use:   "cyoa-cli",
	Short: "",
	Long:  ``,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString("path")
		if strings.TrimSpace(path) == "" {
			return fmt.Errorf("path to the story file must not be empty")
		}

		Logger = slog.Default()
		s, err := cyoa.LoadStory(path)
		if err != nil {
			return fmt.Errorf("loading story: %w", err)
		}
		Story = s
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("path", "p", "gopher.json", "Path to the story file")
}
