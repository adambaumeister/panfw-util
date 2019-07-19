package auth

import (
	"encoding/xml"
	"fmt"
	"github.com/adambaumeister/panfw-util/panos/api"
	"github.com/adambaumeister/panfw-util/panos/errors"
)

// Generate an API key for the given user.
func KeyGen(user string, pass string, fqdn string) string {

	kgq := api.NewParamQuery()
	kgq.SetFqdn(fqdn)
	kgq.SetPath(api.API_ROOT)

	kgq.AddParam("type", "keygen")
	kgq.AddParam("user", user)
	kgq.AddParam("password", pass)
	resp := kgq.Send()

	v := KeyGenResponse{}
	xml.Unmarshal(resp, &v)
	if v.Status == "error" {
		fmt.Printf("Authentication failed.")
		return ""
	}
	errors.LogDebug(fmt.Sprintf("Query status: %v, Key: %v\n", v.Status, v.Result.Key))

	return v.Result.Key
}

type KeyGenResponse struct {
	Status string `xml:"status,attr"`
	Result Key    `xml:"result"`
}

type Key struct {
	Key string `xml:"key"`
}
