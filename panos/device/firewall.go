package device

import (
	"github.com/adamb/panfw-util/panos/api/auth"
)

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
