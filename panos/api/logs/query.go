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

type LogEntry interface {
	Print()
}

type TrafficLogResponse struct {
	Count    int                `xml:"count,attr"`
	Progress int                `xml:"progress,attr"`
	Entries  []*TrafficLogEntry `xml:"entry"`
}
type TrafficLogEntry struct {
	Rule string `xml:"rule"`
	Src  string `xml:"src"`
	Dst  string `xml:"dst"`
}

func (l *TrafficLogEntry) Print() {
	fmt.Printf("%v, %v, %v\n", l.Rule, l.Src, l.Dst)
}

func Query(fqdn string, apikey string, query string) {
	q := api.NewParamQuery()
	q.EnableAuth(apikey)

	q.AddParam("type", "log")
	q.AddParam("log-type", "traffic")
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
	Log              TrafficLogResponse `xml:"result>log>logs"`
}
