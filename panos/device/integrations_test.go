package device

import (
	"os"
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

func _TestLoad(t *testing.T) {
	fw := Connect(USER_DEFAULT, PASS_DEFAULT, TESTING_IP_DEFAULT)
	fw.ImportNamed("panvm.xml")
}

func TestCommit(t *testing.T) {
	if os.Getenv("NO_TEST") != "" {
		t.Skip()
	}
	fw := Connect(USER_DEFAULT, PASS_DEFAULT, TESTING_IP_DEFAULT)
	fw.LoadNamed("panvm.xml", false)
	fw.Commit()
}
