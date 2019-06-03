package deviceconfig

import (
	"encoding/xml"
	"fmt"
	"github.com/adamb/panfw-util/panos/api"
	"github.com/adamb/panfw-util/panos/errors"
	"strings"
)

func ImportNamed(fqdn string, apikey string, fn string, commit bool) *MsgResponse {
	/*
		Imports a named configuration snapshot off the local disk
	*/
	// Read the file into multipart encoded data
	body, boundary := api.HttpMultiPart(fn)
	pq := api.NewPost(body)
	// Set the content-type in the header
	pq.SetMultipart(boundary)

	pq.EnableAuth(apikey)
	pq.SetFqdn(fqdn)
	pq.SetPath(api.API_ROOT)

	pq.AddParam("type", "import")
	pq.AddParam("category", "configuration")

	resp := pq.Send()
	r := MsgResponse{}
	xml.Unmarshal(resp, &r)
	return &r
}

type MsgResponse struct {
	api.Response
	Msg []string `xml:"msg>line"`
}

func (m *MsgResponse) Print() {
	fmt.Printf("%v\n", strings.Join(m.Msg, "\n"))
}

func Commit(fqdn string, apikey string) {
	c := CommitCommand{
		XMLName: xml.Name{Local: "commit"},
	}

	q := api.NewCmd(c)
	q.EnableAuth(apikey)
	q.SetFqdn(fqdn)
	q.SetPath(api.API_ROOT)

	q.AddParam("type", "commit")
	resp := q.Send()

	r := MsgJobResponse{}
	xml.Unmarshal(resp, &r)
	errors.LogDebug(string(resp))
	// If there are no changes
	if r.Job == 0 {
		errors.LogDebug(r.Msg)
		return
	}

	ShowJob(fqdn, apikey, r.Job)
}

func ShowJob(fqdn string, apikey string, jobid int) {
	cmd := showJobs{
		XMLName: xml.Name{Local: "show"},
		Id:      jobid,
	}
	q := api.NewCmd(cmd)

	q.EnableAuth(apikey)
	q.SetFqdn(fqdn)
	q.SetPath(api.API_ROOT)

	q.AddParam("type", "op")
	resp := q.Send()
	errors.LogDebug(string(resp))
}

type MsgJobResponse struct {
	api.Response
	Msg string `xml:"msg"`
	Job int    `xml:"result>job"`
}

type CommitCommand struct {
	XMLName xml.Name
}

type showJobs struct {
	XMLName xml.Name
	Id      int `xml:"jobs>id"`
}

func LoadNamedConfig(fqdn string, apikey string, cn string) {
	c := LoadNamedCommand{
		XMLName: xml.Name{Local: "load"},
		Config:  cn,
	}

	q := api.NewCmd(c)
	q.EnableAuth(apikey)
	q.SetFqdn(fqdn)
	q.SetPath(api.API_ROOT)
	q.AddParam("type", "op")

	resp := q.Send()
	errors.LogDebug(string(resp))
}

type LoadNamedCommand struct {
	XMLName xml.Name
	Config  string `xml:"config>from"`
}
