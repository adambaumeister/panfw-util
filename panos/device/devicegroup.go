package device

import (
	"fmt"
	"github.com/adambaumeister/panfw-util/panos/api/object"
)

type DeviceGroup struct {
	Universal

	Name string
}

func (dg *DeviceGroup) Addresses() []*object.Address {
	/*
		Return all the Address objects
	*/

	xps := dg.PrepQuery()
	xps = append(xps, "address")
	objs := object.GetAddresses(dg.Fqdn, dg.Apikey, xps)
	return objs
}

func (dg *DeviceGroup) PrepQuery() []string {
	device := fmt.Sprintf("entry[@name='%v']", dg.Device)
	group := fmt.Sprintf("device-group/entry[@name='%v']", dg.Name)

	xps := []string{
		DEVICE_XPATH,
		device,
		group,
	}
	return xps
}
