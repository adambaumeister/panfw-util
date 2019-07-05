package object

import (
	"encoding/xml"
	"github.com/adambaumeister/panfw-util/panos/api"
	"github.com/adambaumeister/panfw-util/panos/api/deviceconfig"
	"github.com/adambaumeister/panfw-util/panos/errors"
)

type UidMessage struct {
	XMLName xml.Name `xml:"uid-message"`
	Version string   `xml:"version"`
	Type    string   `xml:"type"`
	Payload Payload  `xml:"payload"`
}

type Payload struct {
	Register   *Register `xml:"register"`
	Unregister *Register `xml:"unregister"`
}

type Register struct {
	UidEntrys []*UidEntry `xml:"entry"`
}

type UidEntry struct {
	Ip   string   `xml:"ip,attr"`
	Tags []string `xml:"tag>member"`
}

func BulkRegister(fqdn string, apikey string, entries []*UidEntry) deviceconfig.MsgJobResponse {
	// Convienience function for registering a bulk of UidEntry objects as dynamic addresses
	register := Register{
		UidEntrys: entries,
	}
	payload := Payload{
		Register: &register,
	}

	message := UidMessage{
		Version: "1.0",
		Type:    "update",
		Payload: payload,
	}
	xaddr, err := xml.Marshal(message)
	errors.LogDebug(string(xaddr))
	errors.DieIf(err)

	q := api.NewParamQuery()
	q.EnableAuth(apikey)

	q.AddParam("type", "user-id")
	q.AddParam("cmd", string(xaddr))
	q.SetPath(api.API_ROOT)
	q.SetFqdn(fqdn)

	resp := q.Send()
	errors.LogDebug(string(resp))
	r := deviceconfig.MsgJobResponse{}
	xml.Unmarshal(resp, &r)
	return r
}

func BulkUnRegister(fqdn string, apikey string, entries []*UidEntry) deviceconfig.MsgJobResponse {
	// Convienience function for unregistering a bulk of UidEntry objects as dynamic addresses
	register := Register{
		UidEntrys: entries,
	}
	payload := Payload{
		Unregister: &register,
	}

	message := UidMessage{
		Version: "1.0",
		Type:    "update",
		Payload: payload,
	}
	xaddr, err := xml.Marshal(message)
	errors.LogDebug(string(xaddr))
	errors.DieIf(err)

	q := api.NewParamQuery()
	q.EnableAuth(apikey)

	q.AddParam("type", "user-id")
	q.AddParam("cmd", string(xaddr))
	q.SetPath(api.API_ROOT)
	q.SetFqdn(fqdn)

	resp := q.Send()
	errors.LogDebug(string(resp))
	r := deviceconfig.MsgJobResponse{}
	xml.Unmarshal(resp, &r)
	return r
}

func (entry *UidEntry) Register(fqdn string, apikey string) deviceconfig.MsgJobResponse {
	// Register a single dynamic address entry
	entries := []*UidEntry{entry}
	register := Register{
		UidEntrys: entries,
	}
	payload := Payload{
		Register: &register,
	}

	message := UidMessage{
		Version: "1.0",
		Type:    "update",
		Payload: payload,
	}
	xaddr, err := xml.Marshal(message)
	errors.LogDebug(string(xaddr))
	errors.DieIf(err)

	q := api.NewParamQuery()
	q.EnableAuth(apikey)

	q.AddParam("type", "user-id")
	q.AddParam("cmd", string(xaddr))
	q.SetPath(api.API_ROOT)
	q.SetFqdn(fqdn)

	resp := q.Send()
	errors.LogDebug(string(resp))
	r := deviceconfig.MsgJobResponse{}
	xml.Unmarshal(resp, &r)
	return r
}

func (entry *UidEntry) UnRegister(fqdn string, apikey string) deviceconfig.MsgJobResponse {
	// Unregister a single UidEntry as a dynamic address object
	entries := []*UidEntry{entry}
	register := Register{
		UidEntrys: entries,
	}
	payload := Payload{
		Unregister: &register,
	}

	message := UidMessage{
		Version: "1.0",
		Type:    "update",
		Payload: payload,
	}
	xaddr, err := xml.Marshal(message)
	errors.LogDebug(string(xaddr))
	errors.DieIf(err)

	q := api.NewParamQuery()
	q.EnableAuth(apikey)

	q.AddParam("type", "user-id")
	q.AddParam("cmd", string(xaddr))
	q.SetPath(api.API_ROOT)
	q.SetFqdn(fqdn)

	resp := q.Send()
	errors.LogDebug(string(resp))
	r := deviceconfig.MsgJobResponse{}
	xml.Unmarshal(resp, &r)
	return r
}
