package cmd

import (
	"fmt"
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

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register IP addresses with dynamic tags.",
	Run: func(cmd *cobra.Command, args []string) {
		username = viper.GetString("user")
		password = viper.GetString("password")
		hostname = viper.GetString("hostname")

		fw := device.ConnectUniversal(username, password, hostname)
		fw.SetDeviceGroup(devicegroup)
		r := fw.Register(args)
		fmt.Printf("Result: %v\n", r.Status)
	},
}

var unregisterCmd = &cobra.Command{
	Use:   "unregister",
	Short: "Unregister dynamic tags associated with an IP address.",
	Run: func(cmd *cobra.Command, args []string) {
		username = viper.GetString("user")
		password = viper.GetString("password")
		hostname = viper.GetString("hostname")

		fw := device.ConnectUniversal(username, password, hostname)
		fw.SetDeviceGroup(devicegroup)
		r := fw.UnRegister(args)
		fmt.Printf("Result: %v\n", r.Status)
	},
}
