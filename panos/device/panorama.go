package device

import (
	"fmt"
	"github.com/adambaumeister/panfw-util/panos/api"
	"github.com/adambaumeister/panfw-util/panos/api/panorama"
)

type Panorama struct {
	Universal

	DeviceGroups []DeviceGroup
}

func (p *Panorama) Print(t string) {
	/*
		Print a given type of object, like "address"
		For panorama, this will iterate through all known object groups.
	*/
	var objs []api.Entry
	switch t {
	case "address":
		// Not sure why this is required, probably golang idiosyncrasy
		for _, dg := range p.DeviceGroups {
			for _, a := range dg.Addresses() {
				objs = append(objs, a)
			}
		}
	}
	for _, o := range objs {
		o.Print()
	}
}

func (p *Panorama) InitDeviceGroups() {
	xps := p.PrepQuery()
	xps = append(xps, "device-group")
	for _, dg := range panorama.GetDeviceGroups(p.Fqdn, p.Apikey, xps) {
		dgobj := DeviceGroup{
			Name:   dg.Name,
			parent: p,
		}
		p.DeviceGroups = append(p.DeviceGroups, dgobj)
	}
}

func (p *Panorama) PrepQuery() []string {
	device := fmt.Sprintf("entry[@name='%v']", p.Device)

	xps := []string{
		DEVICE_XPATH,
		device,
	}
	return xps
}
