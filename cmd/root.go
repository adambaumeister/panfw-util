package cmd

import (
	"fmt"
	"github.com/adambaumeister/panfw-util/panos/errors"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var cfgFile string
var username string
var password string
var hostname string
var commit bool
var devicegroup string
var filename string

var rootCmd = &cobra.Command{
	Use:   "panutil",
	Short: "Golang based PANOS utilities.",
	Long: `Collection of useful and fast functionality for
interacting with PANOS devices.`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of panutil",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("PANUTIL Version 0.0\n")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// Root flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().BoolVar(&errors.DEBUG, "debug", false, "Enable verbose debugging.")
	rootCmd.PersistentFlags().StringVar(&username, "username", "", "Login Username")
	rootCmd.PersistentFlags().StringVar(&password, "password", "", "Login Password")
	rootCmd.PersistentFlags().StringVar(&hostname, "hostname", "", "PANOS device hostname")

	// Load config flags
	loadCmd.Flags().BoolVar(&commit, "commit", false, "Commit the configuration.")

	// Add flaggs
	addCmd.Flags().StringVar(&devicegroup, "devicegroup", "shared", "Set the device group if targeting Panorama.")

	viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("hostname", rootCmd.PersistentFlags().Lookup("hostname"))
	rootCmd.AddCommand(versionCmd)

	rootCmd.AddCommand(loadCmd)
	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(printCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(testPcap)
}

func initConfig() {
	viper.SetEnvPrefix("panutil")
	viper.BindEnv("user")
	viper.BindEnv("hostname")
	viper.BindEnv("password")
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath("./")
		viper.SetConfigName(".panutil")
	}

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		return
		//os.Exit(1)
	}
}
