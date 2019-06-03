package device

import (
	"fmt"
	"github.com/adamb/panfw-util/panos/api/auth"
	"github.com/adamb/panfw-util/panos/api/deviceconfig"
	"github.com/adamb/panfw-util/panos/api/policy"
)

const DEVICE_XPATH = "/config/devices"

type Firewall struct {
	// Basic stuff, required
	Fqdn   string
	User   string
	Pass   string
	Apikey string

	// Optional
	Vsys   string
	Device string
}

func Connect(user string, pass string, fqdn string) *Firewall {
	/*
		Connect to a Firewall and return it's containing Struct
	*/
	fw := Firewall{
		Fqdn:   fqdn,
		Vsys:   "vsys1",
		Device: "localhost.localdomain",
	}
	fw.Apikey = auth.KeyGen(user, pass, fqdn)
	fw.User = user
	fw.Pass = pass
	return &fw
}

func (fw *Firewall) Rules() {
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

func (fw *Firewall) Load(fn string) {
	deviceconfig.Load(fw.Fqdn, fw.Apikey, fn, false)
}
