package policy

import (
	"encoding/xml"
	"fmt"
	"github.com/adambaumeister/panfw-util/panos/api"
	"github.com/adambaumeister/panfw-util/panos/errors"
)

// GetRules retrieves a list of rules, present at xpath
func GetRules(fqdn string, apikey string, xpath []string) []*Rule {
	rq := api.NewXpathQuery()
	rq.EnableAuth(apikey)

	rq.SetXpath(xpath)
	rq.AddParam("type", "config")
	rq.SetPath(api.API_ROOT)
	rq.SetFqdn(fqdn)

	r := RuleResponse{}
	resp := rq.Send()
	xml.Unmarshal(resp, &r)

	errors.LogDebug(string(resp))
	if len(r.Result.Rules.Entries) == 0 {
		r := EntryResponse{}
		xml.Unmarshal(resp, &r)
		return r.Entries
	}

	return r.Result.Rules.Entries
}

type RuleResponse struct {
	Status string   `xml:"status,attr"`
	Result Security `xml:"result>security"`
}

// Special response type, for when we are asking for a specific item
type EntryResponse struct {
	Entries []*Rule `xml:"result>entry"`
}

type Security struct {
	Rules Rules `xml:"rules"`
}

type Rules struct {
	Entries []*Rule `xml:"entry"`
}

type Rule struct {
	Name        string        `xml:"name,attr"`
	To          []MemberField `xml:"to"`
	From        []MemberField `xml:"From"`
	Source      []MemberField `xml:"source"`
	Destination []MemberField `xml:"destination"`
	SourceUser  []MemberField `xml:"source-user"`
	Category    []MemberField `xml:"category"`
	Application []MemberField `xml:"application"`
	Service     []MemberField `xml:"service"`
	HipProfiles []MemberField `xml:"hip-profiles"`
	Action      string        `xml:"action"`
	LogStart    string        `xml:"log-start"`
	LogEnd      string        `xml:"log-end"`
	Description string        `xml:"description"`
}

type MemberField struct {
	Member string `xml:"member"`
}

func (r *Rule) Print() {
	fmt.Printf("%v, %v\n", r.Name, r.Description)
}
