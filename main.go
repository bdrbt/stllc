package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bdrbt/stllc/database"
	"github.com/bdrbt/stllc/internal/config"
	"github.com/bdrbt/stllc/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Cannot configure service %v", err)
	}

	if err = database.MigrateDB(cfg.PgURL()); err != nil {
		log.Fatalf("cannot run migrations:%v", err)
	}

	log.Print("starting")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	svc, err := service.New(cfg)
	if err != nil {
		log.Fatalf("cannot create service:%v", err)
	}

	if err := svc.Run(); err != nil {
		log.Fatalf("error starting service:%v", err)
	}

	svc.Run()

	sig := <-shutdown
	log.Print("shutdown on:", sig)
}
