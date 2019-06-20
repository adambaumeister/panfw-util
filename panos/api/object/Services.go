package object

import (
	"encoding/xml"
	"fmt"
	"github.com/adambaumeister/panfw-util/panos/api"
	"github.com/adambaumeister/panfw-util/panos/api/deviceconfig"
	"github.com/adambaumeister/panfw-util/panos/errors"
)

func GetServices(fqdn string, apikey string, xpath []string) []*Service {
	/*
	   Retrieve all the address objects at the given xpath
	*/
	q := api.NewXpathQuery()
	q.EnableAuth(apikey)

	q.SetXpath(xpath)
	q.AddParam("type", "config")
	q.SetPath(api.API_ROOT)
	q.SetFqdn(fqdn)

	r := ServiceResponse{}
	resp := q.Send()

	xml.Unmarshal(resp, &r)
	errors.LogDebug(string(resp))

	return r.Entries
}

type ServiceResponse struct {
	Status  string     `xml:"status,attr"`
	Entries []*Service `xml:"result>service>entry"`
}

type Service struct {
	XMLName xml.Name        `xml:"entry"`
	Name    string          `xml:"name,attr"`
	Tcp     *PortDefinition `xml:"protocol>tcp,omitempty"`
	Udp     *PortDefinition `xml:"protocol>udp,omitempty"`
}

type PortDefinition struct {
	Port string `xml:"port"`
}

func (s *Service) Print() {
	if s.Tcp != nil {
		fmt.Printf("Name: %v Port: %v\n", s.Name, s.Tcp.Port)
	}

	if s.Udp != nil {
		fmt.Printf("Name: %v Port: %v\n", s.Name, s.Udp.Port)
	}
}

func (s *Service) GetType() string {
	return "service"
}

func (s *Service) GetName() string {
	return s.Name
}

func (s *Service) Add(fqdn string, apikey string, xpath []string) deviceconfig.MsgJobResponse {
	xservice, err := xml.Marshal(s)
	errors.LogDebug(api.MakeXPath(xpath))
	errors.LogDebug(string(xservice))
	errors.DieIf(err)

	q := api.NewXpathQuery()
	q.EnableAuth(apikey)

	q.SetXpath(xpath)
	q.AddParam("type", "config")
	q.AddParam("action", "set")
	q.AddParam("element", string(xservice))
	q.SetPath(api.API_ROOT)
	q.SetFqdn(fqdn)

	resp := q.Send()
	r := deviceconfig.MsgJobResponse{}
	xml.Unmarshal(resp, &r)
	errors.LogDebug(string(resp))
	return r
}
