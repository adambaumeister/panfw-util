package cmd

import (
	"github.com/adambaumeister/panfw-util/panos/device"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Dump various aspects of a PAN device",
	Run: func(cmd *cobra.Command, args []string) {
		username = viper.GetString("user")
		password = viper.GetString("password")
		hostname = viper.GetString("hostname")

		fw := device.Connect(username, password, hostname)
		fw.Print(args[0])
	},
	Args: cobra.MinimumNArgs(1),
}
