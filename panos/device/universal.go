package device

import (
	"fmt"
	"github.com/adamb/panfw-util/panos/api/deviceconfig"
)

/*
Universal represents functionality that is common to both physical firewalls and panorama
*/
type Universal struct {
	Fqdn   string
	User   string
	Pass   string
	Apikey string
}

func (fw *Universal) ImportNamed(fn string) {
	/*
		Import a named configuration file

		Files are imported as the name as they appear on disk
	*/
	r := deviceconfig.ImportNamed(fw.Fqdn, fw.Apikey, fn)
	fmt.Printf("Import complete!\n")
	r.Print()
	return
}

func (fw *Universal) LoadNamed(fn string, commit bool) {
	deviceconfig.LoadNamedConfig(fw.Fqdn, fw.Apikey, fn, commit)
}
