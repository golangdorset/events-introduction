package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	kafkaAddressFlag = "kafka-address"
	inputTopicFlag   = "input-topic"
	outputTopicFlag  = "output-topic"
)

var rootCmd = &cobra.Command{
	Use: "proto-demo",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(processCmd)
	rootCmd.AddCommand(accountsCmd)

	rootCmd.PersistentFlags().String(kafkaAddressFlag, "", "kafka address")
	rootCmd.PersistentFlags().String(inputTopicFlag, "", "input kafka topic")
	rootCmd.PersistentFlags().String(outputTopicFlag, "", "output kafka topic")

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		panic("binding flags")
	}
}
