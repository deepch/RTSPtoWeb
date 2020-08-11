package main

import "log"

func loggingPrintln(v ...interface{}) {
	if debug {
		log.Println(v...)
	}
}
