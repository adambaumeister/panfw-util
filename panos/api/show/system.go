package show

import (
	"encoding/xml"
	"github.com/adambaumeister/panfw-util/panos/api"
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
