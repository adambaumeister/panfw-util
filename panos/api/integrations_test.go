package api

import (
	"github.com/adamb/panfw-util/panos/api/auth"
	"testing"
)

const USER_DEFAULT = "admin"
const PASS_DEFAULT = "admin"
const TESTING_IP_DEFAULT = "localhost:8443"

/*
Panos/api integration tests

These tests rely on a PANOS device (firewall) to be accessible at PANOS_IP (default: localhost:8443)
with username "admin" and password "admin" (the default).
*/

func TestKeyGen(t *testing.T) {
	auth.KeyGen(USER_DEFAULT, PASS_DEFAULT, TESTING_IP_DEFAULT)
}
