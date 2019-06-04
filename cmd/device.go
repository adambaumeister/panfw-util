package cmd

import (
	"github.com/adamb/panfw-util/panos/device"
	"github.com/adamb/panfw-util/panos/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load, and commit, an XML configuration file from disk.",
	Run: func(cmd *cobra.Command, args []string) {
		username = viper.GetString("user")
		password = viper.GetString("password")
		hostname = viper.GetString("hostname")

		fw := device.Connect(username, password, hostname)
		filename := args[0]
		fw.ImportNamed(filename)

		file, err := os.Open(filename)
		errors.DieIf(err)

		fi, err := file.Stat()
		fw.LoadNamed(fi.Name(), commit)
	},
	Args: cobra.MinimumNArgs(1),
}
