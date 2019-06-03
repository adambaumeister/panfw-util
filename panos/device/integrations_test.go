package device

import (
	"testing"
)

const USER_DEFAULT = "admin"
const PASS_DEFAULT = "admin"
const TESTING_IP_DEFAULT = "localhost:8443"

/*
Panos/device integration tests

These tests rely on a PANOS device (firewall) to be accessible at PANOS_IP (default: localhost:8443)
with username "admin" and password "admin" (the default).
*/

func _TestConnect(t *testing.T) {
	Connect(USER_DEFAULT, PASS_DEFAULT, TESTING_IP_DEFAULT)
}

func _TestFirewall_Rules(t *testing.T) {
	fw := Connect(USER_DEFAULT, PASS_DEFAULT, TESTING_IP_DEFAULT)
	fw.Rules()
}

func TestLoad(t *testing.T) {
	fw := Connect(USER_DEFAULT, PASS_DEFAULT, TESTING_IP_DEFAULT)
	fw.Load("panvm.xml")
}
