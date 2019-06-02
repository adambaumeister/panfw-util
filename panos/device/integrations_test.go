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

func TestConnect(t *testing.T) {
	Connect(USER_DEFAULT, PASS_DEFAULT, TESTING_IP_DEFAULT)
}
