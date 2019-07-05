package device

import (
	"fmt"
	"github.com/adambaumeister/panfw-util/Input"
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

func (fw *Firewall) Print(t string) {
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
	}

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

func (fw *Firewall) Register(args []string) deviceconfig.MsgJobResponse {
	// All args are treated as Ip addresses except the last, which is considered the tag
	al := len(args)
	var ips []string
	ips = args[:al-1]
	tag := args[al-1]
	var entries []*object.UidEntry
	for _, ip := range ips {
		o := &object.UidEntry{
			Ip:   ip,
			Tags: []string{tag},
		}
		entries = append(entries, o)
	}

	return object.BulkRegister(fw.Fqdn, fw.Apikey, entries)
}
func (fw *Firewall) UnRegister(args []string) deviceconfig.MsgJobResponse {
	// All args are treated as Ip addresses except the last, which is considered the tag
	al := len(args)
	var ips []string
	ips = args[:al-1]
	tag := args[al-1]
	var entries []*object.UidEntry
	for _, ip := range ips {
		o := &object.UidEntry{
			Ip:   ip,
			Tags: []string{tag},
		}
		entries = append(entries, o)
	}

	return object.BulkUnRegister(fw.Fqdn, fw.Apikey, entries)
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
