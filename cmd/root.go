/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"
	"github.com/spf13/cobra"
	"nomad-image-updater/internal/config"
	"log/slog"
)



// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nomad-image-updater",
	Short: "light tool to update docker image in nomad file",
	Long: `A longer descriptiolication.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nomad-image-updater.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	config:=config.GetConfig()
	lvl := &slog.LevelVar{}
	lvl.UnmarshalText([]byte(config.LoggerOption.Level))
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: lvl,
	}))
	slog.SetDefault(logger)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

