package device

import (
	"fmt"
	"github.com/adambaumeister/panfw-util/panos/api/auth"
	"github.com/adambaumeister/panfw-util/panos/api/deviceconfig"
	"github.com/adambaumeister/panfw-util/panos/api/show"
)

/*
Universal represents functionality that is common to both physical firewalls and panorama
*/
type Universal struct {
	Fqdn   string
	User   string
	Pass   string
	Apikey string

	Device string
}

/*
Generic PANOS interface

Functions that interact with any PAN-OS device.
*/
type Panos interface {
	Print(string)
	Add([]string)

	ImportNamed(string)
	LoadNamed(string, bool)

	SetDeviceGroup(string)
}

func ConnectUniversal(user string, pass string, fqdn string) Panos {
	/*
		Connect to a Panos device and return it
	*/
	// We use Firewall as the container as it works regardless of the underyling type
	fw := Firewall{
		Vsys: "vsys1",
	}
	fw.Device = "localhost.localdomain"
	fw.Fqdn = fqdn
	fw.Apikey = auth.KeyGen(user, pass, fqdn)
	fw.User = user
	fw.Pass = pass

	si := show.ShowSystemInfo(fw.Fqdn, fw.Apikey)
	// Switch based on the model
	if si.Model == "Panorama" {
		p := Panorama{}
		p.Apikey = fw.Apikey
		p.Fqdn = fw.Fqdn
		p.Device = fw.Device
		// If it's panorama, get the device groups.
		p.InitDeviceGroups()
		p.CurrentDeviceGroup = "shared"
		return &p
	}
	return &fw
}

func (fw *Universal) ImportNamed(fn string) {
	/*
		Import a named configuration file

		Files are imported as the name as they appear on disk
	*/
	fmt.Printf("Importing named configuration file %v...", fn)
	r := deviceconfig.ImportNamed(fw.Fqdn, fw.Apikey, fn)
	if r.Status == "success" {
		fmt.Printf("Done!\n")
	} else {
		fmt.Printf("Failed!\n")
	}
	return
}

func (fw *Universal) LoadNamed(fn string, commit bool) {
	fmt.Printf("Loading named configuration file %v...\n", fn)
	deviceconfig.LoadNamedConfig(fw.Fqdn, fw.Apikey, fn, commit)
}
