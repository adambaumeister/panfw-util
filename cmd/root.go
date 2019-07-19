package cmd

import (
	"bufio"
	"fmt"
	"github.com/adambaumeister/panfw-util/panos/errors"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"syscall"
)

var cfgFile string
var username string
var password string
var hostname string
var commit bool
var devicegroup string
var filename string
var maxTests int
var fromZone string
var toZone string
var count int
var logtype string
var joinFilter string
var joinFilterVal string

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

	// print flags
	printCmd.Flags().StringVar(&devicegroup, "devicegroup", "shared", "Set the device group if targeting Panorama.")

	// Logs flags
	logCmd.Flags().IntVar(&count, "count", 20, "Limit the returned count of logs.")
	logCmd.Flags().StringVar(&logtype, "type", "traffic", "Specify the log to query.")

	// Join command flats
	joinCmd.Flags().IntVar(&count, "count", 20, "Limit the returned count of logs.")
	joinCmd.Flags().StringVar(&logtype, "type", "traffic", "Specify the log to query.")
	joinCmd.Flags().StringVar(&joinFilter, "filterkey", "", "Field to filter on in joined object.")
	joinCmd.Flags().StringVar(&joinFilterVal, "filterval", ".*", "Regex supported value to search for within value.")

	viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("hostname", rootCmd.PersistentFlags().Lookup("hostname"))
	rootCmd.AddCommand(versionCmd)

	rootCmd.AddCommand(loadCmd)
	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(printCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(joinCmd)
	rootCmd.AddCommand(logCmd)
	rootCmd.AddCommand(registerCmd)
	rootCmd.AddCommand(unregisterCmd)
	rootCmd.AddCommand(apiCmd)
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

func PromptIfNil(p string, secret bool) string {
	// If the val is nul
	v := viper.GetString(p)
	if v == "" {
		fmt.Printf("Enter value for %v: ", p)
		if secret {
			bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
			errors.DieIf(err)
			password := string(bytePassword)
			return password
		} else {
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			return text
		}
	}
	// otherwise just return the value DUH
	return v

}
