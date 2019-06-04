package api

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"github.com/adambaumeister/panfw-util/panos/errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

/*
Main Query broker
*/
type Query interface {
	Send()
}

/*
Base of a query, all queries need at least this much stuff to work
*/
type QueryBase struct {
	Fqdn    string
	Path    string
	Params  map[string]string
	Headers map[string]string
}

func (q *QueryBase) AddParam(k string, v string) {
	q.Params[k] = v
}
func (q *QueryBase) AddHeader(k string, v string) {
	q.Headers[k] = v
}
func (q *QueryBase) SetPath(path string) {
	q.Path = path
}
func (q *QueryBase) SetFqdn(fqdn string) {
	q.Fqdn = fqdn
}
func (q *QueryBase) EnableAuth(apikey string) {
	q.Params["key"] = apikey
}
func (q *QueryBase) SetMultipart(boundary string) {
	w := multipart.Writer{}
	w.SetBoundary(boundary)
	ct := w.FormDataContentType()

	q.AddHeader("Content-Type", ct)
}

func (q *QueryBase) Setup() *http.Request {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%v%v", q.Fqdn, q.Path), nil)
	errors.DieIf(err)

	// Setup the full url
	u := req.URL.Query()
	for k, v := range q.Params {
		u.Add(k, v)
	}
	for k, v := range q.Headers {
		req.Header.Add(k, v)
	}
	req.URL.RawQuery = u.Encode()
	return req
}

func (q *QueryBase) SetupPost(body io.Reader) *http.Request {
	req, err := http.NewRequest("POST", fmt.Sprintf("https://%v%v", q.Fqdn, q.Path), body)
	errors.DieIf(err)

	// Setup the full url
	u := req.URL.Query()
	for k, v := range q.Params {
		u.Add(k, v)
	}
	req.URL.RawQuery = u.Encode()
	for k, v := range q.Headers {
		req.Header.Add(k, v)
	}

	return req
}

func (q *QueryBase) GetUrlString() string {
	req := q.Setup()
	return req.URL.String()
}

/*
Simple GET query type
*/
type ParamQuery struct {
	QueryBase
}

// Initialize a new Parameter query
func NewParamQuery() *ParamQuery {
	pq := ParamQuery{}
	pq.Params = make(map[string]string)
	pq.Headers = make(map[string]string)
	return &pq
}

func (q *ParamQuery) Send() []byte {
	req := q.Setup()

	errors.LogDebug(req.URL.String())
	return SendHttpReq(req)
}

/*
Query containing an Xpath param
*/
type XpathQuery struct {
	QueryBase
	Xpath []string
}

func NewXpathQuery() *XpathQuery {
	xpq := XpathQuery{}
	xpq.Params = make(map[string]string)
	xpq.Headers = make(map[string]string)
	return &xpq
}

func (q *XpathQuery) SetXpath(xps []string) {
	q.Xpath = xps
}

func (q *XpathQuery) Send() []byte {
	xpath := MakeXPath(q.Xpath)
	q.AddParam("xpath", xpath)

	req := q.Setup()

	errors.LogDebug(req.URL.String())
	return SendHttpReq(req)
}

// Post request.
// Body must be set.
type Post struct {
	QueryBase
	Body io.Reader
}

func NewPost(data io.Reader) *Post {
	p := Post{
		Body: data,
	}
	p.Params = make(map[string]string)
	p.Headers = make(map[string]string)
	return &p
}

func (q *Post) Send() []byte {
	req := q.SetupPost(q.Body)
	errors.LogDebug(req.URL.String())
	return SendHttpReq(req)
}

type Cmd struct {
	QueryBase
	Command interface{}
}

func NewCmd(command interface{}) *Cmd {
	c := Cmd{
		Command: command,
	}
	c.Params = make(map[string]string)
	c.Headers = make(map[string]string)
	return &c
}

func (c *Cmd) Send() []byte {
	cmdxml, err := xml.Marshal(c.Command)
	errors.DieIf(err)
	c.AddParam("cmd", string(cmdxml))
	errors.LogDebug(string(cmdxml))

	req := c.Setup()
	return SendHttpReq(req)
}

/*
HTTP lib wrappers/convenience functions
*/
func SendHttpReq(req *http.Request) []byte {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{
		Transport: tr,
	}

	resp, err := client.Do(req)
	errors.DieIf(err)

	body, err := ioutil.ReadAll(resp.Body)
	errors.DieIf(err)

	return body
}

// Send a file from the local system using a HTTP POST
func HttpMultiPart(fn string) (*bytes.Buffer, string) {

	file, err := os.Open(fn)
	errors.DieIf(err)

	fileContents, err := ioutil.ReadAll(file)
	errors.DieIf(err)

	fi, err := file.Stat()
	errors.DieIf(err)

	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fi.Name())
	errors.DieIf(err)

	part.Write(fileContents)

	err = writer.Close()
	errors.DieIf(err)

	return body, writer.Boundary()
}

// Takes a list of string seps and converts to Xpath
// Not perfect but good enough for interacting with PAN.
func MakeXPath(path []string) string {
	return fmt.Sprintf("%v", strings.Join(path, "/"))
}
