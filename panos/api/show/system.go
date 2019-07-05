package show

import (
	"encoding/xml"
	"fmt"
	"github.com/adambaumeister/panfw-util/panos/api"
	"github.com/adambaumeister/panfw-util/panos/api/object"
	"github.com/adambaumeister/panfw-util/panos/errors"
)

func ShowSystemInfo(fqdn string, apikey string) SystemInfo {
	cmd := showSysInfo{
		XMLName: xml.Name{Local: "show"},
	}
	q := api.NewCmd(cmd)

	q.EnableAuth(apikey)
	q.SetFqdn(fqdn)
	q.SetPath(api.API_ROOT)

	q.AddParam("type", "op")
	resp := q.Send()
	errors.LogDebug(string(resp))

	//fmt.Printf(string(resp))
	r := SystemInfoResponse{}
	xml.Unmarshal(resp, &r)
	return r.SystemInfo
}

type showSysInfo struct {
	XMLName xml.Name
	Id      string `xml:"system>info"`
}

type SystemInfoResponse struct {
	api.Response
	SystemInfo SystemInfo `xml:"result>system"`
}

type SystemInfo struct {
	Devicename string `xml:"devicename"`
	Model      string `xml:"model"`
}

type showRegIps struct {
	XMLName xml.Name `xml:"show"`
	Cmd     string   `xml:"object>registered-ip>all"`
}

func ShowRegisteredIPs(fqdn string, apikey string) *RegisteredIpResponse {
	cmd := showRegIps{}
	q := api.NewCmd(cmd)

	q.EnableAuth(apikey)
	q.SetFqdn(fqdn)
	q.SetPath(api.API_ROOT)

	q.AddParam("type", "op")
	resp := q.Send()
	errors.LogDebug(string(resp))

	//fmt.Printf(string(resp))
	r := RegisteredIpResponse{}
	xml.Unmarshal(resp, &r)
	return &r
}

type RegisteredIpResponse struct {
	Entries []*object.UidEntry `xml:"result>entry"`
}

func (r *RegisteredIpResponse) Print() {
	for _, e := range r.Entries {
		fmt.Printf("%v\n", e.Ip)
		for _, t := range e.Tags {
			fmt.Printf("  %v\n", t)
		}
	}

}
