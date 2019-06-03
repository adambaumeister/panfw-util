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
	errors.LogDebug(r.Msg)
	errors.LogDebug(string(resp))

}

type MsgJobResponse struct {
	api.Response
	Msg string `xml:"msg"`
	Job int    `xml:"job"`
}

type CommitCommand struct {
	XMLName xml.Name
}
