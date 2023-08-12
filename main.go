package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Print("starting")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	sig := <-shutdown
	log.Print("shutdown on:", sig)
}
