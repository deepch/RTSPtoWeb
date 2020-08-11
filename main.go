package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	loggingPrintln("Server start")
	go HTTPAPIServer()
	go Storage.StreamRunAll()
	signalChanel := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(signalChanel, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-signalChanel
		loggingPrintln("Server receive signal", sig)
		done <- true
	}()
	loggingPrintln("Server start success a wait signals")
	<-done
	loggingPrintln("Server stop working by signal")
}
