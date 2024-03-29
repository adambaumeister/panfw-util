package device

import (
	"fmt"
	"github.com/adambaumeister/panfw-util/Input"
	"github.com/adambaumeister/panfw-util/panos/api/object"
	"github.com/adambaumeister/panfw-util/panos/api/policy"
)

type DeviceGroup struct {
	parent *Panorama

	Name string
}

func (dg *DeviceGroup) Addresses() []*object.Address {
	/*
		Return all the Address objects
	*/

	xps := dg.PrepQuery()
	xps = append(xps, "address")
	// Important - use the parent connection details
	objs := object.GetAddresses(dg.parent.Fqdn, dg.parent.Apikey, xps)
	return objs
}

func (dg *DeviceGroup) AddressGroups() []*object.AddressGroup {
	/*
		Return all the Address objects
	*/

	xps := dg.PrepQuery()
	xps = append(xps, "address-group")
	// Important - use the parent connection details
	objs := object.GetAddressGroups(dg.parent.Fqdn, dg.parent.Apikey, xps)
	return objs
}

func (dg *DeviceGroup) Services() []*object.Service {
	xps := dg.PrepQuery()
	xps = append(xps, "service")
	// Important - use the parent connection details
	objs := object.GetServices(dg.parent.Fqdn, dg.parent.Apikey, xps)
	return objs
}

func (dg *DeviceGroup) Rules() []*policy.Rule {
	xps := dg.PrepQuery()
	xps = append(xps, "pre-rulebase")
	xps = append(xps, "security")
	// Important - use the parent connection details
	objs := policy.GetRules(dg.parent.Fqdn, dg.parent.Apikey, xps)
	return objs
}

func (dg *DeviceGroup) Add(args []string) {
	objs := Input.ToObjects(args)

	for _, ob := range objs {
		xps := dg.PrepQuery()
		xps = append(xps, ob.GetType())
		msg := ob.Add(dg.parent.Fqdn, dg.parent.Apikey, xps)
		fmt.Printf("Add %v : %v\n", ob.GetName(), msg.Status)
	}
}

func (dg *DeviceGroup) PrepQuery() []string {

	// Hacking around the fact that "shared" is not a normal DG
	if dg.Name == "shared" {
		xps := []string{
			"/config/shared",
		}
		return xps
	}
	device := fmt.Sprintf("entry[@name='%v']", dg.parent.Device)
	group := fmt.Sprintf("device-group/entry[@name='%v']", dg.Name)

	xps := []string{
		DEVICE_XPATH,
		device,
		group,
	}
	return xps
}
