package cmd

import (
	"github.com/adambaumeister/panfw-util/panos/device"
	"github.com/spf13/cobra"
)

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Dump various aspects of a PAN device",
	Run: func(cmd *cobra.Command, args []string) {
		hostname = PromptIfNil("hostname", false)
		password = PromptIfNil("password", true)
		username = PromptIfNil("user", false)

		d := device.ConnectUniversal(username, password, hostname)
		d.SetDeviceGroup(devicegroup)
		d.Print(args[0])
	},
	Args: cobra.MinimumNArgs(1),
}
