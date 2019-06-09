package cmd

import (
	"github.com/adambaumeister/panfw-util/panos/device"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add objects to the PANOS device.",
	Run: func(cmd *cobra.Command, args []string) {
		username = viper.GetString("user")
		password = viper.GetString("password")
		hostname = viper.GetString("hostname")

		fw := device.ConnectUniversal(username, password, hostname)
		fw.SetDeviceGroup(devicegroup)
		fw.Add(args)
	},
}
