package testcmd

import (
	"encoding/xml"
	"github.com/adambaumeister/panfw-util/panos/api"
	"github.com/adambaumeister/panfw-util/panos/errors"
	"github.com/adambaumeister/panfw-util/pcaptest"
)

type TestSecurityPolicyMatchCmd struct {
	XMLName xml.Name      `xml:"test"`
	Flow    pcaptest.Flow `xml:"security-policy-match"`
}

func TestPolicy(fqdn string, apikey string, flow pcaptest.Flow) {
	cmd := TestSecurityPolicyMatchCmd{
		Flow: flow,
	}
	q := api.NewCmd(cmd)

	q.EnableAuth(apikey)
	q.SetFqdn(fqdn)
	q.SetPath(api.API_ROOT)

	q.AddParam("type", "op")
	resp := q.Send()
	errors.LogDebug(string(resp))
}
