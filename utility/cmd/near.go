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

const nearBinaryPath = "../bazel-bin/external/near/cli/near_binary.sh"

// nearCmd represents the near command
var nearCmd = &cobra.Command{
	Use:   "near",
	Short: "Transfer tokens from Near network to destination",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		configs, err := NewConfig(file_path)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, config := range configs {

			if config.Network == "near" {
				res := TransferFromNear(&config)
				fmt.Println(res)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(nearCmd)
	nearCmd.PersistentFlags().StringVarP(&file_path, "config", "c", "", "Config File containing deatils for token transfer")
	nearCmd.MarkPersistentFlagRequired("config")
}

func doDeposit(config *Config) (string, error) {

	cmd := exec.Command(nearBinaryPath, "call", config.BtsAddress, "deposit", "'{}'", "--accountId", config.Sender, "--amount", "10")
	stdout, err := cmd.Output()

	if err != nil {
		return "", err
	}
	return string(stdout), nil
}

func TransferFromNear(config *Config) string {

	deposit, err := doDeposit(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(deposit)

	params := fmt.Sprintf("'{\"coin_id\":[247,184,188,27,185,62,246,213,77,228,101,193,206,130,171,233,39,148,195,217,177,33,141,212,41,223,44,241,83,209,37,217],\"destination\":\"%s\",\"amount\":\"%s\"}'", config.Reciever, config.Value)
	cmd := exec.Command(nearBinaryPath, "call", config.BtsAddress, "transfer", params, "--accountId", "btp-16.testnet", "--gas", "300000000000000")
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
	}
	return string(stdout)
}
