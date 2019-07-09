package cmd

import (
	"fmt"
	"github.com/adambaumeister/panfw-util/panos/api/logs"
	"github.com/adambaumeister/panfw-util/panos/api/policy"
	"github.com/adambaumeister/panfw-util/panos/device"
	"github.com/adambaumeister/panfw-util/panos/errors"
	"github.com/spf13/cobra"
	"regexp"
	"strings"
)

var logCmd = &cobra.Command{
	Use:   "logs",
	Short: "Search and print logs",
	Run: func(cmd *cobra.Command, args []string) {
		hostname = PromptIfNil("hostname", false)
		password = PromptIfNil("password", true)
		username = PromptIfNil("user", false)

		fw := device.ConnectUniversal(username, password, hostname)
		entries := fw.LogQuery(args, count, logtype)
		for _, e := range entries {
			e.Print()
		}
	},
}

var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "Joins log entries with their underyling PAN objects.",
	Run: func(cmd *cobra.Command, args []string) {
		hostname = PromptIfNil("hostname", false)
		password = PromptIfNil("password", true)
		username = PromptIfNil("user", false)

		fw := device.ConnectUniversal(username, password, hostname)
		entries := fw.LogQuery(args, count, logtype)
		JoinLogsWithRuleObjects(fw, entries)
	},
}

func JoinLogsWithRuleObjects(fw device.Panos, logs []*logs.LogEntry) {
	apikey := fw.GetApiKey()
	fqdn := fw.GetHostname()

	// Print the field names
	if len(logs) > 0 {
		log := logs[0]
		fields, _ := log.ToFields()
		fmt.Printf("%v\n", strings.Join(fields, ", "))
	}

	for _, log := range logs {
		xpath := strings.Split(log.FullPath, "/")
		r := policy.GetRules(fqdn, apikey, xpath)
		if len(r) > 0 {
			rule := r[0]
			if joinFilter != "" {
				fieldVal := rule.Lookup(joinFilter)
				match, err := regexp.MatchString(joinFilterVal, fieldVal)
				errors.DieIf(err)
				if match {
					_, vals := log.ToFields()
					fmt.Printf("%v\n", strings.Join(vals, ", "))
				}
			} else {
				_, vals := log.ToFields()
				fmt.Printf("%v\n", strings.Join(vals, ", "))
			}
		}
	}
}
