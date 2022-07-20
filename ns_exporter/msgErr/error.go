package msgErr

import "log"

func ErrInfo(err error) {
	if err != nil {
		log.Panic(err)
	}
}
