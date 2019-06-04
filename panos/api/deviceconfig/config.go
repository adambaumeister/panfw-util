package deviceconfig

import (
	"encoding/xml"
	"fmt"
	"github.com/adambaumeister/panfw-util/panos/api"
	"github.com/adambaumeister/panfw-util/panos/errors"
	"github.com/schollz/progressbar"
	"log"
	"strings"
	"time"
)

func ImportNamed(fqdn string, apikey string, fn string) *MsgResponse {
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
	/*
		Commits the current configuration.

		This func waits for the commit to succeed before returning by polling the job id.
	*/
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

	job := ShowJob(fqdn, apikey, r.Job)
	bar := progressbar.NewOptions(100, progressbar.OptionSetRenderBlankState(true))
	for job.Status == "ACT" {
		job = ShowJob(fqdn, apikey, r.Job)
		bar.Add(job.Progress)
		time.Sleep(1 * time.Second)
	}
	bar.Finish()
}

func ShowJob(fqdn string, apikey string, jobid int) Job {
	/*
		Retrieve a JOB with id <jobid> on the system, returning a Job object.
	*/
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

	r := ShowJobResponse{}
	xml.Unmarshal(resp, &r)

	return r.Job
}

// Response received from commit or other messages that enqueue a job
type MsgJobResponse struct {
	api.Response
	Msg string `xml:"msg"`
	Job int    `xml:"result>job"`
}

// Representation of a show jobs id <blah>
type ShowJobResponse struct {
	api.Response
	Job Job `xml:"result>job"`
}

type CommitCommand struct {
	XMLName xml.Name
}

type showJobs struct {
	XMLName xml.Name
	Id      int `xml:"jobs>id"`
}

/*
PANOS "Job" API Object -> Returned from the output of certain show commands
*/
type Job struct {
	StartTime string `xml:"tenq"`
	EndTime   string `xml:"tdeq"`
	Id        int    `xml:"id"`
	User      string `xml:"user"`
	Type      string `xml:"type"`
	Status    string `xml:"status"`
	Queued    string `xml:"queued"`
	Stoppable string `xml:"stoppable"`
	Progress  int    `xml:"progress"`
	Warnings  string `xml:"warnings"`
	Details   string `xml:"details"`
}

func LoadNamedConfig(fqdn string, apikey string, cn string, commit bool) {
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

	r := MsgResponse{}
	xml.Unmarshal(resp, &r)
	// If the configuration is loaded, PAN device responds with 200 and status=success
	if r.Status == "success" {
		if commit {
			Commit(fqdn, apikey)
		}
	} else {
		r.Print()
		log.Fatal("Config import failed!\n")
	}
}

type LoadNamedCommand struct {
	XMLName xml.Name
	Config  string `xml:"config>from"`
}
