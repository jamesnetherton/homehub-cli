package homehub

import (
	"log"
)

type debugging bool

func (debug debugging) Printf(format string, args ...interface{}) {
	if debug {
		log.Printf(format, args...)
	}
}

func (debug debugging) Print(args ...interface{}) {
	if debug {
		log.Print(args...)
	}
}

func (debug debugging) Println(args ...interface{}) {
	if debug {
		log.Println(args...)
	}
}
