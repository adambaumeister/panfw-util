package device

import (
	"fmt"
	"github.com/adambaumeister/panfw-util/Input"
	"github.com/adambaumeister/panfw-util/panos/api"
	"github.com/adambaumeister/panfw-util/panos/api/auth"
	"github.com/adambaumeister/panfw-util/panos/api/deviceconfig"
	"github.com/adambaumeister/panfw-util/panos/api/object"
	"github.com/adambaumeister/panfw-util/panos/api/policy"
	"github.com/adambaumeister/panfw-util/panos/api/show"
)

const DEVICE_XPATH = "/config/devices"

// Firewall represents a physical or virtual PANOS firewall
type Firewall struct {
	Universal

	// Optional
	Vsys string
}

func Connect(user string, pass string, fqdn string) *Firewall {
	/*
		Connect to a Firewall and return it's containing Struct
	*/
	fw := Firewall{
		Vsys: "vsys1",
	}
	fw.Device = "localhost.localdomain"
	fw.Fqdn = fqdn
	fw.Apikey = auth.KeyGen(user, pass, fqdn)
	fw.User = user
	fw.Pass = pass
	return &fw
}

func (fw *Firewall) Get(t string) []api.Entry {
	var objs []api.Entry
	switch t {
	case "address":
		// Not sure why this is required, probably golang idiosyncrasy
		for _, a := range fw.Addresses() {
			objs = append(objs, a)
		}
	case "address-group":
		for _, a := range fw.AddressGroups() {
			objs = append(objs, a)
		}
	case "registered-ips":
		objs = append(objs, show.ShowRegisteredIPs(fw.Fqdn, fw.Apikey))
	case "?":
		fmt.Printf("Available options:\n")
		fmt.Printf(" address\n")
		fmt.Printf(" address-group\n")
		fmt.Printf(" registered-ips\n")
	}

	return objs
}

func (fw *Firewall) Print(t string) {
	objs := fw.Get(t)

	for _, o := range objs {
		o.Print()
	}
}

func (fw *Firewall) Add(args []string) {
	objs := Input.ToObjects(args)

	for _, ob := range objs {
		xps := fw.PrepQuery()
		xps = append(xps, ob.GetType())
		ob.Add(fw.Fqdn, fw.Apikey, xps)
	}
}

func (fw *Firewall) Rules() {
	/*
		Return the firewall rulebase
	*/
	xps := fw.PrepQuery()
	xps = append(xps, "rulebase")
	xps = append(xps, "security")

	policy.GetRules(fw.Fqdn, fw.Apikey, xps)
}

func (fw *Firewall) Addresses() []*object.Address {
	/*
		Return all the Address objects
	*/

	xps := fw.PrepQuery()
	xps = append(xps, "address")
	objs := object.GetAddresses(fw.Fqdn, fw.Apikey, xps)
	return objs
}

func (fw *Firewall) AddressGroups() []*object.Address {
	/*
		Return all the Address objects
	*/

	xps := fw.PrepQuery()
	xps = append(xps, "address-group")
	objs := object.GetAddresses(fw.Fqdn, fw.Apikey, xps)
	return objs
}

func (fw *Firewall) Commit() {
	deviceconfig.Commit(fw.Fqdn, fw.Apikey)
}

func (fw *Firewall) PrepQuery() []string {
	device := fmt.Sprintf("entry[@name='%v']", fw.Device)
	vsys := fmt.Sprintf("vsys/entry[@name='%v']", fw.Vsys)

	xps := []string{
		DEVICE_XPATH,
		device,
		vsys,
	}
	return xps
}

func (fw *Firewall) SetDeviceGroup(s string) {
	// This function does nothing as Firewall devices do not have devicegroups
}
