package Input

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/adambaumeister/panfw-util/panos/api/object"
	"github.com/adambaumeister/panfw-util/panos/errors"
	"os"
	"strings"
)

func ToObjects(args []string) []object.ApiObject {
	s := args[0]
	errors.LogDebug(fmt.Sprintf("Converting:\n%v\nto objects", s))
	// If the argument is a CSV
	if len(strings.Split(s, ",")) > 1 {
		errors.LogDebug("Input argument taken as a CSV")
		return CsvToObjects(s)
	}

	if _, err := os.Stat(s); err == nil {
		errors.LogDebug("Input argument taken as a file")
		f, err := os.Open(s)
		errors.DieIf(err)
		buf := new(bytes.Buffer)
		buf.ReadFrom(f)
		csvString := buf.String()
		errors.LogDebug(csvString)
		return CsvToObjects(csvString)
	}

	return nil
}

func CsvToObjects(s string) []object.ApiObject {
	var objs []object.ApiObject

	r := csv.NewReader(strings.NewReader(s))
	str_csv, err := r.ReadAll()
	errors.DieIf(err)

	for _, str_csv_entry := range str_csv {
		o := ListToObjects(str_csv_entry)
		objs = append(objs, o)
	}

	return objs
}

func ListToObjects(vals []string) object.ApiObject {
	t := vals[0]
	var o object.ApiObject
	switch {
	case t == "address":
		name := vals[1]
		ip := vals[2]
		o = &object.Address{
			Name: name,
			Ip:   ip,
		}
	case t == "address-group":
		name := vals[1]
		member := vals[2]
		members := []string{member}
		o = &object.AddressGroup{
			Name:          name,
			StaticMembers: members,
		}
	case t == "service":
		name := vals[1]

		srvObj := &object.Service{
			Name: name,
		}
		protocol := vals[2]
		var pd object.PortDefinition
		if protocol == "tcp" {
			pd = object.PortDefinition{
				Port: vals[3],
			}
			srvObj.Tcp = &pd
		} else {
			pd = object.PortDefinition{
				Port: vals[3],
			}
			srvObj.Udp = &pd
		}
		o = srvObj
	}
	return o
}
