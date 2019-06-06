package device

import (
	"fmt"
	"github.com/adambaumeister/panfw-util/panos/api/panorama"
)

type Panorama struct {
	Universal

	DeviceGroups []DeviceGroup
}

func (p *Panorama) Print(t string) {
	p.InitDeviceGroups()
}

func (p *Panorama) InitDeviceGroups() {
	xps := p.PrepQuery()
	xps = append(xps, "device-group")
	panorama.GetDeviceGroups(p.Fqdn, p.Apikey, xps)
}

func (p *Panorama) PrepQuery() []string {
	device := fmt.Sprintf("entry[@name='%v']", p.Device)

	xps := []string{
		DEVICE_XPATH,
		device,
	}
	return xps
}
