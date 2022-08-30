/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

const iconBinaryPath = "../bazel-bin/external/com_github_icon_project_goloop/cmd/goloop/goloop_/goloop"

// iconCmd represents the icon command
var iconCmd = &cobra.Command{
	Use:   "icon",
	Short: "Transfer token from Icon to Destination",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		configs, err := NewConfig(file_path)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, config := range configs {
			if config.Network == "icon" {
				res := TransferFromIcon(&config)
				fmt.Println(res)

			}
		}
	},
}

func init() {
	rootCmd.AddCommand(iconCmd)
	iconCmd.PersistentFlags().StringVarP(&file_path, "config", "c", "", "Config File containing deatils for token transfer")
	iconCmd.MarkPersistentFlagRequired("config")
}

func TransferFromIcon(config *Config) string {
	reciever := fmt.Sprintf("_to=\"%s\"", config.Reciever)

	cmd := exec.Command(iconBinaryPath, "rpc", "--uri", config.Uri, "sendtx", "call", "--to", config.BtsAddress, "--method", "transferNativeCoin", "--param", reciever, "--value", config.Value, "--nid", config.NetworkId, "--step_limit", fmt.Sprintf("%d", config.StepLimit), "--key_store", config.KeyStore, "--key_secret", config.KeySecret)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err)
	}
	return string(stdout)
}
