package panorama

import (
	"encoding/xml"
	"fmt"
	"github.com/adambaumeister/panfw-util/panos/api"
)

func GetDeviceGroups(fqdn string, apikey string, xpath []string) []DeviceGroup {
	q := api.NewXpathQuery()
	q.EnableAuth(apikey)

	q.SetXpath(xpath)
	q.AddParam("type", "config")
	q.SetPath(api.API_ROOT)
	q.SetFqdn(fqdn)

	resp := q.Send()
	r := DgResponse{}
	xml.Unmarshal(resp, &r)
	fmt.Printf("%v\n", r.Groups[0].Name)
	return r.Groups
}

type DgResponse struct {
	Status string        `xml:"status,attr"`
	Groups []DeviceGroup `xml:"result>device-group>entry"`
}

type DeviceGroup struct {
	Name string `xml:"name,attr"`
}
