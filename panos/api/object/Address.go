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

	//xpath = append(xpath, fmt.Sprintf("entry[@name='%v']", addr.Name))

	xaddr, err := xml.Marshal(addr)
	errors.LogDebug(string(xaddr))
	errors.DieIf(err)
	//xaddr := fmt.Sprintf("<ip-netmask>%v</ip-netmask>", addr.Ip)

	q := api.NewXpathQuery()
	q.EnableAuth(apikey)
	print(api.MakeXPath(xpath))
	fmt.Print("\n")

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
