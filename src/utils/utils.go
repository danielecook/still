package utils

import (
	"crypto/sha1"
	"fmt"
	"log"
)

// Quick error check
func Check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// Generate sha1 from string
func StringHash(Txt string) string {
	h := sha1.New()
	h.Write([]byte(Txt))
	bs := h.Sum(nil)
	sh := string(fmt.Sprintf("%x", bs))
	return sh
}
