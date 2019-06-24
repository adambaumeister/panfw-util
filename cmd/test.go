package cmd

import (
	"github.com/adambaumeister/panfw-util/panos/api/testcmd"
	"github.com/adambaumeister/panfw-util/panos/device"
	"github.com/adambaumeister/panfw-util/pcaptest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var testPcap = &cobra.Command{
	Use:   "testpcap",
	Short: "Read a PCAP file, testing all the flows against ",
	Run: func(cmd *cobra.Command, args []string) {
		username = viper.GetString("user")
		password = viper.GetString("password")
		hostname = viper.GetString("hostname")

		flows := pcaptest.ReadPcap(args[0])
		fw := device.Connect(username, password, hostname)
		testcmd.TestPolicy(hostname, fw.Apikey, flows[0])
	},
}
