package Input

import (
	"github.com/adambaumeister/panfw-util/panos/errors"
	"testing"
)

func TestToObjects(t *testing.T) {
	errors.DEBUG = true
	a := []string{"address,test_addr1,1.1.1.1\naddress,test_addr2,2.2.2.2"}
	objects := ToObjects(a)

	if objects == nil {
		errors.LogDebug("Nil returned.")
		t.Fail()
	}
	if len(objects) != 2 {
		errors.LogDebug("Incorrect object count returned.")
		t.Fail()
	}
}
