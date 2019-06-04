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
