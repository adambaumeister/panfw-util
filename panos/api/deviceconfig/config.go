package deviceconfig

import (
	"fmt"
	"github.com/adamb/panfw-util/panos/api"
)

func Load(fqdn string, apikey string, fn string, commit bool) {
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
	fmt.Printf("%v\n", string(resp))
}
