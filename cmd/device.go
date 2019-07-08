package cmd

import (
	"github.com/adambaumeister/panfw-util/panos/device"
	"github.com/adambaumeister/panfw-util/panos/errors"
	"github.com/spf13/cobra"
	"os"
)

var loadCmd = &cobra.Command{
	Use:   "load [path to config]",
	Short: "Load, and commit, an XML configuration file from disk.",
	Run: func(cmd *cobra.Command, args []string) {
		hostname = PromptIfNil("hostname", false)
		password = PromptIfNil("password", true)
		username = PromptIfNil("user", false)

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

var importCmd = &cobra.Command{
	Use:   "import [path to config]",
	Short: "Import - without loading - a named configuration file.",
	Run: func(cmd *cobra.Command, args []string) {
		hostname = PromptIfNil("hostname", false)
		password = PromptIfNil("password", true)
		username = PromptIfNil("user", false)

		fw := device.Connect(username, password, hostname)
		filename := args[0]
		fw.ImportNamed(filename)
	},
	Args: cobra.MinimumNArgs(1),
}
