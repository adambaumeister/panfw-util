package logs

import (
	"encoding/xml"
	"fmt"
	"github.com/adambaumeister/panfw-util/panos/api"
	"github.com/adambaumeister/panfw-util/panos/api/deviceconfig"
	"github.com/adambaumeister/panfw-util/panos/errors"
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
	Path               string `xml:"path"`
	FullPath           string `xml:"full-path"`
	Serial             string `xml:"serial"`
	DeviceName         string `xml:"device_name"`
	Host               string `xml:"host"`
	Command            string `xml:"cmd"`
	Admin              string `xml:"admin"`
	Result             string `xml:"result"`
	AfterChangePreview string `xml:"after-change-preview"`
}

func (l *LogEntry) Print() {
	_, vals := l.ToFields()
	fmt.Printf("%v\n", vals)
}

func (l *LogEntry) ToFields() ([]string, []string) {
	/*
		Returns all of the fields for the log entry and the values associated with it based on the type
	*/

	if l.Type == "CONFIG" {
		FieldKeys := []string{
			"time",
			"admin",
			"command",
			"result",
			"path",
			"full-path",
		}
		FieldVals := []string{
			l.ReceiveTime,
			l.Admin,
			l.Command,
			l.Result,
			l.Path,
			l.FullPath,
		}
		return FieldKeys, FieldVals
	}

	FieldKeys := []string{
		"time",
		"path",
		"full-path",
	}
	FieldVals := []string{
		l.ReceiveTime,
		l.Src,
		l.Dst,
		l.Rule,
	}
	return FieldKeys, FieldVals
}

func Query(fqdn string, apikey string, query string, count int, logtype string) []*LogEntry {
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
		return nil
	}

	job := ShowQueryJob(fqdn, apikey, r.Job)
	//bar := progressbar.NewOptions(100, progressbar.OptionSetRenderBlankState(true))
	for job.Status == "ACT" {
		job = ShowQueryJob(fqdn, apikey, r.Job)
		//bar.Add(job.Log.Progress)
		time.Sleep(1 * time.Second)
	}
	//bar.Finish()
	//print("\n")
	return job.Log.Entries
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
