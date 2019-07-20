package clientapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestApiHandler_Register(t *testing.T) {
	go Start()
	lr := LoginRequest{
		Username: "admin",
		Password: os.Getenv("TEST_PASSWORD"),
		Hostname: os.Getenv("TEST_HOSTNAME"),
	}
	j, _ := json.Marshal(lr)
	resp, err := http.Post("http://127.0.0.1:8080"+"/login", "application/json", bytes.NewBuffer(j))
	if err != nil {
		panic(err)
	}

	r := LoginResponse{}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &r)
	fmt.Printf("%v %v\n", r.Status, r.ApiKey)

	cr := Command{
		ApiKey:  r.ApiKey,
		Command: "register",
		Args:    []string{"1.1.1.2", "servers"},
	}
	j, _ = json.Marshal(cr)
	resp, err = http.Post("http://127.0.0.1:8080"+"/register", "application/json", bytes.NewBuffer(j))
	if err != nil {
		panic(err)
	}
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Printf(string(body))
}
