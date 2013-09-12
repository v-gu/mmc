package common

import "log"

const Debug Debugging = true

type Debugging bool

func (d Debugging) Printf(format string, args ...interface{}) {
	if d {
		log.Printf(format, args...)
	}
}
