package device

import (
	"fmt"
	"github.com/adambaumeister/panfw-util/panos/api"
	"github.com/adambaumeister/panfw-util/panos/api/panorama"
	"github.com/adambaumeister/panfw-util/panos/errors"
)

type Panorama struct {
	Universal

	DeviceGroups       []DeviceGroup
	CurrentDeviceGroup string
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

func (p *Panorama) Add(args []string) {
	dg := p.getDG(p.CurrentDeviceGroup)
	dg.Add(args)
}

func (p *Panorama) SetDeviceGroup(name string) {
	errors.LogDebug(fmt.Sprintf("Setting active DG to %v", name))
	p.CurrentDeviceGroup = name
}

func (p *Panorama) getDG(name string) *DeviceGroup {
	for _, dg := range p.DeviceGroups {
		if dg.Name == name {
			errors.LogDebug(fmt.Sprintf("Found DG in Panorama: %v", name))
			return &dg
		}
	}

	dgobj := DeviceGroup{
		Name:   "shared",
		parent: p,
	}
	return &dgobj
}
