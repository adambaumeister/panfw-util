package api

import (
	"crypto/tls"
	"fmt"
	"github.com/adamb/panfw-util/panos/errors"
	"io/ioutil"
	"net/http"
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
	Fqdn   string
	Path   string
	Params map[string]string
}

func (q *QueryBase) AddParam(k string, v string) {
	q.Params[k] = v
}
func (q *QueryBase) SetPath(path string) {
	q.Path = path
}
func (q *QueryBase) SetFqdn(fqdn string) {
	q.Fqdn = fqdn
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
	return &pq
}

func (q *ParamQuery) Send() []byte {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%v%v", q.Fqdn, q.Path), nil)
	errors.DieIf(err)

	// Setup the full url
	u := req.URL.Query()
	for k, v := range q.Params {
		u.Add(k, v)
	}
	req.URL.RawQuery = u.Encode()

	errors.LogDebug(req.URL.String())
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
