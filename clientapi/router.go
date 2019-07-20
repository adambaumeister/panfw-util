package clientapi

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/adambaumeister/panfw-util/panos/device"
	"github.com/adambaumeister/panfw-util/panos/errors"
	"io/ioutil"
	"log"
	"net/http"
)

type ApiHandler struct {
	loginkey string

	pandevice device.Panos
}

func Start() {
	a := ApiHandler{}

	http.HandleFunc("/status", Test)
	http.HandleFunc("/login", a.Login)
	http.HandleFunc("/loginstatus", a.LoginStatus)
	http.HandleFunc("/join", a.Join)
	http.HandleFunc("/register", a.Register)

	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func Test(w http.ResponseWriter, r *http.Request) {
	// Return the status of Panutil backend
	jm := StatusMessage{
		Status:  0,
		Message: "Panutil is running.",
	}
	j, err := json.Marshal(jm)
	if err != nil {
		fmt.Println("error:", err)
	}

	w.Write(j)
}

func (a *ApiHandler) LoginStatus(w http.ResponseWriter, r *http.Request) {
	var resp LoginResponse
	if a.pandevice != nil {
		resp = LoginResponse{
			ApiKey: a.pandevice.GetApiKey(),
			Status: 0,
		}
		j, _ := json.Marshal(resp)
		w.Write(j)
		return
	}

	resp = LoginResponse{
		ApiKey: "",
		Status: 1,
	}
	j, _ := json.Marshal(resp)
	w.Write(j)
	return
}

func (a *ApiHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Attempt to login and retrieve an API key from the configured firewall.
	// If no Hostname is configured, return an error.
	var resp LoginResponse

	lr := LoginRequest{}
	body, _ := ioutil.ReadAll(r.Body)
	errors.LogDebug(string(body))
	json.Unmarshal(body, &lr)
	d := device.ConnectUniversal(lr.Username, lr.Password, lr.Hostname)

	if d == nil {
		resp = LoginResponse{
			ApiKey: "",
			Status: 1,
		}
		j, _ := json.Marshal(resp)
		w.Write(j)
		return
	}

	lk := loginUserKey(10)
	resp = LoginResponse{
		ApiKey: lk,
		Status: 0,
	}
	j, _ := json.Marshal(resp)
	w.Write(j)
	a.pandevice = d
	a.loginkey = lk
}

func (a *ApiHandler) Join(w http.ResponseWriter, r *http.Request) {
	cr := Command{}
	body, _ := ioutil.ReadAll(r.Body)
	errors.LogDebug(string(body))
	json.Unmarshal(body, &cr)

	if !a.CheckLogin(cr.ApiKey) {
		sr := StatusMessage{
			Status:  1,
			Message: "Not logged in.",
		}
		j, _ := json.Marshal(sr)
		w.Write(j)
		return
	}
	sr := StatusMessage{
		Status:  0,
		Message: "Command OK.",
	}
	j, _ := json.Marshal(sr)
	w.Write(j)
}

func (a *ApiHandler) Register(w http.ResponseWriter, r *http.Request) {
	cr := Command{}
	body, _ := ioutil.ReadAll(r.Body)
	errors.LogDebug(string(body))
	json.Unmarshal(body, &cr)

	if !a.CheckLogin(cr.ApiKey) {
		sr := StatusMessage{
			Status:  1,
			Message: "Not logged in.",
		}
		j, _ := json.Marshal(sr)
		w.Write(j)
		return
	}

	result := a.pandevice.Register(cr.Args)
	//fmt.Printf("Result: %v\n", result.Status)

	sr := StatusMessage{
		Status:  0,
		Message: result.Status,
	}
	j, _ := json.Marshal(sr)
	w.Write(j)
}

func (a *ApiHandler) CheckLogin(key string) bool {
	if key == a.loginkey {
		return true
	}
	return false
}

func loginUserKey(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	s := fmt.Sprintf("%X", b)
	return s
}
