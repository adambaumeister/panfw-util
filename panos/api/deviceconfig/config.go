package deviceconfig

import (
	"encoding/xml"
	"fmt"
	"github.com/adamb/panfw-util/panos/api"
	"strings"
)

func Load(fqdn string, apikey string, fn string, commit bool) *MsgResponse {
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
