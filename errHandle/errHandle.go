package errhandle

import "log"

func ErrHandling(err error) {
	if err != nil {
		log.Fatal(" : Error is ocured : ", err)
	}
}
