package errors

import (
	"log"
)

var DEBUG = false

func DieIf(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func LogDebug(s string) {
	if DEBUG {
		log.Printf("%v", s)
	}
}
