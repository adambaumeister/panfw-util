package object

import (
	"encoding/xml"
	"fmt"
	"github.com/adambaumeister/panfw-util/panos/api"
	"github.com/adambaumeister/panfw-util/panos/api/deviceconfig"
	"github.com/adambaumeister/panfw-util/panos/errors"
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
			addr.ipObject = ipnet
		} else {
			addr.ipObject = ipnet
		}

	}

	return r.Result.Entries
}

func (addr *Address) Add(fqdn string, apikey string, xpath []string) deviceconfig.MsgJobResponse {
	/*
		Add an address object to fqdn at xpath location

		Returns a msgJobResponse object containing the status
	*/

	xaddr, err := xml.Marshal(addr)
	errors.LogDebug(string(xaddr))
	errors.DieIf(err)

	q := api.NewXpathQuery()
	q.EnableAuth(apikey)

	q.SetXpath(xpath)
	q.AddParam("type", "config")
	q.AddParam("action", "set")
	q.AddParam("element", string(xaddr))
	q.SetPath(api.API_ROOT)
	q.SetFqdn(fqdn)

	resp := q.Send()
	r := deviceconfig.MsgJobResponse{}
	xml.Unmarshal(resp, &r)
	return r
}

type AddressResponse struct {
	Status string           `xml:"status,attr"`
	Result AddressContainer `xml:"result>address"`
}
type AddressContainer struct {
	Entries []*Address `xml:"entry"`
}

type Entries struct {
	Entries []*Address `xml:"entry"`
}
type Address struct {
	XMLName  xml.Name `xml:"entry"`
	Name     string   `xml:"name,attr"`
	Ip       string   `xml:"ip-netmask"`
	ipObject *net.IPNet
}

func (a *Address) Print() {
	fmt.Printf("%v, %v\n", a.Name, a.ipObject.String())
}

func (a *Address) GetType() string {
	return "address"
}

func (a *Address) GetName() string {
	return a.Name
}

func GetAddressGroups(fqdn string, apikey string, xpath []string) []*AddressGroup {
	/*
	   Retrieve all the address objects at the given xpath
	*/
	q := api.NewXpathQuery()
	q.EnableAuth(apikey)

	q.SetXpath(xpath)
	q.AddParam("type", "config")
	q.SetPath(api.API_ROOT)
	q.SetFqdn(fqdn)

	r := AddressGroupResponse{}
	resp := q.Send()

	xml.Unmarshal(resp, &r)

	errors.LogDebug(string(resp))
	errors.LogDebug(r.Status)
	return r.Entries
}

type AddressGroupResponse struct {
	Status  string          `xml:"status,attr"`
	Entries []*AddressGroup `xml:"result>address-group>entry"`
}

type AddressGroupEntries struct {
	XMLName xml.Name        `xml:"address-group"`
	Entries []*AddressGroup `xml:"entry"`
}

type AddressGroup struct {
	XMLName       xml.Name `xml:"entry"`
	Name          string   `xml:"name,attr"`
	StaticMembers []string `xml:"static>member"`
}

func (ag *AddressGroup) Print() {
	fmt.Printf("%v\n", ag.Name)
	for _, entry := range ag.StaticMembers {
		fmt.Printf(" %v\n", entry)
	}
}

func (ag *AddressGroup) GetType() string {
	return "address-group"
}

func (a *AddressGroup) GetName() string {
	return a.Name
}

func (ag *AddressGroup) Add(fqdn string, apikey string, xpath []string) deviceconfig.MsgJobResponse {
	xaddr, err := xml.Marshal(ag)
	errors.LogDebug(api.MakeXPath(xpath))
	errors.LogDebug(string(xaddr))
	errors.DieIf(err)

	q := api.NewXpathQuery()
	q.EnableAuth(apikey)

	q.SetXpath(xpath)
	q.AddParam("type", "config")
	q.AddParam("action", "set")
	q.AddParam("element", string(xaddr))
	q.SetPath(api.API_ROOT)
	q.SetFqdn(fqdn)

	resp := q.Send()
	r := deviceconfig.MsgJobResponse{}
	xml.Unmarshal(resp, &r)
	errors.LogDebug(string(resp))
	return r
}
