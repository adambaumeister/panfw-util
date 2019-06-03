package device

import (
	"fmt"
	"github.com/adamb/panfw-util/panos/api/auth"
	"github.com/adamb/panfw-util/panos/api/policy"
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

func (fw *Firewall) Rules() {
	/*
		Return the firewall rulebase
	*/
	device := fmt.Sprintf("entry[@name='%v']", fw.Device)
	vsys := fmt.Sprintf("vsys/entry[@name='%v']", fw.Vsys)

	xps := []string{
		DEVICE_XPATH,
		device,
		vsys,
		"rulebase", "security",
	}

	policy.GetRules(fw.Fqdn, fw.Apikey, xps)
}
