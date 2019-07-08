package logs

import (
	"encoding/xml"
	"fmt"
	"github.com/adambaumeister/panfw-util/panos/api"
	"github.com/adambaumeister/panfw-util/panos/api/deviceconfig"
	"github.com/adambaumeister/panfw-util/panos/errors"
	"github.com/schollz/progressbar"
	"time"
)

type LogResponse struct {
	Count    int         `xml:"count,attr"`
	Progress int         `xml:"progress,attr"`
	Entries  []*LogEntry `xml:"entry"`
}
type LogEntry struct {
	Type        string `xml:"type"`
	ReceiveTime string `xml:"receive_time"`
	// Traffic log fieldsd
	Rule  string `xml:"rule"`
	Src   string `xml:"src"`
	Dst   string `xml:"dst"`
	Sport string `xml:"sport"`
	Dport string `xml:"dport"`
	// Config log fields
	Path string `xml:"path"`
}

func (l *LogEntry) Print() {
	switch l.Type {
	case "TRAFFIC":
		fmt.Printf("%v: %v, %v:%v, %v:%v\n", l.ReceiveTime, l.Rule, l.Src, l.Sport, l.Dst, l.Dport)
	case "CONFIG":
		fmt.Printf("%v: %v\n", l.ReceiveTime, l.Path)
	}
}

func Query(fqdn string, apikey string, query string, count int, logtype string) {
	q := api.NewParamQuery()
	q.EnableAuth(apikey)

	q.AddParam("type", "log")
	q.AddParam("log-type", logtype)
	q.AddParam("nlogs", fmt.Sprintf("%v", count))
	if query != "" {
		q.AddParam("query", query)
	}

	q.SetPath(api.API_ROOT)
	q.SetFqdn(fqdn)

	resp := q.Send()
	errors.LogDebug(string(resp))
	r := deviceconfig.MsgJobResponse{}
	xml.Unmarshal(resp, &r)

	if r.Job == 0 {
		errors.LogDebug(r.Msg)
		return
	}

	job := ShowQueryJob(fqdn, apikey, r.Job)
	bar := progressbar.NewOptions(100, progressbar.OptionSetRenderBlankState(true))
	for job.Status == "ACT" {
		job = ShowQueryJob(fqdn, apikey, r.Job)
		bar.Add(job.Log.Progress)
		time.Sleep(1 * time.Second)
	}
	bar.Finish()
	print("\n")
	for _, e := range job.Log.Entries {
		e.Print()
	}
}
func ShowQueryJob(fqdn string, apikey string, jobid int) *QueryJobResult {

	jid := fmt.Sprintf("%v", jobid)
	//fmt.Printf("JOBID: %v\n", jid)
	q := api.NewParamQuery()
	q.EnableAuth(apikey)

	q.SetPath(api.API_ROOT)
	q.SetFqdn(fqdn)

	q.AddParam("type", "log")
	q.AddParam("action", "get")
	q.AddParam("job-id", jid)

	resp := q.Send()
	errors.LogDebug(string(resp))

	r := QueryJobResult{}
	xml.Unmarshal(resp, &r)

	return &r
}

type QueryJobResult struct {
	deviceconfig.Job `xml:"result>job"`
	Log              LogResponse `xml:"result>log>logs"`
}
