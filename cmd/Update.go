/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"nomad-image-updater/nid"
	"github.com/spf13/cobra"
)

// UpdateCmd represents the Update command
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updating image from file in parameter if folder parse all file contains ",
	Run: func(cmd *cobra.Command, args []string) {
		nid.Update(args[0])
	},
}

func init() {
	rootCmd.AddCommand(UpdateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// UpdateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// UpdateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
