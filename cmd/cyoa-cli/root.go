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

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		filepath, _ := cmd.Flags().GetString("p")
		if strings.TrimSpace(filepath) == "" {
			fmt.Printf("filepath is %s\n", filepath)
			return
		}

		Logger = slog.Default()
		s, err := cyoa.LoadStory(filepath, Logger)
		if err != nil {
			return
		}
		Story = s
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().String("p", "C:\\Users\\User\\OneDrive\\Рабочий стол\\goProjects\\cyoa\\gopher.json", "Path to the file")
}
