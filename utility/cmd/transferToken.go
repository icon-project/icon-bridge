/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var file_path string

// transferTokenCmd represents the transferToken command
var transferTokenCmd = &cobra.Command{
	Use:   "transferToken",
	Short: "To Transfer token between source and Destination",
	Long:  ``,
}

func init() {
	rootCmd.AddCommand(transferTokenCmd)
	transferTokenCmd.AddCommand(nearCmd)
	transferTokenCmd.AddCommand(iconCmd)
}
