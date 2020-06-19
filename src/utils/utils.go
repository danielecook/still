package utils

import "log"

// Quick error check
func Check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
