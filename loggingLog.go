package main

import "log"

//loggingPrintln logging function
func loggingPrintln(v ...interface{}) {
	if debug {
		log.Println(v...)
	}
}
