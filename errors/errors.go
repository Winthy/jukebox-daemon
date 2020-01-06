package errors

import "log"

func CaptureErr(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		// TODO: sentry
	}
}
