package object

import (
	"encoding/xml"
	"fmt"
	"github.com/adambaumeister/panfw-util/panos/api"
	"net"
)

func GetAddresses(fqdn string, apikey string, xpath []string) []*Address {
	/*
	   Retrieve all the address objects at the given xpath
	*/
	q := api.NewXpathQuery()
	q.EnableAuth(apikey)

	q.SetXpath(xpath)
	q.AddParam("type", "config")
	q.SetPath(api.API_ROOT)
	q.SetFqdn(fqdn)

	r := AddressResponse{}
	resp := q.Send()
	xml.Unmarshal(resp, &r)

	for _, addr := range r.Result.Entries {
		_, ipnet, err := net.ParseCIDR(addr.Ip)
		if err != nil {
			ip := addr.Ip + "/32"
			_, ipnet, err = net.ParseCIDR(ip)
			addr.IpObject = ipnet
		} else {
			addr.IpObject = ipnet
		}

	}

	return r.Result.Entries
}

type AddressResponse struct {
	Status string           `xml:"status,attr"`
	Result AddressContainer `xml:"result>address"`
}
type AddressContainer struct {
	Entries []*Address `xml:"entry"`
}
type Address struct {
	Name     string `xml:"name,attr"`
	Ip       string `xml:"ip-netmask"`
	IpObject *net.IPNet
}

func (a *Address) Print() {
	fmt.Printf("%v, %v\n", a.Name, a.IpObject.String())
}
