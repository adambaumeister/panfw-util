package auth

import (
	"encoding/xml"
	"fmt"
	"github.com/adamb/panfw-util/panos/api"
	"github.com/adamb/panfw-util/panos/errors"
)

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
