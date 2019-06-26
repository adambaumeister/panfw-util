package cmd

import (
	"fmt"
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

		fw := device.Connect(username, password, hostname)
		ParseAndTestPcap(args[0], fw)
	},
}

func ParseAndTestPcap(fn string, fw *device.Firewall) {
	flows := pcaptest.ReadPcap(fn)
	if len(flows) > maxTests {
		fmt.Printf("Number of flows in PCAP exceeds maximum. Testing the first %v seen. Use --max to override.\n", maxTests)
		flows = flows[:maxTests]
	}
	for _, f := range flows {
		fmt.Printf("Testing flow...")
		f.To = &toZone
		f.From = &fromZone
		f.Print()
		rules := testcmd.TestPolicy(hostname, fw.Apikey, flows[0])
		fmt.Printf("name: %v\n", rules[0].Name)
	}
}
