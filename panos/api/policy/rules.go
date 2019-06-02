package policy

import (
	"encoding/xml"
	"fmt"
	"github.com/adamb/panfw-util/panos/api"
)

// GetRules retrieves the rulebase
// This call is only valid on firewalls, Panorama requires setting a devicegroup
func GetRules(fqdn string, apikey string, xpath []string) {
	rq := api.NewXpathQuery()
	rq.EnableAuth(apikey)

	rq.SetXpath(xpath)
	rq.AddParam("type", "config")
	rq.SetPath(api.API_ROOT)
	rq.SetFqdn(fqdn)

	r := RuleResponse{}
	resp := rq.Send()
	xml.Unmarshal(resp, &r)
	fmt.Printf("Source: %v\n", r.Result.Rules.Entries[0].Name)
}

type RuleResponse struct {
	Status string   `xml:"status,attr"`
	Result Security `xml:"result>security"`
}

type Security struct {
	Rules Rules `xml:"rules"`
}

type Rules struct {
	Entries []Rule `xml:"entry"`
}

type Rule struct {
	Name   string `xml:"name,attr"`
	Source string `xml:"source"`
}
