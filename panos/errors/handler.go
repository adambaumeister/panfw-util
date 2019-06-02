package errors

import (
	"log"
)

func DieIf(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func LogDebug(s string) {
	log.Printf("%v", s)
}
