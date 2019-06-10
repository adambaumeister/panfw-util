package Input

import (
	"encoding/csv"
	"fmt"
	"github.com/adambaumeister/panfw-util/panos/api/object"
	"github.com/adambaumeister/panfw-util/panos/errors"
	"io/ioutil"
	"os"
	"strings"
)

func ToObjects(args []string) []object.ApiObject {
	s := args[0]
	errors.LogDebug(fmt.Sprintf("Converting:\n%v\nto objects", s))
	_, err := os.Stat(s)
	if err == nil {
		bytes, err := ioutil.ReadFile(s)
		errors.DieIf(err)
		return CsvToObjects(string(bytes))
	}

	// If the argument is a CSV
	if len(strings.Split(s, ",")) > 1 {
		errors.LogDebug("Input argument taken as a CSV")
		return CsvToObjects(s)
	}

	return nil
}

func CsvToObjects(s string) []object.ApiObject {
	var objs []object.ApiObject

	fmt.Printf("DEBUG: %v\n", s)
	r := csv.NewReader(strings.NewReader(s))
	str_csv, err := r.ReadAll()
	errors.DieIf(err)

	for _, str_csv_entry := range str_csv {
		//fmt.Printf("DEBUG: entry: %v\n", str_csv_entry)
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
	}
	return o
}
