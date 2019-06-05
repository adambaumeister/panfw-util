package device

import (
	"fmt"
	"github.com/adambaumeister/panfw-util/panos/api"
	"github.com/adambaumeister/panfw-util/panos/api/auth"
	"github.com/adambaumeister/panfw-util/panos/api/deviceconfig"
	"github.com/adambaumeister/panfw-util/panos/api/object"
	"github.com/adambaumeister/panfw-util/panos/api/policy"
)

const DEVICE_XPATH = "/config/devices"

// Firewall represents a physical or virtual PANOS firewall
type Firewall struct {
	Universal

	// Optional
	Vsys   string
	Device string
}

func Connect(user string, pass string, fqdn string) *Firewall {
	/*
		Connect to a Firewall and return it's containing Struct
	*/
	fw := Firewall{
		Vsys:   "vsys1",
		Device: "localhost.localdomain",
	}
	fw.Fqdn = fqdn
	fw.Apikey = auth.KeyGen(user, pass, fqdn)
	fw.User = user
	fw.Pass = pass
	return &fw
}

func (fw *Firewall) Print(t string) {
	var objs []api.Entry
	switch t {
	case "address":
		// Not sure why this is required, probably golang idiosyncrasy
		for _, a := range fw.Addresses() {
			objs = append(objs, a)
		}
	}

	for _, o := range objs {
		o.Print()
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
